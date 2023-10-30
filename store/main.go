package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"yyyoichi/Distributed-Task-Management-System/database/store"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("port is not found")
	}
	addr := fmt.Sprintf(":%s", port)

	var store = store.NewStore()

	http.HandleFunc("/", store.Handler)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Println(err)
	}
	log.Println("start database")
}
