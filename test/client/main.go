package main

import (
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"runtime"
	"strconv"
	"strings"
)

const charSet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789ABCD"

func main() {
	cert, err := tls.LoadX509KeyPair("certs/client.crt", "certs/client.key")
	if err != nil {
		log.Fatal("falied to load client certificates", err)
	}

	pool := x509.NewCertPool()

	caCert, err := ioutil.ReadFile("certs/ca.crt")
	if err != nil {
		log.Fatal("failed to load ca cert", err)
	}

	ok := pool.AppendCertsFromPEM(caCert)
	if !ok {
		log.Fatal("failed to load ca cert in ca cert pool")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
		InsecureSkipVerify:true,
	}

	conn, err := tls.Dial("tcp", "34.245.228.255:3333", tlsConfig)
	if err != nil {
		log.Fatal("failed to create connection", err)
	}

	log.Println("client listening...")

	var buf [512]byte
	var store []byte
	var req string
	authData := ""

	for {
		n, err := conn.Read(buf[0:])
		if err == io.EOF {
			fmt.Println("Connection closed")
			return
		}
		if err != nil {
			fmt.Println(err)
		}

		allReceived := false

		for _, val := range buf[0:n] {
			store = append(store, val)
			if val == byte('\n') {
				req = string(store)
				store = []byte{}
				allReceived = true
				break
			}
		}

		done := false

		if allReceived {
			req = strings.TrimSpace(req)
			args := strings.Split(req, " ")

			fmt.Println("From server: ", req)

			switch args[0] {
			case "HELO":
				_, err = conn.Write([]byte("EHLO\n"))
				if err != nil {
					log.Fatal(err)
				}
				break
			case "ERROR":
				log.Fatal(strings.Join(args[1:], " "))
				break
			case "POW":
				authData = args[1]
				difficulty, err := strconv.ParseInt(args[2], 10, 32)
				if err != nil {
					log.Fatal(err)
				}
				suffix := POW(authData, int(difficulty))
				_, err = conn.Write([]byte(suffix + "\n"))
				if err != nil {
					log.Fatal(err)
				}
				break
			case "END":
				conn.Write([]byte("OK\n"))
				done = true
				break
			case "NAME":
				checksum := sha1.Sum([]byte(authData + args[1]))

				data := hex.EncodeToString(checksum[:])+ " " + "Md. Nure Alam Nahid\n"
				_, err = conn.Write([]byte(data))
				if err != nil {
					log.Fatal(err)
				}
				break
			case "MAILNUM":
				checksum := sha1.Sum([]byte(authData + args[1]))
				data := hex.EncodeToString(checksum[:])+ " " + "1\n"
				_, err = conn.Write([]byte(data))
				if err != nil {
					log.Fatal(err)
				}
				break
			case "MAIL1":
				checksum := sha1.Sum([]byte(authData + args[1]))
				data := hex.EncodeToString(checksum[:])+ " " + "knightnahid@gmail.com\n"
				_, err = conn.Write([]byte(data))
				if err != nil {
					log.Fatal(err)
				}
				break
			case "SKYPE":
				checksum := sha1.Sum([]byte(authData + args[1]))
				data := hex.EncodeToString(checksum[:])+ " " + "nure_alam_nahid\n"
				_, err = conn.Write([]byte(data))
				if err != nil {
					log.Fatal(err)
				}
				break
			case "BIRTHDATE":
				checksum := sha1.Sum([]byte(authData + args[1]))
				data := hex.EncodeToString(checksum[:])+ " " + "21.12.1994\n"
				_, err = conn.Write([]byte(data))
				if err != nil {
					log.Fatal(err)
				}
				break
			case "COUNTRY":
				checksum := sha1.Sum([]byte(authData + args[1]))
				data := hex.EncodeToString(checksum[:])+ " " + "Bangladesh\n"
				_, err = conn.Write([]byte(data))
				if err != nil {
					log.Fatal(err)
				}
				break
			case "ADDRNUM":
				checksum := sha1.Sum([]byte(authData + args[1]))
				data := hex.EncodeToString(checksum[:])+ " " + "4\n"
				_, err = conn.Write([]byte(data))
				if err != nil {
					log.Fatal(err)
				}
				break
			case "ADDRLINE1":
				checksum := sha1.Sum([]byte(authData + args[1]))
				data := hex.EncodeToString(checksum[:])+ " " + "Village: Amdure\n"
				_, err = conn.Write([]byte(data))
				if err != nil {
					log.Fatal(err)
				}
				break
			case "ADDRLINE2":
				checksum := sha1.Sum([]byte(authData + args[1]))
				data := hex.EncodeToString(checksum[:])+ " " + "Post office: Ajgara bazar\n"
				_, err = conn.Write([]byte(data))
				if err != nil {
					log.Fatal(err)
				}
				break
			case "ADDRLINE3":
				checksum := sha1.Sum([]byte(authData + args[1]))
				data := hex.EncodeToString(checksum[:])+ " " + "Upzila: Laksam\n"
				_, err = conn.Write([]byte(data))
				if err != nil {
					log.Fatal(err)
				}
				break
			case "ADDRLINE4":
				checksum := sha1.Sum([]byte(authData + args[1]))
				data := hex.EncodeToString(checksum[:])+ " " + "District: Comilla\n"
				_, err = conn.Write([]byte(data))
				if err != nil {
					log.Fatal(err)
				}
				break

			}

		}

		if done {
			break
		}

	}
	conn.Close()

}

func POW(authData string, difficulty int) string {
	nCpu := runtime.NumCPU()
	runtime.GOMAXPROCS(nCpu)

	res := make(chan string, 4)
	stopCh := make(chan bool)

	for i := 0; i < nCpu; i++ {
		go POWC(authData, difficulty, res,stopCh)
	}

	suffix := <-res
	close(stopCh)
	close(res)
	return suffix
}

func randomString(size int) string {
	var buf [512]byte
	n, err := rand.Read(buf[0:size])
	if err != nil {
		fmt.Errorf("%v", err)
	}

	st := ""

	for i := 0; i < n; i++ {
		st = st + string(charSet[buf[i]&63])
	}

	return st
}

func POWC(authData string, diffculty int, res chan string, stopCh <-chan bool) {
	prefix1 := ""
	prefix2 := ""
	for i := 0; i < diffculty; i++ {
		prefix1 = prefix1 + "0"
	}

	prefix2 = prefix1 + "0"

	var (
		str           string
		checkSumInHex string
	)

	for {
		select {
		case <-stopCh:
			return
		default:
			str = randomString(12)
			// fmt.Println(st)
			checkSum := sha1.Sum([]byte(authData + str))
			checkSumInHex = hex.EncodeToString(checkSum[:])

			if strings.HasPrefix(checkSumInHex, prefix1) && !strings.HasPrefix(checkSumInHex, prefix2) {
				fmt.Println(str)
				fmt.Println(checkSumInHex)
				res <- str
				return
			}
		}
	}
}
