package coding

import (
	"testing"
)

var testEncode = []struct {
	input  byte
	answer []byte
}{
	{231, []byte{116, 58}},
	{72, []byte{39, 69}},
}
var testDecode = []struct {
	input  []byte
	answer byte
}{}

func TestCalcSuccess(t *testing.T) {
	for _, e := range testEncode {
		res := encode(e.input)
		if res[0] != e.answer[0] && res[1] != e.answer[1] {
			t.Errorf("Tryed to encode %b, expected :{%b} {%b} got:{%b} {%b}",
				e.input, e.answer[0], e.answer[1], res[0], res[1])
		}
	}
}
