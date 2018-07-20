package main

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func BenchmarkPOW(b *testing.B) {
	POW("hello", 6)
}

func BenchmarkPOWC(b *testing.B) {
	resChan := make(chan string, 1)
	go POWC("hello", 6, resChan)
	go POWC("hello", 6, resChan)
	go POWC("hello", 6, resChan)
	go POWC("hello", 6, resChan)
	fmt.Println(<-resChan)
}

func BenchmarkPOWR(b *testing.B) {
	resChan := make(chan string, 1)
	runtime.GOMAXPROCS(4)
	go POWC("hello", 6, resChan)
	go POWR("hello", 6, 147356834756, resChan)
	go POWR("hello", 6, 8768743567, resChan)
	go POWR("hello", 6, 987682563423, resChan)
	fmt.Println(<-resChan)
}

func TestChannel(t *testing.T) {
	res := make(chan bool, 1)
	go func() {
		for {

			select {
			case <-res:
				fmt.Println("done")
				return
			default:
				fmt.Println("1")
			}
		}
	}()

	time.Sleep(1 * time.Millisecond)
	close(res)
}
