package unicode

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
)

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
