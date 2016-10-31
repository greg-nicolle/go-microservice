package domainManager

import (
  "github.com/greg-nicolle/go-microservice/transport"
  "net/http"
  stdprometheus "github.com/prometheus/client_golang/prometheus"
  "golang.org/x/net/context"
  "github.com/go-kit/kit/log"
  httptransport "github.com/go-kit/kit/transport/http"
  kitprometheus "github.com/go-kit/kit/metrics/prometheus"
  "strconv"
)
// Domain is a
type Domain struct {
  domainRegistreted map[string]transport.Service
  logger            log.Logger
  port              int
  ctx               context.Context
}

// Create return a new Domain instance
func Create(port int, logger            log.Logger) Domain {

  ctx := context.Background()

  return Domain{domainRegistreted:map[string]transport.Service{},
    logger:logger, port:port, ctx:ctx}
}

// AddService add a new transport.Service to the domainManager
func (d *Domain)AddService(service transport.Service) {
  d.domainRegistreted[service.GetServiceName()] = service
}

// LaunchService launch services
func (d *Domain)LaunchService(serviceName string, instancesIP []string) {

  //  instrument stuff
  fieldKeys := []string{"method", "error"}
  requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
    Namespace: "my_group",
    Subsystem: "string_service",
    Name:      "request_count",
    Help:      "Number of requests received.",
  }, fieldKeys)
  requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
    Namespace: "my_group",
    Subsystem: "string_service",
    Name:      "request_latency_microseconds",
    Help:      "Total duration of requests in microseconds.",
  }, fieldKeys)
  countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
    Namespace: "my_group",
    Subsystem: "string_service",
    Name:      "count_result",
    Help:      "The result of each count method.",
  }, []string{})

  if service, isPresent := d.domainRegistreted[serviceName]; isPresent {
    d.domainRegistreted = map[string]transport.Service{service.GetServiceName():service}
  }

  for _, service := range d.domainRegistreted {
    for _, endpoint := range service.GetServiceEndpoints() {
      handler := httptransport.NewServer(
        d.ctx,
        endpoint.MakeEndpoint(service.GetService(d.ctx, instancesIP,
          d.logger,
          requestCount,
          requestLatency,
          countResult)),
        transport.DecodeRequest(endpoint.GetIo().Request),
        transport.EncodeResponse,
      )
      http.Handle(endpoint.GetIo().Path, handler)
    }
  }

  http.Handle("/metrics", stdprometheus.Handler())
  d.logger.Log("msg", "HTTP", "addr", d.port)
  d.logger.Log("err", http.ListenAndServe(":" + strconv.Itoa(d.port), nil))
}