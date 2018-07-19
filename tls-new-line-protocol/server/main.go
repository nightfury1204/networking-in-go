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

/*
	New line based protocol

	Every message should have '\n' at the end of line
	Server will not response until it got full message
*/

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

	var buf [2]byte
	var store []byte

	for {
		log.Println("trying to read")
		for {
			n, err := conn.Read(buf[0:])
			if err == io.EOF {
				fmt.Println("Connection closed")
				return
			}
			if err!=nil {
				fmt.Println(err)
			}

			for _, val := range buf[0:n] {
				store = append(store, val)
				if val == byte('\n') {
					_, err = conn.Write(store)
					if err!=nil {
						fmt.Println(err)
						return
					}

					store = []byte{}
				}
			}
		}
	}
}