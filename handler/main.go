package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	handler(context.Background(), os.Stdin)
	log.Println("start handler")
	select {}
}

var url = fmt.Sprintf("http://%s", os.Getenv("LB_ADDR"))

func handler(cxt context.Context, r io.Reader) {
	log.Printf("[Rq] Get Request\n")
	sc := bufio.NewScanner(r)
	go func() {
		for {
			fmt.Printf("Create TODO list... ")
			sc.Scan()
			task := sc.Text()
			body := []byte(fmt.Sprintf(`{"task":"%s"}`, task))
			resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
			if err != nil {
				log.Println(err)
				continue
			}
			fmt.Printf("TODO is Created!\n")
			resp.Body.Close()
		}
	}()
}
