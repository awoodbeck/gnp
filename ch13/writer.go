package ch13

import (
	"io"

	"go.uber.org/multierr"
)

type sustainedMultiWriter struct {
	writers []io.Writer
}

func (s *sustainedMultiWriter) Write(p []byte) (n int, err error) {
	for _, w := range s.writers {
		i, wErr := w.Write(p)
		n += i
		err = multierr.Append(err, wErr)
	}

	return n, err
}

func SustainedMultiWriter(writers ...io.Writer) io.Writer {
	mw := &sustainedMultiWriter{writers: make([]io.Writer, 0, len(writers))}

	for _, w := range writers {
		if m, ok := w.(*sustainedMultiWriter); ok {
			mw.writers = append(mw.writers, m.writers...)
			continue
		}

		mw.writers = append(mw.writers, w)
	}

	return mw
}
