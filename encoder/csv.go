package encoder

import (
	"encoding/csv"
	"io"
	"log"
)

type CSVTableEncoder struct {
	writer *csv.Writer
	lines  [][]string
	header []string
}

func NewCSVTableEncoder(w io.Writer) *CSVTableEncoder {
	return &CSVTableEncoder{
		writer: csv.NewWriter(w),
		lines:  [][]string{},
		header: []string{},
	}
}

func (cte *CSVTableEncoder) SetHeader(h []string) {
	cte.header = h
}

func (cte *CSVTableEncoder) Append(row []string) {
	cte.lines = append(cte.lines, row)
}

func (cte *CSVTableEncoder) Render() {
	if len(cte.header) > 0 {
		if err := cte.writer.Write(cte.header); err != nil {
			log.Fatalf("can't write csv header: %s\n", err.Error())
		}
	}
	for _, line := range cte.lines {
		if err := cte.writer.Write(line); err != nil {
			log.Fatalf("can't write csv: %s\n", err.Error())
		}
	}
	cte.writer.Flush()
}
