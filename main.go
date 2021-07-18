package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"golang.org/x/text/unicode/runenames"
)

func main() {
	buf := bufio.NewReader(os.Stdin)
	runeTable := tablewriter.NewWriter(os.Stdout)
	runeTable.SetHeader([]string{"Character", "Code Point", "Name", "Hex"})
	for {
		c, bs, err := ReadUtf8Char(buf)
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
