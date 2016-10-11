package main

import (
  "flag"
  "net/http"
  "os"

  stdprometheus "github.com/prometheus/client_golang/prometheus"
  "golang.org/x/net/context"

  "github.com/go-kit/kit/log"
  httptransport "github.com/go-kit/kit/transport/http"
  "github.com/greg-nicolle/kit-test/transport"
  "github.com/greg-nicolle/kit-test/string"
  kitprometheus "github.com/go-kit/kit/metrics/prometheus"
)

func main() {
  var (
    listen = flag.String("listen", ":8080", "HTTP listen address")
    proxy = flag.String("proxy", "", "Optional comma-separated list of URLs to proxy uppercase requests")
  )
  flag.Parse()

  var logger log.Logger
  logger = log.NewLogfmtLogger(os.Stderr)
  logger = log.NewContext(logger).With("listen", *listen).With("caller", log.DefaultCaller)

  ctx := context.Background()

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

  var services []transport.Service
  services = append(services, stringModule.String{})

  for _, service := range services {

    for _, endpoint := range service.GetServiceEndpoints() {
      handler := httptransport.NewServer(
        ctx,
        endpoint.MakeEndpoint(service.GetService(*proxy, ctx, logger, requestCount,
        requestLatency,
        countResult )),
        transport.DecodeRequest(endpoint.GetIo().Request),
        transport.EncodeResponse,
      )
      http.Handle(endpoint.GetIo().Path, handler)
    }
  }

  http.Handle("/metrics", stdprometheus.Handler())
  logger.Log("msg", "HTTP", "addr", *listen)
  logger.Log("err", http.ListenAndServe(*listen, nil))
}
