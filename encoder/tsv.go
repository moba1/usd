package encoder

import (
	"fmt"
	"io"
	"strings"
)

type TSVTableEncoder struct {
	writer io.Writer
	lines  [][]string
	header []string
}

func NewTSVTableEncoder(w io.Writer) *TSVTableEncoder {
	return &TSVTableEncoder{
		writer: w,
	}
}

func (tte *TSVTableEncoder) SetHeader(h []string) {
	tte.header = h
}

func (tte *TSVTableEncoder) Append(row []string) {
	tte.lines = append(tte.lines, row)
}

func (tte *TSVTableEncoder) Render() error {
	if len(tte.header) > 0 {
		if _, err := fmt.Fprintln(tte.writer, strings.Join(tte.header, "\t")); err != nil {
			return fmt.Errorf("cannot write tsv header (reason; %s)", err.Error())
		}
	}
	for _, line := range tte.lines {
		if _, err := fmt.Fprintln(tte.writer, strings.Join(line, "\t")); err != nil {
			return fmt.Errorf("cannot write tsv row (reason; %s)", err.Error())
		}
	}
	return nil
}
