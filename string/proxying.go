package stringModule

import (
  "errors"
  "fmt"
  "net/url"
  "strings"
  "time"

  jujuratelimit "github.com/juju/ratelimit"
  "github.com/sony/gobreaker"
  "golang.org/x/net/context"

  "github.com/go-kit/kit/circuitbreaker"
  "github.com/go-kit/kit/endpoint"
  "github.com/go-kit/kit/log"
  "github.com/go-kit/kit/ratelimit"
  "github.com/go-kit/kit/sd"
  "github.com/go-kit/kit/sd/lb"
  httptransport "github.com/go-kit/kit/transport/http"
  "github.com/greg-nicolle/kit-test/transport"
)

func proxyingMiddleware(instances string, ctx context.Context, logger log.Logger) ServiceMiddleware {
  // If instances is empty, don't proxy.
  if instances == "" {
    logger.Log("proxy_to", "none")
    return func(next StringService) StringService {
      return next
    }
  }

  var instanceList = split(instances)
  logger.Log("proxy_to", fmt.Sprint(instanceList))

  // And finally, return the ServiceMiddleware, implemented by proxymw.
  return func(next StringService) StringService {
    return proxymw{ctx,
      next,
      creadteProxiingEndpoint(instanceList, ctx, transport.DecodeResponse(uppercaseResponse{}), "/uppercase"),
      creadteProxiingEndpoint(instanceList, ctx, transport.DecodeResponse(countResponse{}), "/count")}
  }
}

func creadteProxiingEndpoint(instances []string, ctx context.Context, decodeFunc httptransport.DecodeResponseFunc, path string) endpoint.Endpoint {
  var (
    qps = 100                         // beyond which we will return an error
    maxAttempts = 3                   // per request, before giving up
    maxTime = 250 * time.Millisecond  // wallclock time, before giving up
  )
  var subscriberUppercase sd.FixedSubscriber
  for _, instance := range instances {
    var e endpoint.Endpoint

    if !strings.HasPrefix(instance, "http") {
      instance = "http://" + instance
    }
    u, err := url.Parse(instance)
    if err != nil {
      panic(err)
    }
    if u.Path == "" {
      u.Path = path
    }
    e = httptransport.NewClient(
      "GET",
      u,
      transport.EncodeRequest,
      decodeFunc,
    ).Endpoint()

    e = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(e)
    e = ratelimit.NewTokenBucketLimiter(jujuratelimit.NewBucketWithRate(float64(qps), int64(qps)))(e)
    subscriberUppercase = append(subscriberUppercase, e)
  }
  balancerUppercase := lb.NewRoundRobin(subscriberUppercase)
  retry := lb.Retry(maxAttempts, maxTime, balancerUppercase)

  return retry
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

func split(s string) []string {
  a := strings.Split(s, ",")
  for i := range a {
    a[i] = strings.TrimSpace(a[i])
  }
  return a
}
