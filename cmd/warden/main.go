package main

import (
	"log"
	"net/http"

	proxy "github.com/DeveloperAromal/warden/internal/proxy"
)

func main() {
	p := proxy.NewProxy()

	rp, err := p.NewReverseProxy("http://localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", p.ProxyRequestHandler(rp))

	log.Println("Warden listening on :9000")
	if err := http.ListenAndServe(":9000", nil); err != nil {
		log.Fatal(err)
	}
}
