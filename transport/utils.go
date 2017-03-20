package transport

import (
  "reflect"
  "encoding/json"
  "io/ioutil"
  "bytes"
  "net/http"
  "context"
  httptransport "github.com/go-kit/kit/transport/http"
  "github.com/go-kit/kit/endpoint"
  "github.com/Sirupsen/logrus"
)
// DecodeRequest is a generic implementation for decoding a request
func DecodeRequest(i interface{}) httptransport.DecodeRequestFunc {
  request := reflect.New(reflect.TypeOf(i)).Interface()
  return func(_ context.Context, r *http.Request) (interface{}, error) {
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
      return nil, err
    }
    i2 := reflect.ValueOf(request).Elem().Interface()
    return i2, nil
  }
}

// DecodeResponse is a generic implementation for decoding a response
func DecodeResponse(i interface{}) httptransport.DecodeResponseFunc {
  request := reflect.New(reflect.TypeOf(i)).Interface()
  return func(_ context.Context, r *http.Response) (interface{}, error) {
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
      return nil, err
    }
    i2 := reflect.ValueOf(request).Elem().Interface()
    return i2, nil
  }
}

// EncodeResponse is a generic implementation for encoding response
func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
  return json.NewEncoder(w).Encode(response)
}

// EncodeRequest is a generic implementation for encoding request
func EncodeRequest(_ context.Context, r *http.Request, request interface{}) error {
  var buf bytes.Buffer
  if err := json.NewEncoder(&buf).Encode(request); err != nil {
    return err
  }
  r.Body = ioutil.NopCloser(&buf)
  return nil
}

// Io is an interface that discribe the io of an endpoint
type Io struct {
  Request  interface{}
  Response interface{}
  Path     string
}

// GEndpoint is an interface that describe an endpoint
type GEndpoint interface {
  GetIo() Io
  MakeEndpoint(interface{}) endpoint.Endpoint
}

// Service is an interface that describe a service
type Service interface {
  GetServiceEndpoints() []GEndpoint
  GetService(ctx context.Context, instances string, logger logrus.Entry) interface{}
  GetServiceName() string
}