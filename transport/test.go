package transport

import (
  "github.com/stretchr/testify/mock"
  "context"
  "github.com/Sirupsen/logrus"

)

// MockService is an mock implementation of service
type MockService struct {
  mock.Mock
}

// GetServiceEndpoints is an mock implementation of GetServiceEndpoints
func (m *MockService) GetServiceEndpoints() []GEndpoint {
  args := m.Called()
  return args.Get(0).([]GEndpoint)
}

// GetService is an mock implementation of GetService
func (m *MockService)GetService(ctx context.Context, instancesIP string, logger logrus.Logger) interface{} {
  args := m.Called()
  return args.Get(0).(*[]GEndpoint)
}

// GetServiceName is a mock implementation of GetServiceName
func (m *MockService)GetServiceName() string{
  args := m.Called()
  return args.String(0)
}