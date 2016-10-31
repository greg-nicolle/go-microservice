package stringModule

import (
  "errors"
  "golang.org/x/net/context"

  "github.com/go-kit/kit/endpoint"
  "github.com/go-kit/kit/log"
  "github.com/greg-nicolle/go-microservice/transport"
  "github.com/greg-nicolle/go-microservice/proxy"
)

func proxyingMiddleware(ctx context.Context, instances []string, logger log.Logger) ServiceMiddleware {
  // If instances is empty, don't proxy.
  if len(instances) == 0 {
    logger.Log("proxy_to", "none")
    return func(next StringService) StringService {
      return next
    }
  }

  logger.Log("proxy_to", instances[0])

  // And finally, return the ServiceMiddleware, implemented by proxymw.
  return func(next StringService) StringService {
    return proxymw{ctx: ctx,
      next: next,
      uppercaseProxy: proxy.CreatedProxiingEndpoint(instances, transport.DecodeResponse(uppercaseResponse{}), "/uppercase"),
      countProxy: proxy.CreatedProxiingEndpoint(instances, transport.DecodeResponse(countResponse{}), "/count")}
  }
}

// proxymw implements StringService, forwarding Uppercase requests to the
// provided endpoint, and serving all other (i.e. Count) requests via the
// next StringService.
type proxymw struct {
  ctx            context.Context
  next           StringService     // Serve most requests via this service...
  uppercaseProxy endpoint.Endpoint // ...except Uppercase, which gets served by this endpoint
  countProxy     endpoint.Endpoint // ...except Uppercase, which gets served by this endpoint
}

func (mw proxymw) count(s string) (int, error) {
  response, err := mw.countProxy(mw.ctx, countRequest{S: s})
  if err != nil {
    return -1, err
  }

  resp := response.(countResponse)
  if resp.Err != "" {
    return resp.V, errors.New(resp.Err)
  }
  return resp.V, nil
}

func (mw proxymw) uppercase(s string) (string, error) {
  response, err := mw.uppercaseProxy(mw.ctx, uppercaseRequest{S: s})
  if err != nil {
    return "", err
  }

  resp := response.(uppercaseResponse)
  if resp.Err != "" {
    return resp.V, errors.New(resp.Err)
  }
  return resp.V, nil
}