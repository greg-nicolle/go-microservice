package stringModule

import (
  "errors"
  "fmt"
  "context"

  "github.com/go-kit/kit/endpoint"
  "github.com/greg-nicolle/go-microservice/transport"
  "github.com/greg-nicolle/go-microservice/proxy"
  "github.com/Sirupsen/logrus"
  "time"
)

func proxyingMiddleware(ctx context.Context, instances string, logger logrus.Logger) ServiceMiddleware {
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
    return proxymw{
      logger,
      ctx,
      next,
      proxy.CreatedProxiingEndpoint(ctx, instanceList, transport.DecodeResponse(uppercaseResponse{}), "/uppercase"),
      proxy.CreatedProxiingEndpoint(ctx, instanceList, transport.DecodeResponse(countResponse{}), "/count")}
  }
}

// proxymw implements StringService, forwarding Uppercase requests to the
// provided endpoint, and serving all other (i.e. Count) requests via the
// next StringService.
type proxymw struct {
  logger    logrus.Logger
  ctx       context.Context
  next      StringService     // Serve most requests via this service...
  uppercase endpoint.Endpoint // ...except Uppercase, which gets served by this endpoint
  count     endpoint.Endpoint // ...except Uppercase, which gets served by this endpoint
}

func (mw proxymw) Count(s string) (output int, err error) {
  defer func(begin time.Time) {
    mw.logger.WithFields(logrus.Fields{
      "method": "uppercase",
      "input": s,
      "output": output,
      "err": err,
      "took": time.Since(begin),
    }).Info("proxy !", )
  }(time.Now())

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

func (mw proxymw) Uppercase(s string) (output string, err error) {
  defer func(begin time.Time) {
    mw.logger.WithFields(logrus.Fields{
      "method": "uppercase",
      "input": s,
      "output": output,
      "err": err,
      "took": time.Since(begin),
    }).Info("proxy !", )
  }(time.Now())

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