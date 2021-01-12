package echo

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"testing"
	"time"
)

func TestEchoServerUDP(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	serverAddr, err := echoServerUDP(ctx, "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()

	client, err := net.ListenPacket("udp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = client.Close() }()

	msg := []byte("ping")
	_, err = client.WriteTo(msg, serverAddr)
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 1024)
	n, addr, err := client.ReadFrom(buf)
	if err != nil {
		t.Fatal(err)
	}

	if addr.String() != serverAddr.String() {
		t.Fatalf("received reply from %q instead of %q", addr, serverAddr)
	}

	if !bytes.Equal(msg, buf[:n]) {
		t.Errorf("expected reply %q; actual reply %q", msg, buf[:n])
	}
}

func TestDropLocalhostUDPPackets(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	s, err := net.ListenPacket("udp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		defer cancel()

		buf := make([]byte, 1024)
		for {
			n, clientAddr, err := s.ReadFrom(buf) // client to server
			if err != nil {
				return
			}

			_, err = s.WriteTo(buf[:n], clientAddr) // server to client
			if err != nil {
				return
			}
		}
	}()

	server, ok := s.(*net.UDPConn)
	if !ok {
		t.Fatal("not a UDPConn")
	}
	err = server.SetWriteBuffer(2)
	if err != nil {
		t.Fatal(err)
	}

	lAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	client, err := net.ListenUDP("udp", lAddr)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = client.Close() }()

	err = client.SetReadBuffer(2)
	if err != nil {
		t.Fatal(err)
	}

	pings := 50
	for i := 0; i < pings; i++ {
		msg := []byte(fmt.Sprintf("%2d", i))
		_, err = client.WriteTo(msg, s.LocalAddr())
		if err != nil {
			t.Fatal(err)
		}
	}

	err = client.SetReadDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		t.Fatal(err)
	}

	recv := make(chan []byte)
	go func() {
		for {
			buf := make([]byte, 1024)
			n, _, err := client.ReadFrom(buf)
			if err != nil {
				_ = s.Close()
			}
			recv <- buf[:n]
		}
	}()

	replies := 0
OUTER:
	for {
		select {
		case m := <-recv:
			replies++
			t.Logf("%s", m)
		case <-ctx.Done():
			break OUTER
		}
	}

	if replies >= pings {
		t.Fatal("no replies were dropped")
	}
	t.Logf("received %d replies", replies)
}
