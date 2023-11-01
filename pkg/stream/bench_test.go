package stream

import (
	"context"
	"math/rand"
	"testing"
)

// 効率の悪いアルゴリズムを利用した素数判定機
func findPrimer(number int) bool {
	if number <= 1 {
		return false
	}

	for i := 2; i < number; i++ {
		if number%i == 0 {
			return false
		}
	}
	return true
}

func genRand(n int) []int {
	rands := make([]int, n)
	for i := 0; i < n; i++ {
		rands = append(rands, rand.Intn(5*1000*1000))
	}
	return rands
}

func Benchmark_FunInOut(b *testing.B) {
	rands := genRand(b.N)
	b.ResetTimer()
	cxt := context.Background()
	randStream := Generator[int](cxt, rands...)
	outStream := Out[int, bool](cxt, randStream, findPrimer)
	inStream := In[bool](cxt, outStream...)
	for range inStream {
	}
}

func Benchmark_FunIO(b *testing.B) {
	rands := genRand(b.N)
	b.ResetTimer()
	cxt := context.Background()
	randStream := Generator[int](cxt, rands...)
	inStream := FunIO[int, bool](cxt, randStream, findPrimer)
	for range inStream {
	}
}

func Benchmark_NoUseFUNIO(b *testing.B) {
	rands := genRand(b.N)
	b.ResetTimer()
	cxt := context.Background()
	randStream := Generator[int](cxt, rands...)
	findRandStream := Line[int, bool](cxt, randStream, findPrimer)
	for range findRandStream {
	}
}
