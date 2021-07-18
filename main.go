package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"golang.org/x/text/unicode/runenames"
)

var reader func(*bufio.Reader) (rune, []byte, error)

func parseCmd() {
	const (
		utf8CmdName  = "utf8"
		utf16CmdName = "utf16"
		utf32CmdName = "utf32"
	)

	flag.Usage = func() {
		prog := path.Base(os.Args[0])
		stmts := []string{
			fmt.Sprintf("Usage of %s:", prog),
			fmt.Sprintf("  %s <subcommand>", prog),
			fmt.Sprintf("  %s -help", prog),
			"Sub commands:",
			fmt.Sprintf("  %s", utf8CmdName),
			"        dump UTF-8",
			fmt.Sprintf("  %s", utf16CmdName),
			"        dump UTF-16",
			fmt.Sprintf("  %s", utf32CmdName),
			"        dump UTF-32",
			"Options:",
			"  -help",
			"       show help",
		}
		for _, stmt := range stmts {
			fmt.Fprintln(flag.CommandLine.Output(), stmt)
		}
		flag.PrintDefaults()
	}
	flag.Parse()
	args := flag.Args()

	utf8Cmd := flag.NewFlagSet(utf8CmdName, flag.ExitOnError)
	utf16Cmd := flag.NewFlagSet(utf16CmdName, flag.ExitOnError)
	utf32Cmd := flag.NewFlagSet(utf32CmdName, flag.ExitOnError)
	var (
		subCmd     string
		subCmdArgs []string
	)
	if len(args) < 2 {
		subCmd = utf8CmdName
		subCmdArgs = []string{}
	} else {
		subCmd = args[0]
		subCmdArgs = args[1:]
	}
	var parseEndian = func(endianHolder *Endian, s string) error {
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
	switch subCmd {
	case utf8CmdName:
		utf8Cmd.Usage = func() {
			stmts := []string{
				fmt.Sprintf("Usage of %s:", utf8CmdName),
				fmt.Sprintf("  %s [option]", utf8CmdName),
				"Options:",
				"  -help",
				"        show help",
			}
			for _, stmt := range stmts {
				fmt.Fprintln(utf8Cmd.Output(), stmt)
			}
			utf8Cmd.PrintDefaults()
		}
		utf8Cmd.Parse(subCmdArgs)
		reader = ReadUtf8Char
	case utf16CmdName:
		var endian Endian = BigEndian
		utf16Cmd.Func("endian", "UTF16 endian. default is `Big` (value: Big|Little)", func(s string) error {
			return parseEndian(&endian, s)
		})
		utf16Cmd.Usage = func() {
			stmts := []string{
				fmt.Sprintf("Usage of %s:", utf16CmdName),
				fmt.Sprintf("  %s [option]", utf16CmdName),
				"Options:",
				"  -help",
				"        show help",
			}
			for _, stmt := range stmts {
				fmt.Fprintln(utf16Cmd.Output(), stmt)
			}
			utf16Cmd.PrintDefaults()
		}
		utf16Cmd.Parse(subCmdArgs)
		reader = func(buf *bufio.Reader) (rune, []byte, error) {
			return ReadUtf16Char(endian, buf)
		}
	case utf32CmdName:
		var endian Endian = BigEndian
		utf32Cmd.Func("endian", "UTF32 endian. default is `Big` (value: Big|Little)", func(s string) error {
			return parseEndian(&endian, s)
		})
		utf32Cmd.Usage = func() {
			stmts := []string{
				fmt.Sprintf("Usage of %s:", utf32CmdName),
				fmt.Sprintf("  %s [option]", utf32CmdName),
				"Options:",
				"  -help",
				"        show help",
			}
			for _, stmt := range stmts {
				fmt.Fprintln(utf32Cmd.Output(), stmt)
			}
			utf32Cmd.PrintDefaults()
		}
		utf32Cmd.Parse(subCmdArgs)
		reader = func(buf *bufio.Reader) (rune, []byte, error) {
			return ReadUtf32Char(endian, buf)
		}
	default:
		println(fmt.Sprintf("invalid command: %s", subCmd))
		os.Exit(1)
	}
}

func main() {
	parseCmd()

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
