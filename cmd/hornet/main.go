package main

import (
	"log"
	"net/http"

	proxy "github.com/DeveloperAromal/hornet/internal/proxy"
	"github.com/DeveloperAromal/hornet/internal/ratelimiter"
	"github.com/DeveloperAromal/hornet/internal/requestguard"

	banner "github.com/DeveloperAromal/hornet/internal/cli"
)

func main() {

	banner.Banner()

	limiter, err := ratelimiter.NewLimiter(ratelimiter.Config{
		RatePerSecond: 10,
		Burst:         20,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer limiter.Close()

	p := proxy.NewProxy()

	rp, err := p.NewReverseProxy("http://localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	handler := ratelimiter.RateLimiter(limiter, nil)(
		requestguard.Limits{
			MaxURLLength: 8 * 1024,
			MaxBodySize:  10 * 1024 * 1024,
		}.RequestGuard(http.HandlerFunc(p.ProxyRequestHandler(rp))),
	)
	http.Handle("/", handler)

	log.Println("hornet listening on :9000")
	if err := http.ListenAndServe(":9000", nil); err != nil {
		log.Fatal(err)
	}
}
