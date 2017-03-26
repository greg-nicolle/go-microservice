package wikilog

import (
	"github.com/go-kit/kit/endpoint"
	"context"
	"github.com/greg-nicolle/go-microservice/transport"
)

type searchPageNameEndpoint struct{}

func (searchPageNameEndpoint) MakeEndpoint(svc interface{}) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(searchPageNameRequest)
		v, err := svc.(Domain).searchPageName(req.S)
		if err != nil {
			return searchPageNameResponse{v, err.Error()}, nil
		}
		return searchPageNameResponse{v, ""}, nil
	}
}

func (searchPageNameEndpoint) GetIo() transport.Io {
	return transport.Io{
		Request:  searchPageNameRequest{},
		Response: searchPageNameResponse{},
		Path:     "/getPageName"}
}

type searchPageNameRequest struct {
	S string `json:"s"`
}

type searchPageNameResponse struct {
	V   []string `json:"v"`
	Err string `json:"err,omitempty"`
}
