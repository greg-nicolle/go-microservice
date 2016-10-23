package proxy

import (
  "strings"
  "net/url"
  "github.com/go-kit/kit/circuitbreaker"
  "github.com/sony/gobreaker"
  "github.com/go-kit/kit/sd/lb"
  "github.com/go-kit/kit/endpoint"
  "time"
  "github.com/go-kit/kit/sd"

  jujuratelimit "github.com/juju/ratelimit"
  "golang.org/x/net/context"

  "github.com/go-kit/kit/ratelimit"
  httptransport "github.com/go-kit/kit/transport/http"
  "github.com/greg-nicolle/go-microservice/transport"
)

func Split(s string) []string {
  a := strings.Split(s, ",")
  for i := range a {
    a[i] = strings.TrimSpace(a[i])
  }
  return a
}

func CreadteProxiingEndpoint(instances []string, ctx context.Context, decodeFunc httptransport.DecodeResponseFunc, path string) endpoint.Endpoint {
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
