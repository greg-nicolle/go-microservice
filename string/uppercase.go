package stringModule

import (
  "golang.org/x/net/context"

  "github.com/go-kit/kit/endpoint"
  "github.com/greg-nicolle/go-microservice/transport"
)

type UppercaseEndpoint struct{}

func (UppercaseEndpoint) MakeEndpoint(svc interface{}) endpoint.Endpoint {
  return func(ctx context.Context, request interface{}) (interface{}, error) {
    req := request.(uppercaseRequest)
    v, err := svc.(StringService).Uppercase(req.S)
    if err != nil {
      return uppercaseResponse{v, err.Error()}, nil
    }
    return uppercaseResponse{v, ""}, nil
  }
}

func (UppercaseEndpoint) GetIo() transport.Io {
  return transport.Io{
    Request: uppercaseRequest{},
    Response: uppercaseResponse{},
    Path: "/uppercase"}
}

type uppercaseRequest struct {
  S string `json:"s"`
}

type uppercaseResponse struct {
  V   string `json:"v"`
  Err string `json:"err,omitempty"`
}