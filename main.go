package main

import (
  "flag"
  "github.com/greg-nicolle/go-microservice/string"
  "strings"
  "github.com/greg-nicolle/go-microservice/domainManager"
  "os"
  "github.com/go-kit/kit/log"
)

func main() {
  var (
    listen = flag.Int("listen", 8080, "HTTP listen address")
    proxy = flag.String("proxy", "", "Optional comma-separated list of URLs to proxy uppercase requests")
    serviceName = flag.String("service", "", "Service you want to start")
  )
  flag.Parse()

  var logger log.Logger
  logger = log.NewLogfmtLogger(os.Stderr)
  logger = log.NewContext(logger).With("listen", *listen).With("caller", log.DefaultCaller)

  domains := domainManager.Create(*listen, logger)

  domains.AddService(stringModule.String{})
  domains.LaunchService(*serviceName, split(*proxy))
}

func split(s string) []string {
  a := strings.Split(s, ",")
  for i := range a {
    a[i] = strings.TrimSpace(a[i])
  }
  return a
}