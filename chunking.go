package hoard

import (
	"bytes"
	"fmt"
	"io"
)

func CopyChunked(dest func(chunk []byte) error, src func() ([]byte, error), chunkSize int64) error {
	// The internal Buffer will ensure we write in chunks
	_, err := io.CopyBuffer(NewPusher(dest), NewPuller(src), make([]byte, chunkSize))
	if err != io.EOF {
		return err
	}
	return nil
}

type PullReader struct {
	pull func() ([]byte, error)
	buf  bytes.Buffer
}

func NewPuller(pull func() ([]byte, error)) *PullReader {
	return &PullReader{
		pull: pull,
	}
}

func (pr *PullReader) Read(p []byte) (n int, err error) {
	var bs []byte

	if pr.pull == nil {
		return 0, io.EOF
	}

	// Attempt to fill read buffer
	for pr.buf.Len() < len(p) {
		bs, err = pr.pull()
		if err != nil {
			if err == io.EOF {
				// Signal end of pull stream to all subsequent calls
				pr.pull = nil
				break
			}
			return n, fmt.Errorf("PullBuffer: could not pull bytes: %w", err)
		}
		_, err = pr.buf.Write(bs)
		if err != nil {
			return n, fmt.Errorf("PullBuffer: could not write bytes into buffer: %w", err)
		}

	}
	n, err = pr.buf.Read(p)
	if err != nil && err != io.EOF {
		return n, fmt.Errorf("PullBuffer: could not read bytes from buffer: %w", err)
	}

	// Buffer filled or EOF
	return
}

type PushWriter struct {
	push func([]byte) error
}

func NewPusher(push func([]byte) error) *PushWriter {
	return &PushWriter{
		push: push,
	}
}

func (pw *PushWriter) Write(p []byte) (_ int, err error) {
	err = pw.push(p)
	if err != nil {
		return 0, fmt.Errorf("PushWriter: could not push bytes: %w", err)
	}
	return len(p), nil
}
