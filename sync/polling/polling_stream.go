package polling

import (
	"context"
	"yyyoichi/Distributed-Task-Management-System/pkg/stream"
	"yyyoichi/Distributed-Task-Management-System/sync/api"
)

func generateURL(cxt context.Context, urls []string) <-chan string {
	return stream.Generator[string](cxt, urls...)
}

func getDifferenceUseFunIO(cxt context.Context, urlCh <-chan string, fn func(url string) api.DiffResponse) <-chan api.DiffResponse {
	return stream.FunIO[string, api.DiffResponse](cxt, urlCh, fn)
}
