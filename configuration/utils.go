package configuration

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"strconv"
)

type remoteHost struct {
	Name string
	Host string
	Port int
}

// Configuration is the server configuration
type Configuration struct {
	Name string
	Host string
	Port int
	Databases struct {
		Elasticsearch remoteHost
	}
	Services []remoteHost
}

// GetConfig retreve the server configuration
func GetConfig(path string) Configuration {

	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	conf := Configuration{}
	yamlErr := yaml.Unmarshal([]byte(file), &conf)
	if err != nil {
		panic(yamlErr)
	}

	return conf
}

// GetServicesHost return all services ip and host coma separated
func (c *Configuration) GetServicesHost() string {
	var host string

	for key, service := range c.Services {
		if key != 0 {
			host += ","
		}
		host += service.Host + ":" + strconv.Itoa(service.Port)
	}

	return host
}

// GetServicePort return the port of the given service name
func (c *Configuration) GetServicePort(service string) int {
	var port int

	var nameToSearch string

	if service == "all" {
		nameToSearch = "master"
	} else {
		nameToSearch = service
	}

	for _, service := range c.Services {
		if nameToSearch == service.Name {
			port = service.Port
		}
	}

	return port
}
