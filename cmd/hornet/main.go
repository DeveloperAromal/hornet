package main

import (
	"log"
	"net/http"
	"os"

	"github.com/DeveloperAromal/hornet/internal/auth"
	banner "github.com/DeveloperAromal/hornet/internal/cli"
	proxy "github.com/DeveloperAromal/hornet/internal/proxy"
	"github.com/DeveloperAromal/hornet/internal/ratelimiter"
	"github.com/DeveloperAromal/hornet/internal/requestguard"
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

	verifier, err := auth.NewJwtVerifier(
		os.Getenv("HORNET_JWT_SECRET"),
		"hornet-auth",
	)
	if err != nil {
		log.Fatal(err)
	}

	p := proxy.NewProxy()

	rp, err := p.NewReverseProxy("http://localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	var handler http.Handler = http.HandlerFunc(
		p.ProxyRequestHandler(rp),
	)

	handler = verifier.Authentication(handler)

	handler = (requestguard.Limits{
		MaxURLLength: 8 * 1024,
		MaxBodySize:  10 * 1024 * 1024,
	}).RequestGuard(handler)

	handler = ratelimiter.RateLimiter(limiter, nil)(handler)

	http.Handle("/", handler)

	log.Println("HORNET listening on :9000")
	if err := http.ListenAndServe(":9000", nil); err != nil {
		log.Fatal(err)
	}
}
