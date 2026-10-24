package main

import (
	"log"
	"net/http"

	"github.com/DeveloperAromal/hornet/internal/middleware/ratelimiter"
	proxy "github.com/DeveloperAromal/hornet/internal/proxy"
)

func main() {
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
		http.HandlerFunc(p.ProxyRequestHandler(rp)),
	)
	http.Handle("/", handler)

	log.Println("hornet listening on :9000")
	if err := http.ListenAndServe(":9000", nil); err != nil {
		log.Fatal(err)
	}
}
