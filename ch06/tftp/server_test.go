package tftp

import (
	"bytes"
	"io"
	"io/ioutil"
	"net"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	t.Parallel()

	p1, err := ioutil.ReadFile("./tftp/payload.svg")
	if err != nil {
		t.Fatal(err)
	}

	conn, err := net.ListenPacket("udp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	done := make(chan struct{})
	s := Server{Payload: p1}

	go func() {
		_ = s.Serve(conn)
		close(done)
	}()

	rrq := ReadReq{Filename: "test"}

	client, err := net.ListenPacket("udp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	b, err := rrq.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}

	n, err := client.WriteTo(b, conn.LocalAddr())
	if err != nil {
		t.Fatal(err)
	}

	if n != len(b) {
		t.Fatalf("expected %d bytes; wrote %d bytes", len(b), n)
	}

	p2 := new(bytes.Buffer)

	for {
		_ = client.SetReadDeadline(time.Now().Add(time.Second))

		buf := make([]byte, DatagramSize)

		n, addr, err := client.ReadFrom(buf)
		if err != nil {
			t.Fatal(err)
		}

		var data Data

		err = data.UnmarshalBinary(buf[:n])
		if err != nil {
			t.Fatal(err)
		}

		_, err = io.Copy(p2, data.Payload)
		if err != nil {
			t.Fatal(err)
		}

		ack := Ack(data.Block)

		b, err = ack.MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}

		_, err = client.WriteTo(b, addr)
		if err != nil {
			t.Fatal(err)
		}

		if n < DatagramSize {
			break
		}
	}

	_ = client.Close()
	_ = conn.Close()

	<-done

	if !bytes.Equal(p1, p2.Bytes()) {
		t.Fatal("sent payload not equal to received payload")
	}
}
