package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"bufio"
)

func main() {
	cert, err := tls.LoadX509KeyPair("../certs/client.crt", "../certs/client.key")
	if err != nil {
		log.Fatal("falied to load client certificates", err)
	}

	pool := x509.NewCertPool()

	caCert, err := ioutil.ReadFile("../certs/ca.crt")
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
		ServerName:   "127.0.0.1",
	}

	conn, err := tls.Dial("tcp", "127.0.0.1:18888", tlsConfig)
	if err != nil {
		log.Fatal("failed to create connection", err)
	}

	log.Println("client listening...")
	scanner := bufio.NewScanner(os.Stdin)


	for {

		var st string

		scanner.Scan()
		st = scanner.Text()

		fmt.Println("client:",st)

		if st == "exit" {
			break
		}

		if st == "" {
			continue
		}

		conn.Write([]byte(st + "\n"))

		var buf [2]byte
		var store []byte

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
					fmt.Print("server response:", string(store))
					store = []byte{}
					allReceived = true
				}
			}

			if allReceived {
				break
			}
		}

	}

	conn.Close()

}
