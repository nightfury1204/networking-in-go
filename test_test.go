package main

import (
	"fmt"
	"runtime"
	"testing"
	"time"
	"crypto/sha1"
	"encoding/hex"
	"strings"
)

func BenchmarkPOW(b *testing.B) {
	POW("hello", 6)
}

func BenchmarkPOWC(b *testing.B) {
	resChan := make(chan string, 1)
	prefix := "FPMREQeIgVElLaklbNbSPrAktwhgfjYrDVBFZdiAzNNneAodlZdbQCkaDnDNrLtw"
	difficulty := 4
	go POWC(prefix, difficulty, resChan)
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

func TestR(t *testing.T) {
	resChan := make(chan string, 1)
	prefix := "FPMREQeIgVElLaklbNbSPrAktwhgfjYrDVBFZdiAzNNneAodlZdbQCkaDnDNrLtw"
	//prefix = prefix + "nG2NojGWtgXMaq54"
	//prefix = prefix + "ihV8NUr5Dg1FGa"
	difficulty := 6
	go POWC(prefix, difficulty, resChan)
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

func BenchmarkRandStringBytesMaskImpr(b *testing.B) {
	for i:=1;i<=20000000;i++ {
		RandStringBytesMaskImpr(60)
	}
	// 13387575883 ns
}

func BenchmarkRandomString(b *testing.B) {
	for i:=1;i<=20000000;i++ {
		randomString(60)
	}
}

func BenchmarkSha1(b *testing.B) {
	st := RandStringBytesMaskImpr(60)
	prefix := "000000000"
	for i:=1;i<=20000000;i++ {
		ch := sha1.Sum([]byte(st))
		hex.EncodeToString(ch[:])
		if strings.HasPrefix(hex.EncodeToString(ch[:]),prefix) {

		}
	}
	// 13387575883 ns
}

func BenchmarkSha11(b *testing.B) {
	st := RandStringBytesMaskImpr(60)
	for i:=1;i<=20000000;i++ {
		ch := sha1.Sum([]byte(st))
		if ch[0]==0 && ch[1]==0 && ch[2]==0 && ch[3]==0 && (ch[4]&0xf0)==0 {

		}
	}
	// 13387575883 ns
}

func TestPOWC2(t *testing.T) {
	resChan := make(chan string, 1)
	prefix := ""
	difficulty := 9
	POWC(prefix, difficulty, resChan)
	fmt.Println(<-resChan)
}


/*
zWjTAPkt1xJofuV8VCAQDM
0011e4da72b0a4b7b0655a49e7daa4b61d40b931
xRdKiPB5B4PeWmNehbLTdx
00065def9af3083c55ea4d466303ff6038eb59e4
bJzEQgE4jBHapBVMYYNdzc
00b693f607e5b6a8f342167d58f6fef7c8e0334e
HuK7pggEYPOs7vBNYCr9KN
00cfddaf740b7fc37a1fcf39a2e89a70f98827c5
*/
