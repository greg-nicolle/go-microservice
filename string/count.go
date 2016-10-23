package stringModule

import (
  "golang.org/x/net/context"

  "github.com/go-kit/kit/endpoint"
  "github.com/greg-nicolle/kit-test/transport"
)

type CountEndpoint struct{}

func (CountEndpoint) MakeEndpoint(svc interface{}) endpoint.Endpoint {
  return func(ctx context.Context, request interface{}) (interface{}, error) {
    req := request.(countRequest)
    v, err := svc.(StringService).Count(req.S)
    if err != nil {
      return countResponse{v, err.Error()}, nil
    }
    return countResponse{v, ""}, nil
  }
}

func (CountEndpoint) GetIo() transport.Io {
  return transport.Io{
    Request: countRequest{},
    Response: countResponse{},
    Path: "/count"}
}

type countRequest struct {
  S string `json:"s"`
}

type countResponse struct {
  V   int `json:"v"`
  Err string `json:"err,omitempty"`
}
