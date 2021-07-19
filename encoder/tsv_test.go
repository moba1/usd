package encoder_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/moba1/usd/encoder"
)

func TestTSVTableEncoder_Render(t *testing.T) {
	var bs []byte = nil
	buf := bytes.NewBuffer(bs)
	cte := encoder.NewTSVTableEncoder(buf)
	header := []string{"HeaderA", "HeaderB"}
	rows := [][]string{
		{"row1A", "row1B"},
		{"row2A", "row2B"},
	}
	cte.SetHeader(header)
	for _, r := range rows {
		cte.Append(r)
	}
	if err := cte.Render(); err != nil {
		t.Fatalf("error occured at TSVTableEncoder.Render (%v)", err)
	}

	expectedTSVTable := fmt.Sprintln(strings.Join(header, "\t"))
	for _, r := range rows {
		expectedTSVTable = fmt.Sprintf("%s%s\n", expectedTSVTable, strings.Join(r, "\t"))
	}
	returnTSVTable, err := ioutil.ReadAll(buf)
	if err != nil {
		t.Errorf("error occured at ioutil.ReadAll(buf) (%v)", err)
	}
	if expectedTSVTable != string(returnTSVTable) {
		t.Errorf("expected CSV Table: %q, but CSVTableEncoder.Render return: %q", expectedTSVTable, returnTSVTable)
	}
}
