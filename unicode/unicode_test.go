package unicode_test

type Char struct {
	char       rune
	byteStream []byte
}

type TestCases struct {
	success []Char
	fail    struct {
		invalidSequeces [][]byte
		lackedSeqences  [][]byte
	}
}
