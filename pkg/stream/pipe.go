package stream

import "context"

func Line[I interface{}, O interface{}](cxt context.Context, inCh <-chan I, fn func(I) O) <-chan O {
	outCh := make(chan O)
	go func() {
		defer close(outCh)
		for {
			select {
			case <-cxt.Done():
				return
			case in, ok := <-inCh:
				if !ok {
					return
				}
				select {
				case <-cxt.Done():
					return
				case outCh <- fn(in):
				}
			}
		}
	}()
	return outCh
}

// 一つのチャネルから複数のチャネルを送信する。
//
// [fn]には第一引数にI型のデータから、第二引数の`producer`関数に複数のO型を返す関数を実装することで、
// デマルチプレクサとして機能するようになる。
func Demulti[I interface{}, O interface{}](cxt context.Context, inCh <-chan I, fn func(I, func(O))) <-chan O {
	outCh := make(chan O)
	go func() {
		defer close(outCh)
		producer := func(o O) {
			select {
			case <-cxt.Done():
			default:
				outCh <- o
			}
		}
		for {
			select {
			case <-cxt.Done():
				return
			case in, ok := <-inCh:
				if !ok {
					return
				}
				fn(in, producer)
			}
		}
	}()
	return outCh
}
