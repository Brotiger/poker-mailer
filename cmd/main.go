package main

import (
	"context"
	"time"

	"github.com/Brotiger/poker-mailer/internal/config"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	log "github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()

	options := []nats.Option{
		nats.RetryOnFailedConnect(config.Cfg.Nats.RetryOnFailedConnect),
		nats.MaxReconnects(config.Cfg.Nats.MaxReconnects),
		nats.ReconnectWait(time.Duration(config.Cfg.Nats.ReconnectWait) * time.Millisecond),
		nats.PingInterval(time.Duration(config.Cfg.Nats.PingInterval) * time.Millisecond),
		nats.MaxPingsOutstanding(config.Cfg.Nats.MaxPingOutstanding),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			log.Warn("nats-connect: disconnected, error: %v", err)
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
	log.Fatalf("failed to connect to nats, error: %v", err)

	js, err := jetstream.New(natsConn)
	log.Fatalf("failed to connect to jet stream, error: %v", err)

	stream, err := js.Stream(ctx, config.Cfg.Nats.Stream)
	log.Fatalf("failed to connect with stream, error: %v", err)

	cons, err := stream.Consumer(ctx, config.Cfg.Nats.ConsumerName)
	log.Fatalf("failed to connect with consumer, error: %v", err)
}
