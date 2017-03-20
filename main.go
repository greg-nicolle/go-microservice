package main

import (
  "flag"
  "os"
  "github.com/Sirupsen/logrus"
  "github.com/greg-nicolle/go-microservice/string"
  "github.com/greg-nicolle/go-microservice/domainManager"
)

func main() {
  var (
    listen = flag.Int("listen", 8080, "HTTP listen address")
    proxy = flag.String("proxy", "", "Optional comma-separated list of URLs to proxy uppercase requests")
    serviceName = flag.String("service", "", "Service you want to start")
  )
  flag.Parse()

  var log = logrus.New()

  log.Out = os.Stdout

  domains := domainManager.Create(*listen, *log)

  domains.AddService(stringModule.String{})
  domains.LaunchService(*serviceName, *proxy)
}
