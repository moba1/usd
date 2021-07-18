package unicode

import (
	"bufio"
	"encoding/binary"
	"io"
	"unicode/utf16"
)

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
