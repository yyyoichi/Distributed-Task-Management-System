package main

import (
	"bufio"
	"bytes"
	"context"
	"errors"
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
			input := sc.Text()
			output := request(cxt, input)
			log.Println(output)
		}
	}()
}

var (
	ErrNetwork        = errors.New("cannnot 'post' request")
	ErrResponseDecode = errors.New("cannot decode response body")
)

func request(cxt context.Context, input string) string {
	body := []byte(fmt.Sprintf(`{"task":"%s"}`, input))
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Sprintf("Error\n%s:%s", ErrNetwork, err)
	}
	defer resp.Body.Close()

	output, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("Error\n%s:%s", ErrResponseDecode, err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("Error:\n%s", string(output))
	}
	return string(output)
}
