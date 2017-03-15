package stringModule

import (
  "errors"
  "fmt"
  "golang.org/x/net/context"

  "github.com/go-kit/kit/endpoint"
  "github.com/greg-nicolle/go-microservice/transport"
  "github.com/greg-nicolle/go-microservice/proxy"
  "github.com/Sirupsen/logrus"
)

func proxyingMiddleware(ctx context.Context, instances string, logger logrus.Entry) ServiceMiddleware {
  // If instances is empty, don't proxy.
  if instances == "" {
    logger.Info("proxy_to", "none")
    return func(next StringService) StringService {
      return next
    }
  }

  var instanceList = proxy.Split(instances)
  logger.Info("proxy_to", fmt.Sprint(instanceList))

  // And finally, return the ServiceMiddleware, implemented by proxymw.
  return func(next StringService) StringService {
    return proxymw{ctx,
      next,
      proxy.CreatedProxiingEndpoint(ctx, instanceList, transport.DecodeResponse(uppercaseResponse{}), "/uppercase"),
      proxy.CreatedProxiingEndpoint(ctx, instanceList, transport.DecodeResponse(countResponse{}), "/count")}
  }
}

// proxymw implements StringService, forwarding Uppercase requests to the
// provided endpoint, and serving all other (i.e. Count) requests via the
// next StringService.
type proxymw struct {
  ctx       context.Context
  next      StringService     // Serve most requests via this service...
  uppercase endpoint.Endpoint // ...except Uppercase, which gets served by this endpoint
  count     endpoint.Endpoint // ...except Uppercase, which gets served by this endpoint
}

func (mw proxymw) Count(s string) (int, error) {
  response, err := mw.count(mw.ctx, countRequest{S: s})
  if err != nil {
    return -1, err
  }

  resp := response.(countResponse)
  if resp.Err != "" {
    return resp.V, errors.New(resp.Err)
  }
  return resp.V, nil
}

func (mw proxymw) Uppercase(s string) (string, error) {
  response, err := mw.uppercase(mw.ctx, uppercaseRequest{S: s})
  if err != nil {
    return "", err
  }

  resp := response.(uppercaseResponse)
  if resp.Err != "" {
    return resp.V, errors.New(resp.Err)
  }
  return resp.V, nil
}