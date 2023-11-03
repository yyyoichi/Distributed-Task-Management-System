package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/yyyoichi/Distributed-Task-Management-System/store/handler"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("port is not found")
	}
	addr := fmt.Sprintf(":%s", port)

	sh := handler.NewStoreHandler()

	log.Println("start key-value store")
	http.HandleFunc("/differences", sh.DifferencesHandler)
	http.HandleFunc("/sync", sh.SynchronizeHandler)
	http.HandleFunc("/", sh.CommandsHandler)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Println(err)
	}
}
