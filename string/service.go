package stringModule

import (
  "github.com/greg-nicolle/kit-test/transport"
  "github.com/go-kit/kit/log"
  "golang.org/x/net/context"
  "github.com/go-kit/kit/metrics"
)

type String struct{}

func (String)GetServiceEndpoints() []transport.GEndpoint {
  return []transport.GEndpoint{UppercaseEndpoint{},
    CountEndpoint{}}
}

func (String) GetService(instances string,
ctx context.Context,
logger log.Logger,
requestCount metrics.Counter,
requestLatency metrics.Histogram,
countResult metrics.Histogram) interface{} {

  var svc StringService
  svc = stringService{}
  svc = proxyingMiddleware(instances, ctx, logger)(svc)
  svc = loggingMiddleware(logger)(svc)
  svc = instrumentingMiddleware(requestCount, requestLatency, countResult)(svc)

  return svc
}

func (String)GetServiceName() string {
  return "string"
}