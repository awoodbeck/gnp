package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/awoodbeck/gnp/ch13/instrumentation/metrics"
)

var (
	metricsAddr = flag.String("metrics", "127.0.0.1:8081",
		"metrics listen address")
	webAddr = flag.String("web", "127.0.0.1:8082", "web listen address")
)

func helloHandler(w http.ResponseWriter, _ *http.Request) {
	metrics.Requests.Add(1)
	defer func(start time.Time) {
		metrics.RequestDuration.Observe(time.Since(start).Seconds())
	}(time.Now())

	_, err := w.Write([]byte("Hello!"))
	if err != nil {
		metrics.WriteErrors.Add(1)
	}
}

func newHTTPServer(addr string, mux http.Handler,
	stateFunc func(net.Conn, http.ConnState)) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		IdleTimeout:       time.Minute,
		ReadHeaderTimeout: 30 * time.Second,
		ConnState:         stateFunc,
	}

	go func() { log.Fatal(srv.Serve(l)) }()

	return nil
}

func connStateMetrics(_ net.Conn, state http.ConnState) {
	switch state {
	case http.StateNew:
		metrics.OpenConnections.Add(1)
	case http.StateClosed:
		metrics.OpenConnections.Add(-1)
	}
}

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	mux := http.NewServeMux()
	mux.Handle("/metrics/", promhttp.Handler())
	if err := newHTTPServer(*metricsAddr, mux, nil); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Metrics listening on %q ...\n", *metricsAddr)

	if err := newHTTPServer(*webAddr, http.HandlerFunc(helloHandler),
		connStateMetrics); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Web listening on %q ...\n\n", *webAddr)

	clients := 500
	gets := 100
	wg := new(sync.WaitGroup)

	fmt.Printf("Spawning %d connections to make %d requests each ...",
		clients, gets)
	for i := 0; i < clients; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			c := &http.Client{
				Transport: http.DefaultTransport.(*http.Transport).Clone(),
			}

			for j := 0; j < gets; j++ {
				resp, err := c.Get(fmt.Sprintf("http://%s/", *webAddr))
				if err != nil {
					log.Fatal(err)
				}
				_, _ = io.Copy(ioutil.Discard, resp.Body)
				_ = resp.Body.Close()
			}
		}()
	}
	wg.Wait()
	fmt.Print(" done.\n\n")

	resp, err := http.Get(fmt.Sprintf("http://%s/metrics", *metricsAddr))
	if err != nil {
		log.Fatal(err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()

	metricsPrefix := fmt.Sprintf("%s_%s", *metrics.Namespace,
		*metrics.Subsystem)
	fmt.Println("Current Metrics:")
	for _, line := range bytes.Split(b, []byte("\n")) {
		if bytes.HasPrefix(line, []byte(metricsPrefix)) {
			fmt.Printf("%s\n", line)
		}
	}
}
