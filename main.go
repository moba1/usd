package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"golang.org/x/text/unicode/runenames"
)

func main() {
	const (
		utf8CmdName  = "utf8"
		utf16CmdName = "utf16"
		utf32CmdName = "utf32"
	)
	utf8Cmd := flag.NewFlagSet(utf8CmdName, flag.ExitOnError)
	utf16Cmd := flag.NewFlagSet(utf16CmdName, flag.ExitOnError)
	utf32Cmd := flag.NewFlagSet(utf32CmdName, flag.ExitOnError)
	var (
		subCmd     string
		subCmdArgs []string
	)
	if len(os.Args) < 2 {
		subCmd = utf8CmdName
		subCmdArgs = []string{}
	} else {
		subCmd = os.Args[1]
		subCmdArgs = os.Args[2:]
	}
	var (
		reader      func(*bufio.Reader) (rune, []byte, error)
		parseEndian = func(endianHolder *Endian, s string) error {
			switch s {
			case "Big":
				*endianHolder = BigEndian
			case "Little":
				*endianHolder = LittleEndian
			default:
				return fmt.Errorf("invalid endian: %s", s)
			}
			return nil
		}
	)
	switch subCmd {
	case utf8CmdName:
		utf8Cmd.Parse(subCmdArgs)
		reader = ReadUtf8Char
	case utf16CmdName:
		var endian Endian = BigEndian
		utf16Cmd.Func("endian", "UTF16 endian", func(s string) error {
			return parseEndian(&endian, s)
		})
		utf16Cmd.Parse(subCmdArgs)
		reader = func(buf *bufio.Reader) (rune, []byte, error) {
			return ReadUtf16Char(endian, buf)
		}
	case utf32CmdName:
		var endian Endian = BigEndian
		utf32Cmd.Func("endian", "UTF32 endian", func(s string) error {
			return parseEndian(&endian, s)
		})
		utf32Cmd.Parse(subCmdArgs)
		reader = func(buf *bufio.Reader) (rune, []byte, error) {
			return ReadUtf32Char(endian, buf)
		}
	default:
		println(fmt.Sprintf("invalid command: %s", subCmd))
		os.Exit(1)
	}

	buf := bufio.NewReader(os.Stdin)
	runeTable := tablewriter.NewWriter(os.Stdout)
	runeTable.SetHeader([]string{"Character", "Code Point", "Name", "Hex"})
	for {
		c, bs, err := reader(buf)
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
			return
		}

		toHexString := func(bs []byte) string {
			s := ""
			for _, b := range bs {
				s = fmt.Sprintf("%s\\x%x", s, b)
			}
			return s
		}
		runeTable.Append([]string{
			strconv.QuoteToGraphic(string(c)),
			fmt.Sprintf("%U", c),
			runenames.Name(c),
			toHexString(bs),
		})
	}
	runeTable.Render()
}
