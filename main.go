package main

import (
	"bufio"
	// "bytes"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"golang.org/x/text/unicode/runenames"
)

func main() {
	// b := []byte{
	//   0x30, 0x42,
	//   0x30, 0x44,
	//   0x30, 0x46,
	//   0x30, 0x48,
	//   0x30, 0x4a,
	//   0xd8, 0x69, 0xde, 0xca,
	// }
	// b := []byte("あいうえお🐧\n")
	buf := bufio.NewReader(os.Stdin)
	runeTable := tablewriter.NewWriter(os.Stdout)
	runeTable.SetHeader([]string{"Character", "Code Point", "Name", "Hex"})
	for {
		// c, err := ReadUtf16Char(BigEndian, buf)
		c, err := ReadUtf8Char(buf)
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
			toHexString([]byte(string(c))),
		})
	}
	runeTable.Render()
}
