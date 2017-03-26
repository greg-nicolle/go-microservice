package main

import (
	"flag"
	"os"
	"github.com/Sirupsen/logrus"
	"github.com/greg-nicolle/go-microservice/services/string"
	"github.com/greg-nicolle/go-microservice/domainManager"
	"github.com/greg-nicolle/go-microservice/services/wikilog"
	"github.com/greg-nicolle/go-microservice/configuration"
)

func main() {
	var (
		serviceName = flag.String("service", "", "Service you want to start")
		configPath  = flag.String("configPath", "", "Config path")
	)
	flag.Parse()

	config := configuration.GetConfig(*configPath)

	var log = logrus.New()

	log.Out = os.Stdout

	log.Info(*serviceName)

	domains := domainManager.Create(log)

	domains.AddService(stringModule.String{})
	domains.AddService(wikilog.Wikilog{})
	domains.LaunchService(*serviceName, config)
}
