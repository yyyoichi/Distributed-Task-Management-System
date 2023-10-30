package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("port is not found")
	}
	addr := fmt.Sprintf(":%s", port)

	sh := NewStoreHandler()

	log.Println("start key-value store")
	http.HandleFunc("/", sh.handler)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Println(err)
	}
}
