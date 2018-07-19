package main

import (
	"crypto/tls"
	"log"
	"crypto/x509"
	"io/ioutil"
	"net"
	"fmt"
	"io"
)

func main() {
	cert, err := tls.LoadX509KeyPair("../certs/server.crt", "../certs/server.key")
	if err!=nil {
		log.Fatal("falied to load server certificates",err)
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
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs: pool,
	}

	listener, err := tls.Listen("tcp","127.0.0.1:18888",tlsConfig)
	if err!=nil {
		log.Fatal("failed to start tls listener",err)
	}

	log.Println("server listening...")

	for {
		conn, err := listener.Accept()
		if err!= nil {
			log.Println("filed to accept client connection", err)
		} else {
			log.Println("connection accpected:")
			log.Println("local addr:", conn.LocalAddr())
			log.Println("remote addr:", conn.RemoteAddr())

			go handleClient(conn)
		}
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	var buf [512]byte

	for {
		log.Println("trying to read")
		n, err := conn.Read(buf[0:])
		if err == io.EOF {
			fmt.Println("Connection closed")
			return
		}
		if err!=nil {
			fmt.Println(err)
		}

		fmt.Println("from client", string(buf[0:n]))


		_, err2 := conn.Write(buf[0:n])
		if err2!=nil {
			fmt.Println(err2)
			return
		}
	}
}