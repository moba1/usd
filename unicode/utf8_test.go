package unicode_test

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"reflect"
	"strconv"
	"testing"
	"testing/iotest"

	"github.com/moba1/usd/unicode"
)

func TestReadUtf8Char(t *testing.T) {
	/*********************
	 * success
	 *********************/
	chars := []Char{
		{
			char:       '\n',
			byteStream: []byte{0x0A},
		},
		{
			char:       '🛀',
			byteStream: []byte{0xF0, 0x9F, 0x9B, 0x80},
		},
	}
	var bs []byte
	for _, c := range chars {
		bs = append(bs, c.byteStream...)
	}
	buf := bufio.NewReader(bytes.NewBuffer(bs))
	for _, c := range chars {
		r, bs, err := unicode.ReadUtf8Char(buf)
		if err != nil {
			t.Errorf("unicode.ReadUtf8Char returns error: %v", err)
		}
		if r != c.char {
			t.Errorf("ReadUtf16Char returns character %s, but expected character is %s", strconv.QuoteRuneToGraphic(r), strconv.QuoteRuneToGraphic(c.char))
		}
		if !reflect.DeepEqual(bs, c.byteStream) {
			t.Errorf("unicode.ReadUtf8Char returns %v, but expected value is %v", bs, c.byteStream)
		}
	}
	r, bs, err := unicode.ReadUtf8Char(buf)
	if err != io.EOF {
		t.Errorf("expected EOF, but read character: %c (byte stream: %v)", r, bs)
	}

	// when unicode.ReadUtf8Char read empty stream, return io.EOF
	_, _, err = unicode.ReadUtf8Char(bufio.NewReader(bytes.NewBuffer([]byte{})))
	if err != io.EOF {
		t.Errorf("unicode.ReadUtf8Char read empty byte stream, but return non EOF error")
	}

	/*********************
	 * fail
	 *********************/
	// including invalid byte
	origSeqs := []byte{
		0xFE, // invalid sequence
	}
	buf = bufio.NewReader(bytes.NewBuffer(origSeqs))
	_, _, err = unicode.ReadUtf8Char(buf)
	var invalidSequenceErr *unicode.InvalidSequenceErr
	if !errors.As(err, &invalidSequenceErr) {
		t.Errorf("unicode.ReadUtf8Char read invalid sequences: %v", origSeqs)
	}

	// lack needed byte
	origSeqs = []byte{
		0xF0, 0x9F, 0x9B,
	}
	buf = bufio.NewReader(bytes.NewBuffer(origSeqs))
	_, _, err = unicode.ReadUtf8Char(buf)
	var unexpectedEofErr *unicode.UnexpectedEofErr
	if !errors.As(err, &unexpectedEofErr) {
		t.Errorf("unicode.ReadUtf8Char read invalid sequences: %v", origSeqs)
	}

	// I/O error occured
	buf = bufio.NewReader(iotest.ErrReader(errors.New("general I/O error")))
	_, _, err = unicode.ReadUtf8Char(buf)
	if err == nil {
		t.Error("unicode.ReadUtf8Char ignore I/O error")
	}
}
