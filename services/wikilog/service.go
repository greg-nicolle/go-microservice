package wikilog

import (
	"github.com/greg-nicolle/go-microservice/transport"
	"context"
	"github.com/Sirupsen/logrus"
	"github.com/greg-nicolle/go-microservice/configuration"
)

// Wikilog implement Service
type Wikilog struct{}

// GetServiceEndpoints implement GetServiceEndpoints of String
func (Wikilog) GetServiceEndpoints() []transport.GEndpoint {
	return []transport.GEndpoint{
		searchPageNameEndpoint{}}
}

// GetService implement GetService of String
func (Wikilog) GetService(ctx context.Context,
	instances string,
	logger *logrus.Logger,
	config configuration.Configuration) interface{} {

	var svc Domain
	svc = wikilogDomain{config}
	svc = proxyingMiddleware(ctx, instances, logger)(svc)

	return svc
}

// GetServiceName implement GetServiceName of String
func (Wikilog) GetServiceName() string {
	return "wikilog"
}
