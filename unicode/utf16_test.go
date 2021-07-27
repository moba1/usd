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
	"unicode/utf16"

	"github.com/moba1/usd/unicode"
)

func testReadUtf16Char(endian unicode.Endian, cases TestCases, t *testing.T) {
	reader := func(buf *bufio.Reader) (rune, []byte, error) {
		return unicode.ReadUtf16Char(endian, buf)
	}

	/*********************
	 * success
	 *********************/
	var bs []byte
	for _, c := range cases.success {
		bs = append(bs, c.byteStream...)
	}
	buf := bufio.NewReader(bytes.NewBuffer(bs))
	for _, c := range cases.success {
		r, bs, err := reader(buf)
		if err != nil {
			t.Errorf("ReadUtf16Char returns error: %v", err)
		}
		if r != c.char {
			t.Errorf("ReadUtf16Char returns character %s, but expected character is %s", strconv.QuoteRuneToGraphic(r), strconv.QuoteRuneToGraphic(c.char))
		}
		if !reflect.DeepEqual(bs, c.byteStream) {
			t.Errorf("ReadUtf16Char returns %v, but expected value is %v", bs, c.byteStream)
		}
	}

	r, bs, err := reader(buf)
	if err != io.EOF {
		t.Errorf("expected EOF, but read character: %c (byte stream: %v)", r, bs)
	}

	// when ReadUtf8Char read empty stream, return io.EOF
	_, _, err = reader(bufio.NewReader(bytes.NewBuffer([]byte{})))
	if err != io.EOF {
		t.Errorf("ReadUtf16Char read empty byte stream, but return non EOF error")
	}

	/*********************
	 * fail
	 *********************/
	// including invalid byte
	for _, c := range cases.fail.invalidSequeces {
		buf = bufio.NewReader(bytes.NewBuffer(c))
		_, _, err = reader(buf)
		var invalidSequenceErr *unicode.InvalidSequenceErr
		if !errors.As(err, &invalidSequenceErr) {
			t.Errorf("ReadUtf16Char read invalid sequences: %v", cases.fail.invalidSequeces)
		}
	}

	// lack needed byt_Litte
	for _, c := range cases.fail.lackedSeqences {
		buf = bufio.NewReader(bytes.NewBuffer(c))
		_, _, err = reader(buf)
		var unexpectedEofErr *unicode.UnexpectedEofErr
		if !errors.As(err, &unexpectedEofErr) {
			t.Errorf("ReadUtf16Char read short sequences: %v", c)
		}
	}

	// I/O error occured
	buf = bufio.NewReader(iotest.ErrReader(errors.New("general I/O error")))
	_, _, err = reader(buf)
	if err == nil {
		t.Errorf("ReadUtf16Char ignore I/O error")
	}
}

func TestReadUtf16Char_BigEndian(t *testing.T) {
	testCases := TestCases{
		success: []Char{
			{
				char:       utf16.Decode([]uint16{0xFEFF})[0],
				byteStream: []byte{0xFE, 0xFF},
			},
			{
				char:       'a',
				byteStream: []byte{0x00, 0x61},
			},
			{
				char:       '🐧',
				byteStream: []byte{0xD8, 0x3D, 0xDC, 0x27},
			},
		},
		fail: struct {
			invalidSequeces [][]byte
			lackedSeqences  [][]byte
		}{
			invalidSequeces: [][]byte{
				{
					0xD8, 0x00, 0x00, 0x61, // invalid surrogate pair
				},
			},
			lackedSeqences: [][]byte{
				{
					0x00, // lack one byte
				},
				{
					0xD8, 0x00, // lack surrogate pair byte
				},
			},
		},
	}
	testReadUtf16Char(unicode.BigEndian, testCases, t)
}

func TestReadUtf16Char_LittleEndian(t *testing.T) {
	testCases := TestCases{
		success: []Char{
			{
				char:       utf16.Decode([]uint16{0xFEFF})[0],
				byteStream: []byte{0xFF, 0xFE},
			},
			{
				char:       '\n',
				byteStream: []byte{0x0a, 0x00},
			},
			{
				char:       '🛀',
				byteStream: []byte{0x3D, 0xD8, 0xC0, 0xDE},
			},
			{
				char:       'は',
				byteStream: []byte{0x6F, 0x30},
			},
		},
		fail: struct {
			invalidSequeces [][]byte
			lackedSeqences  [][]byte
		}{
			invalidSequeces: [][]byte{
				{
					0x00, 0xD8, 0x0a, 0x00, //invalid surrogate pair
				},
			},
			lackedSeqences: [][]byte{
				{
					0x61, // lack one byte
				},
				{
					0x00, 0xD8, // lack surrogate pair byte
				},
			},
		},
	}
	testReadUtf16Char(unicode.LittleEndian, testCases, t)
}
