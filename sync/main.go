package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	ENV_STORES     = os.Getenv("STORES")
	ENV_POLLING_MS = os.Getenv("POLLING_MS")
)

func main() {
	if ENV_STORES == "" {
		log.Println("export 'ENV_STORE' is not found")
		return
	}
	if ENV_POLLING_MS == "" {
		log.Println("export 'ENV_POLLING_MS' is not found")
		return
	}

	urls := []string{}
	for _, schema := range strings.Split(ENV_STORES, ",") {
		urls = append(urls, fmt.Sprintf("http://%s", schema))
	}
	pollingMillSecond, err := strconv.Atoi(ENV_POLLING_MS)
	if err != nil {
		log.Println("'ENV_POLLING_MS' is not int")
	}

	cxt := context.WithoutCancel(context.Background())

	log.Println("start sync")
	go func() {
		for {
			select {
			case <-cxt.Done():
				return
			case <-time.After(time.Duration(pollingMillSecond) * time.Millisecond):
				// polling
			}
		}
	}()
}
