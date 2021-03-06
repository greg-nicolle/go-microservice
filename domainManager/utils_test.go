package domainManager

import (
	"testing"
	"github.com/greg-nicolle/go-microservice/transport"
	"github.com/stretchr/testify/assert"
	"github.com/Sirupsen/logrus"
	"os"
)

func TestAddService(t *testing.T) {

	// Given
	service := new(transport.MockService)
	service.On("GetServiceName").Return("test")
	expected := map[string]transport.Service{}
	expected["test"] = service
	d := Domain{domainRegistreted: map[string]transport.Service{}}

	// When
	d.AddService(service)

	// Then
	assert.Equal(t, expected, d.domainRegistreted)
}

func TestCreate(t *testing.T) {

	// Given
	var log = logrus.New()
	log.Out = os.Stdout

	expected := Domain{domainRegistreted: map[string]transport.Service{}, logger: log}

	// When
	actual := Create(log)

	// Then
	assert.Equal(t, expected.logger, actual.logger)
	assert.Equal(t, expected.domainRegistreted, actual.domainRegistreted)
}

// TODO add tests for LaunchService
