package domainManager

import (
  "github.com/greg-nicolle/go-microservice/transport"
  "net/http"
  stdprometheus "github.com/prometheus/client_golang/prometheus"
  "context"
  httptransport "github.com/go-kit/kit/transport/http"
  "strconv"
  "github.com/Sirupsen/logrus"
)

// Domain is a
type Domain struct {
  domainRegistreted map[string]transport.Service
  logger            logrus.Logger
  port              int
  ctx               context.Context
}

// Create return a new Domain instance
func Create(port int, logger logrus.Logger) Domain {

  ctx := context.Background()

  return Domain{domainRegistreted:map[string]transport.Service{},
    logger:logger, port:port, ctx:ctx}
}

// AddService add a new transport.Service to the domainManager
func (d *Domain)AddService(service transport.Service) {
  d.domainRegistreted[service.GetServiceName()] = service
}

// LaunchService launch services
func (d *Domain)LaunchService(serviceName string, instancesIP string) {

  if service, isPresent := d.domainRegistreted[serviceName]; isPresent {
    d.domainRegistreted = map[string]transport.Service{service.GetServiceName():service}
  }

  for _, service := range d.domainRegistreted {
    for _, endpoint := range service.GetServiceEndpoints() {
      handler := httptransport.NewServer(
        endpoint.MakeEndpoint(service.GetService(d.ctx, instancesIP,
          d.logger)),
        transport.DecodeRequest(endpoint.GetIo().Request),
        transport.EncodeResponse,
      )
      http.Handle(endpoint.GetIo().Path, handler)
    }
  }

  http.Handle("/metrics", stdprometheus.Handler())
  d.logger.Info("msg", "HTTP", "addr", d.port)
  d.logger.Info("err", http.ListenAndServe(":" + strconv.Itoa(d.port), nil))
}