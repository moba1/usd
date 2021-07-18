package encoder

import (
	"io"

	"github.com/olekukonko/tablewriter"
)

type PrettyTableEncoder struct {
	tableWriter *tablewriter.Table
}

func NewPrettyTableEncoder(w io.Writer) *PrettyTableEncoder {
	return &PrettyTableEncoder{
		tableWriter: tablewriter.NewWriter(w),
	}

}

func (pte *PrettyTableEncoder) SetHeader(h []string) {
	pte.tableWriter.SetHeader(h)
}

func (pte *PrettyTableEncoder) Append(row []string) {
	pte.tableWriter.Append(row)
}

func (pte *PrettyTableEncoder) Render() {
	pte.tableWriter.Render()
}
