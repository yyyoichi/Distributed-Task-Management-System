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

	var store = NewStore()

	http.HandleFunc("/", store.handler)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Println(err)
	}
	log.Println("start database")
}
