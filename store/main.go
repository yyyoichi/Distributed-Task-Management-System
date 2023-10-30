package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"yyyoichi/Distributed-Task-Management-System/store/store"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("port is not found")
	}
	addr := fmt.Sprintf(":%s", port)

	var store = store.NewStore()

	log.Println("start key-value store")
	http.HandleFunc("/", store.Handler)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Println(err)
	}
}
