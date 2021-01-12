// +build darwin linux

package echo

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"testing"
)

func TestEchoServerUnixDatagram(t *testing.T) {
	dir, err := ioutil.TempDir("", "echo_unixgram")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if rErr := os.RemoveAll(dir); rErr != nil {
			t.Error(rErr)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	sSocket := filepath.Join(dir, fmt.Sprintf("s%d.sock", os.Getpid()))
	serverAddr, err := datagramEchoServer(ctx, "unixgram", sSocket)
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()

	err = os.Chmod(sSocket, os.ModeSocket|0622)
	if err != nil {
		t.Fatal(err)
	}

	cSocket := filepath.Join(dir, fmt.Sprintf("c%d.sock", os.Getpid()))
	client, err := net.ListenPacket("unixgram", cSocket)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = client.Close() }()

	err = os.Chmod(cSocket, os.ModeSocket|0622)
	if err != nil {
		t.Fatal(err)
	}

	msg := []byte("ping")
	for i := 0; i < 3; i++ { // write 3 "ping" messages
		_, err = client.WriteTo(msg, serverAddr)
		if err != nil {
			t.Fatal(err)
		}
	}

	buf := make([]byte, 1024)

	for i := 0; i < 3; i++ { // read 3 "ping" messages
		n, addr, err := client.ReadFrom(buf)
		if err != nil {
			t.Fatal(err)
		}

		if addr.String() != serverAddr.String() {
			t.Fatalf("received reply from %q instead of %q", addr,
				serverAddr)
		}

		if !bytes.Equal(msg, buf[:n]) {
			t.Fatalf("expected reply %q; actual reply %q", msg, buf[:n])
		}
	}
}

func BenchmarkEchoServerUnixDatagram(b *testing.B) {
	dir, err := ioutil.TempDir("", "echo_unixgram_bench")
	if err != nil {
		b.Fatal(err)
	}
	defer func() {
		if rErr := os.RemoveAll(dir); rErr != nil {
			b.Error(rErr)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	socket := filepath.Join(dir, fmt.Sprintf("s%d.sock", os.Getpid()))
	serverAddr, err := datagramEchoServer(ctx, "unixgram", socket)
	if err != nil {
		b.Fatal(err)
	}
	defer cancel()

	cSocket := filepath.Join(dir, fmt.Sprintf("c%d.sock", os.Getpid()))
	client, err := net.ListenPacket("unixgram", cSocket)
	if err != nil {
		b.Fatal(err)
	}
	defer func() { _ = client.Close() }()

	msg := []byte("ping")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = client.WriteTo(msg, serverAddr)
		if err != nil {
			b.Fatal(err)
		}

		buf := make([]byte, 1024)
		n, addr, err := client.ReadFrom(buf)
		if err != nil {
			b.Fatal(err)
		}

		if addr.String() != serverAddr.String() {
			b.Fatalf("received reply from %q instead of %q", addr,
				serverAddr)
		}

		if !bytes.Equal(msg, buf[:n]) {
			b.Fatalf("expected reply %q; actual reply %q", msg, buf[:n])
		}
	}
}
