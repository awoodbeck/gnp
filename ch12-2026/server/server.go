package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/awoodbeck/gnp/ch12-2026/housework/v1"
)

var addr, certFn, keyFn string

func init() {
	flag.StringVar(&addr, "address", "localhost:34443",
		"listen address")
	flag.StringVar(&certFn, "cert", "cert.pem", "certificate file")
	flag.StringVar(&keyFn, "key", "key.pem", "private key file")
}

func main() {
	flag.Parse()

	cert, err := tls.LoadX509KeyPair(certFn, keyFn)
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer(
		grpc.Creds(credentials.NewTLS(
			&tls.Config{
				Certificates:     []tls.Certificate{cert},
				CurvePreferences: []tls.CurveID{tls.CurveP256},
				MinVersion:       tls.VersionTLS12,
			},
		)),
	)
	rosie := new(Rosie)
	housework.RegisterRobotMaidServer(server, rosie)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Listening for TLS connections on %s ...", addr)
	log.Fatal(server.Serve(listener))
}
