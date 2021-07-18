package unicode

import (
	"bufio"
	"fmt"
	"io"
)

type Endian int

const (
	BigEndian Endian = iota
	LittleEndian
)

type InvalidSequenceErr struct {
	sequences []byte
}

func (e *InvalidSequenceErr) Error() string {
	return fmt.Sprintf("invalid sequences: %#v", e.sequences)
}

func (e *InvalidSequenceErr) Sequences() []byte {
	return e.sequences
}

type UnknownEndianErr struct{}

func (*UnknownEndianErr) Error() string {
	return "unknown endian"
}

type UnexpectedEofErr struct{}

func (*UnexpectedEofErr) Error() string {
	return "unexpected eof"
}

func readMultiByte(buf *bufio.Reader, bs []byte) error {
	n, err := buf.Read(bs)
	if 0 < n && n < len(bs) && err == io.EOF {
		return &InvalidSequenceErr{
			sequences: bs,
		}
	}
	return err
}
