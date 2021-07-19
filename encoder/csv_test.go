package encoder_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/moba1/usd/encoder"
)

func TestCSVTableEncoder_Render(t *testing.T) {
	var bs []byte = nil
	buf := bytes.NewBuffer(bs)
	cte := encoder.NewCSVTableEncoder(buf)
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
		t.Fatalf("error occured at CSVTableEncoder.Render (%v)", err)
	}

	expectedCSVTable := fmt.Sprintln(strings.Join(header, ","))
	for _, r := range rows {
		expectedCSVTable = fmt.Sprintf("%s%s\n", expectedCSVTable, strings.Join(r, ","))
	}
	returnCSVTable, err := ioutil.ReadAll(buf)
	if err != nil {
		t.Errorf("error occured at ioutil.ReadAll(buf) (%v)", err)
	}
	if expectedCSVTable != string(returnCSVTable) {
		t.Errorf("expected CSV Table: %q, but CSVTableEncoder.Render return: %q", expectedCSVTable, returnCSVTable)
	}
}
