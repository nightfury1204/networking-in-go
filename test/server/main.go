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

	caCert, err := ioutil.ReadFile("certs/ca.crt")
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
	write(conn, "HELO\n")
	read(conn)
	write(conn, "POW auth 4\n")
	read(conn)
	write(conn, "NAME auth\n")
	read(conn)
	write(conn, "MAILNUM auth\n")
	read(conn)
	write(conn, "MAIL1 auth\n")
	read(conn)
	write(conn, "SKYPE auth\n")
	read(conn)
	write(conn, "BIRTHDATE auth\n")
	read(conn)
	write(conn, "COUNTRY auth\n")
	read(conn)
	write(conn, "ADDRNUM auth\n")
	read(conn)
	write(conn, "ADDRLINE1 auth\n")
	read(conn)
	write(conn, "ADDRLINE2 auth\n")
	read(conn)
	write(conn, "ADDRLINE3 auth\n")
	read(conn)
	write(conn, "ADDRLINE4 auth\n")
	read(conn)
	write(conn, "END\n")
	read(conn)

	return
}

func write(conn net.Conn, data string) {
	_, err := conn.Write([]byte(data))
	if err!=nil {
		fmt.Errorf("%v\n",err)
	}
}

func read(conn net.Conn) {
	var buf [512]byte
	var store []byte
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
				fmt.Println("From client: ",string(store))
				return
			}
		}
	}
}