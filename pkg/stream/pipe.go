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

func Demulti[I interface{}, O interface{}](cxt context.Context, inCh <-chan I, fn func(I, func(O))) <-chan O {
	outCh := make(chan O)
	go func() {
		defer close(outCh)
		send := func(o O) {
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
				fn(in, send)
			}
		}
	}()
	return outCh
}
