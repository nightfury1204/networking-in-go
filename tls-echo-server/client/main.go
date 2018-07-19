package main

import (
	"crypto/tls"
	"log"
	"crypto/x509"
	"io/ioutil"
	"fmt"
	"time"
)

func main() {
	cert, err := tls.LoadX509KeyPair("../certs/client.crt", "../certs/client.key")
	if err!=nil {
		log.Fatal("falied to load client certificates",err)
	}

	pool := x509.NewCertPool()

	caCert, err := ioutil.ReadFile("../certs/ca.crt")
	if err!=nil {
		log.Fatal("failed to load ca cert",err)
	}

	ok := pool.AppendCertsFromPEM(caCert)
	if !ok {
		log.Fatal("failed to load ca cert in ca cert pool")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs: pool,
		ServerName: "127.0.0.1",
	}



	conn, err := tls.Dial("tcp","127.0.0.1:18888",tlsConfig)
	if err!=nil {
		log.Fatal("failed to create connection",err)
	}

	log.Println("client listening...")



	for {

		var st string

		fmt.Scanln(&st)

		if st == "exit" {
			break
		}

		if st == "" {
			continue
		}

		conn.Write([]byte(st))

		var buf [512]byte

		conn.SetReadDeadline(time.Now().Add(5*time.Second))

		n, err := conn.Read(buf[0:])
		if err!=nil {
			fmt.Errorf("%v\n",err)
		}

		fmt.Println("server response:", string(buf[0:n]))
	}

	conn.Close()


}

