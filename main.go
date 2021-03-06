package main

import (
	"bufio"
	_ "embed"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/moba1/usd/encoder"
	"github.com/moba1/usd/unicode"
	"golang.org/x/text/unicode/runenames"
)

//go:embed version
var version string

var (
	reader      func(*bufio.Reader) (rune, []byte, error)
	fileType    encoder.FileType
	noHeader    bool
	showVersion bool
)

func init() {
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
	flag.BoolVar(&showVersion, "version", false, "show version")
	flag.Func("fileType", "output file type. default is None (value: CSV|TSV|None)", func(s string) error {
		switch s {
		case "CSV":
			fileType = encoder.CSV
		case "TSV":
			fileType = encoder.TSV
		case "None":
			fileType = encoder.None
		default:
			return fmt.Errorf("invalid file type: %s", s)
		}
		return nil
	})
	flag.BoolVar(&noHeader, "noHeader", false, "no header")
	flag.Parse()
	if showVersion {
		fmt.Println(version)
		os.Exit(0)
	}
	args := flag.Args()

	utf8Cmd := flag.NewFlagSet(utf8CmdName, flag.ExitOnError)
	utf16Cmd := flag.NewFlagSet(utf16CmdName, flag.ExitOnError)
	utf32Cmd := flag.NewFlagSet(utf32CmdName, flag.ExitOnError)
	var (
		subCmd     string
		subCmdArgs []string
	)
	if len(args) < 1 {
		subCmd = utf8CmdName
		subCmdArgs = []string{}
	} else {
		subCmd = args[0]
		subCmdArgs = args[1:]
	}
	var parseEndian = func(endianHolder *unicode.Endian, s string) error {
		switch s {
		case "Big":
			*endianHolder = unicode.BigEndian
		case "Little":
			*endianHolder = unicode.LittleEndian
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
		if err := utf8Cmd.Parse(subCmdArgs); err != nil {
			log.Fatalln(err)
		}
		reader = unicode.ReadUtf8Char
	case utf16CmdName:
		var endian unicode.Endian = unicode.BigEndian
		utf16Cmd.Func("endian", "UTF16 `endian`. default is 'Big' (value: Big|Little)", func(s string) error {
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
		if err := utf16Cmd.Parse(subCmdArgs); err != nil {
			log.Fatalln(err)
		}
		reader = func(buf *bufio.Reader) (rune, []byte, error) {
			return unicode.ReadUtf16Char(endian, buf)
		}
	case utf32CmdName:
		var endian unicode.Endian = unicode.BigEndian
		utf32Cmd.Func("endian", "UTF32 `endian`. default is 'Big' (value: Big|Little)", func(s string) error {
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
		if err := utf32Cmd.Parse(subCmdArgs); err != nil {
			log.Fatalln(err)
		}
		reader = func(buf *bufio.Reader) (rune, []byte, error) {
			return unicode.ReadUtf32Char(endian, buf)
		}
	default:
		println(fmt.Sprintf("invalid command: %s", subCmd))
		os.Exit(1)
	}
}

func main() {
	runeTable := fileType.Encoder(os.Stdout)
	if !noHeader {
		runeTable.SetHeader([]string{"Character", "Code Point", "Name", "Hex"})
	}

	buf := bufio.NewReader(os.Stdin)
	for {
		c, bs, err := reader(buf)
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
			return
		}

		toHexString := func(bs []byte) string {
			hexes := []string{}
			for _, b := range bs {
				hexes = append(hexes, fmt.Sprintf("0x%02X", b))
			}
			return strings.Join(hexes, " ")
		}
		graphic := strings.Trim(strconv.QuoteRuneToGraphic(c), "'")
		if c == '\'' {
			graphic = "'"
		}
		runeTable.Append([]string{
			graphic,
			fmt.Sprintf("%U", c),
			runenames.Name(c),
			toHexString(bs),
		})
	}
	if err := runeTable.Render(); err != nil {
		log.Fatalln(err)
	}
}
