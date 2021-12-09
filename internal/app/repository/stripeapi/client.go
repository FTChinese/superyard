package stripeapi

import (
	"github.com/FTChinese/superyard/pkg/config"
	"github.com/stripe/stripe-go/v72/client"
	"go.uber.org/zap"
)

type Client struct {
	sc     *client.API
	logger *zap.Logger
}

func newClient(key string, logger *zap.Logger) Client {
	return Client{
		sc:     client.New(key, nil),
		logger: logger,
	}
}

type Clients struct {
	test Client
	live Client
}

func NewClients(logger *zap.Logger) Clients {
	secrets := config.MustStripeSecret()

	return Clients{
		test: newClient(secrets.Pick(false), logger),
		live: newClient(secrets.Pick(true), logger),
	}
}

func (c Clients) Select(live bool) Client {
	if live {
		return c.live
	}

	return c.test
}
