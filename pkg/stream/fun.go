package stream

import (
	"context"
	"runtime"
	"sync"
)

func FunIO[I interface{}, O interface{}](cxt context.Context, inCh <-chan I, fn func(I) O) <-chan O {
	outCh := Out[I, O](cxt, inCh, fn)
	return In[O](cxt, outCh...)
}

func In[T interface{}](cxt context.Context, channels ...<-chan T) <-chan T {
	var wg sync.WaitGroup
	multiplexedCh := make(chan T)
	multiplex := func(c <-chan T) {
		defer wg.Done()
		for i := range c {
			select {
			case <-cxt.Done():
				return
			case multiplexedCh <- i:
			}
		}
	}

	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c)
	}

	go func() {
		wg.Wait()
		close(multiplexedCh)
	}()

	return multiplexedCh
}

func Out[I interface{}, O interface{}](cxt context.Context, inCh <-chan I, fn func(I) O) []<-chan O {
	num := runtime.NumCPU()
	funOut := make([]<-chan O, num)
	for i := 0; i < num; i++ {
		funOut[i] = Line[I, O](cxt, inCh, fn)
	}
	return funOut
}
