package encoder

import (
	"io"
)

type TableEncoder interface {
	SetHeader([]string)
	Append([]string)
	Render()
}

type FileType int

const (
	None FileType = iota
	CSV
	TSV
)

func (f FileType) Encoder(w io.Writer) TableEncoder {
	switch f {
	case None:
		return NewPrettyTableEncoder(w)
	case CSV:
		return NewCSVTableEncoder(w)
	case TSV:
		return NewTSVTableEncoder(w)
	}
	return nil
}
