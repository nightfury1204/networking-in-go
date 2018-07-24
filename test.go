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
	var (
		str string
		// checkSumInHex string
	)
	sz := 23
	for  {
		str = RandStringBytesMaskImpr(int(sz))
		// fmt.Println(st)
		checkSum := sha1.Sum([]byte(authData+str))

		if checkSum[0] == 0 && checkSum[1]==0 && checkSum[2]==0 && (checkSum[3])==0 && (checkSum[4]&0xf0)==0 {
			fmt.Println(str)
			fmt.Println(hex.EncodeToString(checkSum[:]))
			res <- str
			return
		}

		//sz = rand.Uint32() & 0xff
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

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

//RandStringBytesMaskImpr - generate random string using masking improved
func RandStringBytesMaskImpr(n int) string {
	b := make([]byte, n)
	l := len(letterBytes)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < l {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
