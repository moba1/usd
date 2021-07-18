package encoder_test

import (
	"reflect"
	"testing"

	"github.com/moba1/usd/encoder"
)

func TestFileType_Encoder(t *testing.T) {
	e := encoder.None.Encoder(nil)
	if _, ok := e.(*encoder.PrettyTableEncoder); !ok {
		t.Errorf("expected Type is *encoder.PrettyTableEncoder, but actual type is %v", reflect.TypeOf(e))
	}
	e = encoder.CSV.Encoder(nil)
	if _, ok := e.(*encoder.CSVTableEncoder); !ok {
		t.Errorf("expected Type is *encoder.CSVTableEncoder, but actual type is %v", reflect.TypeOf(e))
	}
	e = encoder.TSV.Encoder(nil)
	if _, ok := e.(*encoder.TSVTableEncoder); !ok {
		t.Errorf("expected Type is *encoder.TSVTableEncoder, but actual type is %v", reflect.TypeOf(e))
	}
	e = encoder.FileType(-1).Encoder(nil)
	if e != nil {
		t.Errorf("expected nil, but actual type is %v", reflect.TypeOf(e))
	}
}
