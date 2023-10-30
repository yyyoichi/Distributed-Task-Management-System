package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync"
)

func main() {
	var proxy = NewProxy()

	// export server ports
	ports := os.Getenv("EXPORTS")
	if ports == "" {
		log.Println("export ports is not found")
		return
	}
	for _, schema := range strings.Split(ports, ",") {
		proxy.backend = append(proxy.backend, fmt.Sprintf("http://%s", schema))
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Println("port is not found")
		return
	}
	addr := fmt.Sprintf(":%s", port)

	http.HandleFunc("/", proxy.handler)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
	log.Println("start loadbrancer")
}

func NewProxy() TProxy {
	return TProxy{
		[]string{},
		0,
		&sync.Mutex{},
	}
}

type TProxy struct {
	backend []string
	current int
	mu      *sync.Mutex
}

func (p *TProxy) handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[LB] Get Request\n")
	p.mu.Lock()
	backendUrl := p.backend[p.current]
	p.current = (p.current + 1) % len(p.backend)
	p.mu.Unlock()

	// Parse the backend URL
	backend, err := url.Parse(backendUrl)
	if err != nil {
		http.Error(w, "Bad gateway", http.StatusBadGateway)
		return
	}

	// Create a reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(backend)

	// Serve the request using the reverse proxy
	proxy.ServeHTTP(w, r)
}
