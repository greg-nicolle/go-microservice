package domainManager

import (
  "testing"
  "github.com/greg-nicolle/go-microservice/transport"
  "github.com/stretchr/testify/assert"
  "github.com/go-kit/kit/log"
)

func TestAddService(t *testing.T) {

  // Given
  service := new(transport.MockService)
  service.On("GetServiceName").Return("test")
  expected := map[string]transport.Service{}
  expected["test"] = service
  d := Domain{domainRegistreted:map[string]transport.Service{}}

  // When
  d.AddService(service)

  // Then
  assert.Equal(t, expected, d.domainRegistreted)
}

func TestCreate(t *testing.T) {

  // Given
  port := 8080
  var logger log.Logger
  expected := Domain{domainRegistreted:map[string]transport.Service{},logger:logger,port:port}

  // When
  actual := Create(port, logger)

  // Then
  assert.Equal(t, expected.port, actual.port)
  assert.Equal(t, expected.logger, actual.logger)
  assert.Equal(t, expected.domainRegistreted, actual.domainRegistreted)
}

// TODO add tests for LaunchService