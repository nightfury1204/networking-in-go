package main

import (
	"math/rand"
	"fmt"
	"crypto/sha1"
	"encoding/hex"
	"strings"
	"runtime"
)

const charSet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789ABCD"

func main() {
	//resChan := make(chan string, 1)

	// POW("skldfj", 6)
	fmt.Println(runtime.NumCPU())
	fmt.Println(runtime.GOMAXPROCS(-1))

}

func POW(authData string, diffculty int) string{
	prefix1 := ""
	prefix2 := ""
	for i:=0;i<diffculty;i++ {
		prefix1 = prefix1 + "0"
	}

	prefix2 = prefix1 + "0"

	var (
		str string
		checkSumInHex string
	)

	for  {
		str = randomString(12)
		// fmt.Println(st)
		checkSum := sha1.Sum([]byte(authData+str))
		checkSumInHex = hex.EncodeToString(checkSum[:])

		if strings.HasPrefix(checkSumInHex, prefix1) && !strings.HasPrefix(checkSumInHex, prefix2) {
			fmt.Println(str)
			fmt.Println(checkSumInHex)
			return str
		}
	}

}

func randomString(size int) string {
	var buf [512]byte
	n, err := rand.Read(buf[0:size])
	if err!=nil {
		fmt.Errorf("%v",err)
	}

	st:=""

	for i:=0;i<n;i++ {
		st = st+string(charSet[buf[i]&63])
	}

	return st
}

func randomStringR(size int, r *rand.Rand) string {
	var buf [512]byte
	n, err := r.Read(buf[0:size])
	if err!=nil {
		fmt.Errorf("%v",err)
	}

	st:=""

	for i:=0;i<n;i++ {
		st = st+string(charSet[buf[i]&63])
	}

	return st
}


func POWC(authData string, diffculty int, res chan string){
	prefix1 := ""
	prefix2 := ""
	for i:=0;i<diffculty;i++ {
		prefix1 = prefix1 + "0"
	}

	prefix2 = prefix1 + "0"

	var (
		str string
		checkSumInHex string
	)

	for  {
		str = randomString(12)
		// fmt.Println(st)
		checkSum := sha1.Sum([]byte(authData+str))
		checkSumInHex = hex.EncodeToString(checkSum[:])

		if strings.HasPrefix(checkSumInHex, prefix1) && !strings.HasPrefix(checkSumInHex, prefix2) {
			fmt.Println(str)
			fmt.Println(checkSumInHex)
			res <- str
			return
		}
	}

}

func POWR(authData string, diffculty int, seed int64, res chan string){
	fmt.Println("seed", seed)
	prefix1 := ""
	prefix2 := ""
	for i:=0;i<diffculty;i++ {
		prefix1 = prefix1 + "0"
	}

	prefix2 = prefix1 + "0"

	var (
		str string
		checkSumInHex string
	)
	r := rand.New(rand.NewSource(seed))

	for  {
		str = randomStringR(12, r)
		// fmt.Println(st)
		checkSum := sha1.Sum([]byte(authData+str))
		checkSumInHex = hex.EncodeToString(checkSum[:])

		if strings.HasPrefix(checkSumInHex, prefix1) && !strings.HasPrefix(checkSumInHex, prefix2) {
			fmt.Println(str)
			fmt.Println(checkSumInHex)
			res <- str
			return
		}
	}

}