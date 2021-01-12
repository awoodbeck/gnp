package protobuf

import (
	"io"
	"io/ioutil"

	"google.golang.org/protobuf/proto"

	"github.com/awoodbeck/gnp/ch12/housework/v1"
)

func Load(r io.Reader) ([]*housework.Chore, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var chores housework.Chores

	return chores.Chores, proto.Unmarshal(b, &chores)
}

func Flush(w io.Writer, chores []*housework.Chore) error {
	b, err := proto.Marshal(&housework.Chores{Chores: chores})
	if err != nil {
		return err
	}

	_, err = w.Write(b)

	return err
}
