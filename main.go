package main

import (
  "flag"
  "net/http"
  "os"
  "context"
  "github.com/Sirupsen/logrus"
  httptransport "github.com/go-kit/kit/transport/http"
  "github.com/greg-nicolle/go-microservice/transport"
  "github.com/greg-nicolle/go-microservice/string"
)

func main() {
  var (
    listen = flag.String("listen", ":8080", "HTTP listen address")
    proxy = flag.String("proxy", "", "Optional comma-separated list of URLs to proxy uppercase requests")
    serviceName = flag.String("service", "", "Service you want to start")
  )
  flag.Parse()

  var log = logrus.New()

  log.Out = os.Stdout
  logger := log.WithFields(logrus.Fields{
    "animal": "walrus",
    "size":   10,
  })

  ctx := context.Background()

  services := map[string]transport.Service{}
  services[stringModule.String{}.GetServiceName()] = stringModule.String{}

  if service, isPresent := services[*serviceName]; isPresent {
    services = map[string]transport.Service{service.GetServiceName():service}
  }

  for _, service := range services {
    for _, endpoint := range service.GetServiceEndpoints() {
      handler := httptransport.NewServer(
        endpoint.MakeEndpoint(service.GetService(ctx, *proxy, *logger)),
        transport.DecodeRequest(endpoint.GetIo().Request),
        transport.EncodeResponse,
      )
      http.Handle(endpoint.GetIo().Path, handler)
    }
  }

  logger.Info("msg", "HTTP", "addr", *listen)
  logger.Info("err", http.ListenAndServe(*listen, nil))
}
