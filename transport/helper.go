package transport

import (
  "reflect"
  "encoding/json"
  "io/ioutil"
  "bytes"
  "net/http"
  "golang.org/x/net/context"
  "github.com/go-kit/kit/log"

  httptransport "github.com/go-kit/kit/transport/http"
  "github.com/go-kit/kit/endpoint"
  "github.com/go-kit/kit/metrics"
)

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

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
  return json.NewEncoder(w).Encode(response)
}

func EncodeRequest(_ context.Context, r *http.Request, request interface{}) error {
  var buf bytes.Buffer
  if err := json.NewEncoder(&buf).Encode(request); err != nil {
    return err
  }
  r.Body = ioutil.NopCloser(&buf)
  return nil
}

type Io struct {
  Request  interface{}
  Response interface{}
  Path     string
}

type GEndpoint interface {
  GetIo() Io
  MakeEndpoint(interface{}) endpoint.Endpoint
}

type Service interface {
  GetServiceEndpoints() []GEndpoint
  GetService(instances string, ctx context.Context, logger log.Logger, requestCount metrics.Counter,
  requestLatency metrics.Histogram,
  countResult metrics.Histogram) interface{}
}