package tftp

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"io/ioutil"
	"testing"
)

func TestAck(t *testing.T) {
	t.Parallel()

	a1 := Ack(42)

	b, err := a1.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}

	var a2 Ack

	err = a2.UnmarshalBinary(b)
	if err != nil {
		t.Fatal(err)
	}

	if a1 != a2 {
		t.Fatalf("expected %d; actual %d", a1, a2)
	}
}

func TestDataUnmarshalBinary(t *testing.T) {
	t.Parallel()

	b := make([]byte, BlockSize)

	_, err := rand.Read(b)
	if err != nil {
		t.Fatal(err)
	}

	pbuf := new(bytes.Buffer)

	err = binary.Write(pbuf, binary.BigEndian, OpData)
	if err != nil {
		t.Fatal(err)
	}

	err = binary.Write(pbuf, binary.BigEndian, uint16(1))
	if err != nil {
		t.Fatal(err)
	}

	d := new(Data)

	err = d.UnmarshalBinary(append(pbuf.Bytes(), b...))
	if err != nil {
		t.Fatal(err)
	}

	if d.Block != 1 {
		t.Errorf("expected block number 1; actual block number %d", d.Block)
	}

	p, err := ioutil.ReadAll(d.Payload)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(b, p) {
		t.Fatal("input payload does not equal output payload")
	}
}

func TestDataMarshalBinary(t *testing.T) {
	t.Parallel()

	var (
		previous uint16
		b1       = make([]byte, 1.5*BlockSize)
		bbuf     = new(bytes.Buffer)
	)

	_, err := rand.Read(b1)
	if err != nil {
		t.Fatal(err)
	}

	d := Data{Payload: bytes.NewReader(b1)}

	for { // each iteration simulates a read from the client
		p, err := d.MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}

		actual := len(p)
		if actual < 4 {
			t.Fatalf("expected at least 4 bytes; read %d bytes", actual)
		}

		if actual > DatagramSize {
			t.Fatalf("expected no more than %d bytes; read %d bytes", DatagramSize, actual)
		}

		opcode := OpCode(binary.BigEndian.Uint16(p[:2]))
		if opcode != OpData {
			t.Fatalf("expected operation code %d; actual %d", OpData, opcode)
		}

		blockcode := binary.BigEndian.Uint16(p[2:4])
		if blockcode != previous+1 {
			t.Fatalf("block code %d is not sequential; previous = %d", blockcode, previous)
		}

		previous = blockcode

		bbuf.Write(p[4:])

		if actual < BlockSize {
			// indicates this is the last block of data
			break
		}
	}

	b2 := bbuf.Bytes()

	if !bytes.Equal(b1, b2) {
		t.Fatalf("input payload does not equal output payload: %d vs %d", len(b1), len(b2))
	}
}

func TestError(t *testing.T) {
	t.Parallel()

	e1 := Err{Error: ErrFileExists, Message: "file already exists"}

	b, err := e1.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}

	// operation code + error code + message + 0 byte
	if expected := (2 + 2 + len(e1.Message) + 1); len(b) != expected {
		t.Fatalf("expected %d bytes; read %d bytes", expected, len(b))
	}

	var e2 Err

	err = e2.UnmarshalBinary(b)
	if err != nil {
		t.Fatal(err)
	}

	if e1.Error != e2.Error {
		t.Errorf("expected error code %d; actual error code %d",
			e1.Error, e2.Error)
	}

	if e1.Message != e2.Message {
		t.Errorf("expected message %q; actual message %q",
			e1.Message, e2.Message)
	}
}

func TestReadReq(t *testing.T) {
	t.Parallel()

	r1 := ReadReq{Filename: "/etc/shadow", Mode: "octet"}

	b, err := r1.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}

	// operation code + filename + 0 byte + mode + 0 byte
	expected := (2 + len(r1.Filename) + 1 + len(r1.Mode) + 1)
	if len(b) != expected {
		t.Fatalf("expected %d bytes; read %d bytes", expected, len(b))
	}

	var r2 ReadReq

	err = r2.UnmarshalBinary(b)
	if err != nil {
		t.Fatal(err)
	}

	if r1.Filename != r2.Filename {
		t.Errorf("expected filename %q; actual filename %q",
			r1.Filename, r2.Filename)
	}

	if r1.Mode != r2.Mode {
		t.Errorf("expected mode %q; actual mode %q", r1.Mode, r2.Mode)
	}
}
