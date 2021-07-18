package encoder

type TableEncoder interface {
	SetHeader([]string)
	Append([]string)
	Render()
}
