package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"unicode/utf16"
	"unicode/utf8"
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

func ReadUtf8Char(buf *bufio.Reader) (rune, []byte, error) {
	b1, err := buf.ReadByte()
	if err != nil {
		return 0, nil, err
	}

	if b1 <= 0b0111_1111 {
		r, _ := utf8.DecodeRune([]byte{b1})
		return r, nil, nil
	}

	readRemainBytes := func(n int) ([]byte, error) {
		bs := make([]byte, n)
		err := readMultiByte(buf, bs)
		if err == io.EOF {
			err = &UnexpectedEofErr{}
		}
		if err != nil {
			return nil, err
		}
		return bs, nil
	}
	toRune := func(head byte, tail []byte) (rune, error) {
		seqs := append([]byte{head}, tail...)
		for _, b := range tail {
			if b&0b1100_0000 != 0b1000_0000 {
				return 0, &InvalidSequenceErr{
					sequences: seqs,
				}
			}
		}
		r, _ := utf8.DecodeRune(seqs)
		return r, nil
	}
	var readByte int
	if b1&0b1110_0000 == 0b1100_0000 {
		readByte = 1
	} else if b1&0b1111_0000 == 0b1110_0000 {
		readByte = 2
	} else if b1&0b1111_1000 == 0b1111_0000 {
		readByte = 3
	} else {
		return 0, nil, &InvalidSequenceErr{
			sequences: []byte{b1},
		}
	}
	remainBytes, err := readRemainBytes(readByte)
	if err == io.EOF {
		err = &UnexpectedEofErr{}
	}
	if err != nil {
		return 0, nil, err
	}
	r, err := toRune(b1, remainBytes)
	if err != nil {
		return 0, nil, err
	}
	return r, append([]byte{b1}, remainBytes...), nil
}

func ReadUtf16Char(endian Endian, buf *bufio.Reader) (rune, []byte, error) {
	readChar := func() ([]byte, error) {
		bs := make([]byte, 2)
		if err := readMultiByte(buf, bs); err != nil {
			return nil, err
		}
		return bs, nil
	}
	toUint16 := func(bs []byte) (uint16, error) {
		var rawChar uint16
		switch endian {
		case BigEndian:
			rawChar = binary.BigEndian.Uint16(bs)
		case LittleEndian:
			rawChar = binary.LittleEndian.Uint16(bs)
		default:
			return 0, &UnknownEndianErr{}
		}
		return rawChar, nil
	}

	r1Bytes, err := readChar()
	if err != nil {
		return 0, nil, err
	}
	r1, err := toUint16(r1Bytes)
	if err != nil {
		return 0, nil, err
	}
	if 0xD800 <= r1 && r1 <= 0xDBFF {
		r2Bytes, err := readChar()
		if err == io.EOF {
			err = &UnexpectedEofErr{}
		}
		if err != nil {
			return 0, nil, err
		}
		r2, err := toUint16(r2Bytes)
		if err != nil {
			return 0, nil, err
		}
		if !(0xDC00 <= r2 && r2 <= 0xDFFF) {
			return 0, nil, &InvalidSequenceErr{
				sequences: append(r1Bytes, r2Bytes...),
			}
		}
		return utf16.Decode([]uint16{r1, r2})[0], append(r1Bytes, r2Bytes...), nil
	}

	return utf16.Decode([]uint16{r1})[0], r1Bytes, nil
}

func ReadUtf32Char(endian Endian, buf *bufio.Reader) (rune, []byte, error) {
	bs := make([]byte, 4)
	if err := readMultiByte(buf, bs); err != nil {
		return 0, nil, err
	}

	readInt32 := func(data []byte, endian binary.ByteOrder) (int32, error) {
		var r int32
		buf := bytes.NewBuffer(data)
		err := binary.Read(buf, endian, &r)
		if err == io.EOF {
			err = &UnexpectedEofErr{}
		}
		if err != nil {
			return 0, err
		}
		return r, nil
	}
	var e binary.ByteOrder
	switch endian {
	case LittleEndian:
		e = binary.LittleEndian
	case BigEndian:
		e = binary.BigEndian
	default:
		return 0, nil, &UnknownEndianErr{}
	}
	r, readInt32Err := readInt32(bs, e)
	if readInt32Err != nil {
		return 0, nil, readInt32Err
	}
	return r, bs, nil
}
