package wikilog

import (
	"time"
	"github.com/Sirupsen/logrus"
	"errors"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"context"
	"github.com/greg-nicolle/go-microservice/proxy"

	"github.com/greg-nicolle/go-microservice/transport"
)

func proxyingMiddleware(ctx context.Context, instances string, logger *logrus.Logger) ServiceMiddleware {
	// If instances is empty, don't proxy.
	if instances == "" {
		logger.Info("proxy_to", "none")
		return func(next Domain) Domain {
			return next
		}
	}

	var instanceList = proxy.Split(instances)
	logger.Info("proxy_to", fmt.Sprint(instanceList))

	// And finally, return the ServiceMiddleware, implemented by proxymw.
	return func(next Domain) Domain {
		return proxymw{
			logger,
			ctx,
			next,
			proxy.CreatedProxiingEndpoint(ctx, instanceList, transport.DecodeResponse(searchPageNameResponse{}), "/getPageName")}
	}
}

// proxymw implements StringService, forwarding Uppercase requests to the
// provided endpoint, and serving all other (i.e. Count) requests via the
// next StringService.
type proxymw struct {
	logger         *logrus.Logger
	ctx            context.Context
	next           Domain            // Serve most requests via this service...
	SearchPageName endpoint.Endpoint // ...except Uppercase, which gets served by this endpoint
}

func (mw proxymw) searchPageName(s string) (output []string, err error) {
	defer func(begin time.Time) {
		mw.logger.WithFields(logrus.Fields{
			"method": "getPageName",
			"input":  s,
			"output": output,
			"err":    err,
			"took":   time.Since(begin),
		}).Info("proxy !", )
	}(time.Now())

	response, err := mw.SearchPageName(mw.ctx, searchPageNameRequest{S: s})
	if err != nil {
		return nil, err
	}

	resp := response.(searchPageNameResponse)
	if resp.Err != "" {
		return resp.V, errors.New(resp.Err)
	}
	return resp.V, nil
}
