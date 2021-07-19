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

func testReadUtf32Char(endian unicode.Endian, cases TestCases, t *testing.T) {
	reader := func(buf *bufio.Reader) (rune, []byte, error) {
		return unicode.ReadUtf32Char(endian, buf)
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
			t.Errorf("ReadUtf32Char returns error: %v", err)
		}
		if r != c.char {
			t.Errorf("ReadUtf32Char returns character %s, but expected character is %s", strconv.QuoteRuneToGraphic(r), strconv.QuoteRuneToGraphic(c.char))
		}
		if !reflect.DeepEqual(bs, c.byteStream) {
			t.Errorf("ReadUtf32Char returns %v, but expected value is %v", bs, c.byteStream)
		}
	}

	r, bs, err := reader(buf)
	if err != io.EOF {
		t.Errorf("expected EOF, but read character: %c (byte stream: %v)", r, bs)
	}

	// when ReadUtf8Char read empty stream, return io.EOF
	r, bs, err = reader(bufio.NewReader(bytes.NewBuffer([]byte{})))
	if err != io.EOF {
		t.Errorf("ReadUtf32Char read empty byte stream, but return non EOF error")
	}

	/*********************
	 * fail
	 *********************/
	// including invalid byte
	for _, c := range cases.fail.invalidSequeces {
		buf = bufio.NewReader(bytes.NewBuffer(c))
		r, bs, err = reader(buf)
		var invalidSequenceErr *unicode.InvalidSequenceErr
		if !errors.As(err, &invalidSequenceErr) {
			t.Errorf("ReadUtf32Char read invalid sequences: %v", cases.fail.invalidSequeces)
		}
	}

	// lack needed byt_Litte
	for _, c := range cases.fail.lackedSeqences {
		buf = bufio.NewReader(bytes.NewBuffer(c))
		r, bs, err = reader(buf)
		var unexpectedEofErr *unicode.UnexpectedEofErr
		if !errors.As(err, &unexpectedEofErr) {
			t.Errorf("ReadUtf32Char read short sequences: %v", c)
		}
	}

	// I/O error occured
	buf = bufio.NewReader(iotest.ErrReader(errors.New("general I/O error")))
	_, _, err = reader(buf)
	if err == nil {
		t.Errorf("ReadUtf32Char ignore I/O error")
	}
}

func TestReadUtf32Char_Error(t *testing.T) {
	buf := bufio.NewReader(bytes.NewBuffer([]byte{0x00, 0x00, 0x00, 0x00}))
	_, _, err := unicode.ReadUtf32Char(-1, buf)
	var unknownEndianErr *unicode.UnknownEndianErr
	if !errors.As(err, &unknownEndianErr) {
		t.Errorf("ReadUtf32Char return non-UnknownEndianErr: %v", err)
	}
}

func TestReadUtf32Char_BigEndian(t *testing.T) {
	testCases := TestCases{
		success: []Char{
			{
				char:       utf16.Decode([]uint16{0xFFFE})[0],
				byteStream: []byte{0x00, 0x00, 0xFF, 0xFE},
			},
			{
				char:       '1',
				byteStream: []byte{0x00, 0x00, 0x00, 0x31},
			},
			{
				char:       '\x00',
				byteStream: []byte{0x00, 0x00, 0x00, 0x00},
			},
			{
				char:       '△',
				byteStream: []byte{0x00, 0x00, 0x25, 0xB3},
			},
		},
		fail: struct {
			invalidSequeces [][]byte
			lackedSeqences  [][]byte
		}{
			lackedSeqences: [][]byte{
				{
					0x00, // lack 3 byte
				},
				{
					0x00, 0x00, 0x25, // lack 1 byte
				},
			},
		},
	}
	testReadUtf32Char(unicode.BigEndian, testCases, t)
}

func TestReadUtf32Char_LittleEndian(t *testing.T) {
	testCases := TestCases{
		success: []Char{
			{
				char:       utf16.Decode([]uint16{0xFFFE})[0],
				byteStream: []byte{0xFE, 0xFF, 0x00, 0x00},
			},
			{
				char:       'あ',
				byteStream: []byte{0x42, 0x30, 0x00, 0x00},
			},
			{
				char:       '🐧',
				byteStream: []byte{0x27, 0xF4, 0x01, 0x00},
			},
		},
		fail: struct {
			invalidSequeces [][]byte
			lackedSeqences  [][]byte
		}{
			lackedSeqences: [][]byte{
				{
					0x00, // lack 3 byte
				},
				{
					0xB3, 0x25, 0x00, // lack 1 byte
				},
			},
		},
	}
	testReadUtf32Char(unicode.LittleEndian, testCases, t)
}
