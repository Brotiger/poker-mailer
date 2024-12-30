package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Brotiger/poker-mailer/internal/config"
	"github.com/Brotiger/poker-mailer/internal/controller"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	log "github.com/sirupsen/logrus"
)

func main() {
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	options := []nats.Option{
		nats.RetryOnFailedConnect(config.Cfg.Nats.RetryOnFailedConnect),
		nats.MaxReconnects(config.Cfg.Nats.MaxReconnects),
		nats.ReconnectWait(time.Duration(config.Cfg.Nats.ReconnectWait) * time.Millisecond),
		nats.PingInterval(time.Duration(config.Cfg.Nats.PingInterval) * time.Millisecond),
		nats.MaxPingsOutstanding(config.Cfg.Nats.MaxPingOutstanding),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			log.Warnf("nats-connect: disconnected, error: %v", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Trace("nats-connect: reconnected")
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			if nc.Status() == nats.CLOSED && nc.LastError() != nil {
				log.Fatalf("nats-connect: connection closed after max reconnects: %v", nc.LastError())
			} else if nc.LastError() != nil {
				log.Errorf("nats-connect: connection closed, error: %v", nc.LastError())
			} else {
				log.Error("nats-connect: connection closed")
			}
		}),
	}

	if config.Cfg.Nats.ClientCert != "" && config.Cfg.Nats.ClientKey != "" {
		options = append(options, nats.ClientCert(config.Cfg.Nats.ClientCert, config.Cfg.Nats.ClientKey))
	}

	if config.Cfg.Nats.RootCA != "" {
		options = append(options, nats.RootCAs(config.Cfg.Nats.RootCA))
	}

	natsConn, err := nats.Connect(config.Cfg.Nats.Addr, options...)
	if err != nil {
		log.Fatalf("failed to connect to nats, error: %v", err)
	}

	js, err := jetstream.New(natsConn)
	if err != nil {
		log.Fatalf("failed to connect to jet stream, error: %v", err)
	}

	stream, err := js.Stream(ctx, config.Cfg.Nats.Stream)
	if err != nil {
		log.Fatalf("failed to connect with stream, error: %v", err)
	}

	cons, err := stream.Consumer(ctx, config.Cfg.Nats.ConsumerName)
	if err != nil {
		log.Fatalf("failed to connect with consumer, error: %v", err)
	}

	mailController, err := controller.NewMailController()
	if err != nil {
		log.Fatalf("failed to create mail controller, error: %v", err)
	}

	log.Info("application started")

	var wg sync.WaitGroup
	consumeCtx, err := cons.Consume(func(msg jetstream.Msg) {
		wg.Add(1)
		defer wg.Done()

		reqId := (uuid.New()).String()

		log.Tracef("%s, incoming message: %s", reqId, msg.Data())

		if err := msg.Ack(); err != nil {
			log.Errorf("%s failed to ack msg, error: %v", reqId, err)
		}

		if err := mailController.Send(ctx, msg); err != nil {
			log.Errorf("%s failed to send, error: %v", reqId, err)
		}
	}, jetstream.PullMaxMessages(config.Cfg.Nats.BatchSize))
	if err != nil {
		log.Fatalf("failed to consume, error: %v", err)
	}

	shutdown := make(chan os.Signal, 1)
	defer close(shutdown)

	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	gracefulShutdown(consumeCtx, &wg)
}

func gracefulShutdown(consumeCtx jetstream.ConsumeContext, wg *sync.WaitGroup) {
	consumeCtx.Stop()

	done := make(chan struct{})
	defer close(done)
	go func() {
		wg.Wait()
		done <- struct{}{}
	}()

	log.Info("graceful shutdown: waiting for task completion")
	select {
	case <-done:
		log.Info("graceful shutdown: task completed successfully")
	case <-time.Tick(time.Duration(config.Cfg.App.GracefulShutdownTimeoutMS) * time.Millisecond):
		log.Info("graceful shutdown: timeout, forcing shutdown")
	}
}
