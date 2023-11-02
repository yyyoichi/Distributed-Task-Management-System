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

// int型のmapをループしてO型のチャネルを送信する関数
//
// [fn]は引数で受け取ったmapの[key]と[val]をO型に変換する関数が必要。
func GeneratorWithMapIntKey[V interface{}, O interface{}](cxt context.Context, m map[int]V, fn func(k int, v V) O) <-chan O {
	ch := make(chan O, len(m))
	go func() {
		defer close(ch)
		for key, val := range m {
			select {
			case <-cxt.Done():
				return
			case ch <- fn(key, val):
			}
		}
	}()

	return ch
}
