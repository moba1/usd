package encoder

import (
	"encoding/csv"
	"fmt"
	"io"
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

func (cte *CSVTableEncoder) Render() error {
	if len(cte.header) > 0 {
		if err := cte.writer.Write(cte.header); err != nil {
			return fmt.Errorf("can't write csv header (reason; %s)", err.Error())
		}
	}
	for _, line := range cte.lines {
		if err := cte.writer.Write(line); err != nil {
			return fmt.Errorf("can't write csv row (reason; %s)", err.Error())
		}
	}
	cte.writer.Flush()
	return nil
}
