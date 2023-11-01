package stream

import "context"

func Generator[T interface{}](cxt context.Context, seeds ...T) <-chan T {
	ch := make(chan T, len(seeds))
	go func() {
		defer close(ch)
		for _, s := range seeds {
			select {
			case <-cxt.Done():
				return
			case ch <- s:
			}
		}
	}()
	return ch
}

func GeneratorWithFn[I interface{}, O interface{}](cxt context.Context, fn func(I) O, seeds ...I) <-chan O {
	ch := make(chan O, len(seeds))
	go func() {
		defer close(ch)
		for _, s := range seeds {
			select {
			case <-cxt.Done():
				return
			case ch <- fn(s):
			}
		}
	}()
	return ch
}
