package stringModule

import (
  "github.com/greg-nicolle/go-microservice/transport"
  "golang.org/x/net/context"
  "github.com/go-kit/kit/metrics"
  "github.com/Sirupsen/logrus"
)

// String implement Service
type String struct{}

// GetServiceEndpoints implement GetServiceEndpoints of String
func (String)GetServiceEndpoints() []transport.GEndpoint {
  return []transport.GEndpoint{
    uppercaseEndpoint{},
    countEndpoint{}}
}

// GetService implement GetService of String
func (String) GetService(ctx context.Context,
instances string,
logger logrus.Entry,
requestCount metrics.Counter,
requestLatency metrics.Histogram,
countResult metrics.Histogram) interface{} {

  var svc StringService
  svc = stringService{}
  svc = proxyingMiddleware(ctx, instances, logger)(svc)
  svc = loggingMiddleware(logger)(svc)
  svc = instrumentingMiddleware(requestCount, requestLatency, countResult)(svc)

  return svc
}

// GetServiceName implement GetServiceName of String
func (String)GetServiceName() string {
  return "string"
}