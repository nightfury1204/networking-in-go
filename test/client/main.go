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
	"strings"
	"crypto/sha1"
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

	for {

		var buf [512]byte
		var store []byte
		var resp string
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
					resp =  string(store)
					store = []byte{}
					allReceived = true
				}
			}

			if allReceived {
				resp = strings.TrimSpace(resp)
				args := strings.Split(resp, " ")

				switch args[0] {
				case "HELO":
					_, err = conn.Write([]byte("EHLO\n"))
					if err!=nil {
						log.Fatal(err)
					}
					break
				case "ERROR":
					log.Fatal(strings.Join(args[1:], " "))
					break
				case "POW":
					// implement later
					break
				case "END":
					conn.Write([]byte("OK\n"))
					break
				case "NAME":
					checksum := sha1.Sum([]byte(authData+ args[1]))
					data := string(checksum[:])+" "+"Md. Nure Alam Nahid\n"
					_, err = conn.Write([]byte(data))
					if err!=nil {
						log.Fatal(err)
					}
					break
				case "MAILNUM":
					checksum := sha1.Sum([]byte(authData+ args[1]))
					data := string(checksum[:])+" "+"1\n"
					_, err = conn.Write([]byte(data))
					if err!=nil {
						log.Fatal(err)
					}
					break
				case "MAIL1":
					checksum := sha1.Sum([]byte(authData+ args[1]))
					data := string(checksum[:])+" "+"knightnahid@gmail.com\n"
					_, err = conn.Write([]byte(data))
					if err!=nil {
						log.Fatal(err)
					}
					break
				case "SKYPE":
					checksum := sha1.Sum([]byte(authData+ args[1]))
					data := string(checksum[:])+" "+"nure_alam_nahid\n"
					_, err = conn.Write([]byte(data))
					if err!=nil {
						log.Fatal(err)
					}
					break
				case "BIRTHDATE":
					checksum := sha1.Sum([]byte(authData+ args[1]))
					data := string(checksum[:])+" "+"21.12.1994\n"
					_, err = conn.Write([]byte(data))
					if err!=nil {
						log.Fatal(err)
					}
					break
				case "COUNTRY":
					checksum := sha1.Sum([]byte(authData+ args[1]))
					data := string(checksum[:])+" "+"Bangladesh\n"
					_, err = conn.Write([]byte(data))
					if err!=nil {
						log.Fatal(err)
					}
					break
				case "ADDRNUM":
					checksum := sha1.Sum([]byte(authData+ args[1]))
					data := string(checksum[:])+" "+"1\n"
					_, err = conn.Write([]byte(data))
					if err!=nil {
						log.Fatal(err)
					}
					break
				case "ADDRLINE1":
					checksum := sha1.Sum([]byte(authData+ args[1]))
					data := string(checksum[:])+" "+"Bangladesh\n"
					_, err = conn.Write([]byte(data))
					if err!=nil {
						log.Fatal(err)
					}
					break

				}

			}
		}

	}

	conn.Close()

}
