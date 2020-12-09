package streamer

import (
	"context"
	"fmt"
	"io"
)

const KiB = 1 << 10

const DefaultChunkSize = 64 * KiB

// Streamer is an abstraction for GRPC-ish flows
type Streamer struct {
	chunkSize    int64
	reader       io.Reader
	bytesRead    int64
	send         func(chunk []byte) error
	closeSend    func() error
	recv         func() ([]byte, error)
	writer       io.Writer
	bytesWritten int64
}

// Streamer provides provides non-blocking or bi/uni-directional stream from a io.Reader into a send() function
// and from a recv() into a io.Writer. Each of reader, send(), recv(), writer are optional, so it is possible to
// Implement only a uni-directional scream, or one that does not read from or write to io but instead just uses
// the functions to accumulate state as a side-effect
func New() *Streamer {
	return &Streamer{
		chunkSize: DefaultChunkSize,
		send: func(chunk []byte) error {
			return io.EOF
		},
		closeSend: func() error {
			return nil
		},
		recv: func() ([]byte, error) {
			return nil, io.EOF
		},
	}
}

func (s *Streamer) WithChunkSize(chunkSize int64) *Streamer {
	s.chunkSize = chunkSize
	return s
}

func (s *Streamer) WithInput(reader io.Reader) *Streamer {
	s.reader = reader
	return s
}

func (s *Streamer) WithSend(send func(chunk []byte) error) *Streamer {
	s.send = send
	return s
}

func (s *Streamer) WithCloseSend(closeSend func() error) *Streamer {
	s.closeSend = closeSend
	return s
}

func (s *Streamer) WithRecv(recv func() ([]byte, error)) *Streamer {
	s.recv = recv
	return s
}

func (s *Streamer) WithOutput(writer io.Writer) *Streamer {
	s.writer = writer
	return s
}

// Perform a blocking read from reader to send()
func (s *Streamer) FromReader(ctx context.Context) (err error) {
	out := make([]byte, s.chunkSize)
	n := 0
	for ctx.Err() == nil {
		if s.reader != nil {
			n, err = s.reader.Read(out)
			s.bytesRead += int64(n)
			if err != nil {
				if err == io.EOF {
					return s.closeSend()
				}
				return err
			}
		}
		err = s.send(out[:n])
		if err != nil {
			if err == io.EOF {
				return s.closeSend()
			}
			return err
		}
	}
	return ctx.Err()
}

// A blocking read from recv() to writer
func (s *Streamer) ToWriter(ctx context.Context) error {
	for ctx.Err() == nil {
		chunk, err := s.recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		if s.writer != nil {
			n, err := s.writer.Write(chunk)
			s.bytesWritten += int64(n)
			if err != nil {
				return err
			} else if n != len(chunk) {
				return fmt.Errorf("failed to write data")
			}
		}
	}
	return ctx.Err()
}

func (s *Streamer) Stream(ctx context.Context) error {
	_, _, err := s.StreamCount(ctx)
	return err
}

// Performs a non-blocking stream across the pipeline: reader -> send() -> recv() -> writer
// Though note the connection between send() and recv() is implicit to allow for object streams
// and uni-directional streams. Returns the bytes read and the bytes written (if any) or an error
func (s *Streamer) StreamCount(ctx context.Context) (read int64, written int64, err error) {
	readCh := make(chan error)
	writeCh := make(chan error)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		readCh <- s.FromReader(ctx)
	}()
	go func() {
		writeCh <- s.ToWriter(ctx)
	}()

	wait := 2
	for wait > 0 {
		select {
		case err := <-readCh:
			if err != nil {
				cancel()
				return s.bytesRead, s.bytesWritten, err
			}
			wait--
		case err := <-writeCh:
			if err != nil {
				cancel()
				return s.bytesRead, s.bytesWritten, err
			}
			wait--
		default:
		}
	}

	return s.bytesRead, s.bytesWritten, nil
}
