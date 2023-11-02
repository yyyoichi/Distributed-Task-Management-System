package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"yyyoichi/Distributed-Task-Management-System/sync/polling"
)

var (
	ENV_STORES     = os.Getenv("STORES")
	ENV_POLLING_MS = os.Getenv("POLLING_MS")
)

func main() {
	// 環境変数バリデーション
	if ENV_STORES == "" {
		log.Println("export 'ENV_STORE' is not found")
		return
	}
	if ENV_POLLING_MS == "" {
		log.Println("export 'ENV_POLLING_MS' is not found")
		return
	}

	// データノードのリクエストURK
	urls := []string{}
	for _, schema := range strings.Split(ENV_STORES, ",") {
		urls = append(urls, fmt.Sprintf("http://%s", schema))
	}
	// 同期機構のポーリング間隔
	pollingMillSecond, err := strconv.Atoi(ENV_POLLING_MS)
	if err != nil {
		log.Println("'ENV_POLLING_MS' is not int")
	}

	// 処理開始
	cxt := context.WithoutCancel(context.Background())
	pm := polling.NewPollingManager(urls)

	log.Println("start sync")
	go func() {
		for {
			select {
			case <-cxt.Done():
				return
			case <-time.After(time.Duration(pollingMillSecond) * time.Millisecond):
				// polling
				pm.Polling(cxt)
			}
		}
	}()
	select {}
}
