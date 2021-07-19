package unicode

import (
	"bufio"
	"io"
	"unicode/utf8"
)

func ReadUtf8Char(buf *bufio.Reader) (rune, []byte, error) {
	b1, err := buf.ReadByte()
	if err != nil {
		return 0, nil, err
	}

	if b1 <= 0b0111_1111 {
		seqs := []byte{b1}
		r, _ := utf8.DecodeRune(seqs)
		return r, seqs, nil
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
	if err != nil {
		return 0, nil, err
	}
	r, err := toRune(b1, remainBytes)
	if err != nil {
		return 0, nil, err
	}
	return r, append([]byte{b1}, remainBytes...), nil
}
