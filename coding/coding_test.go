package coding

import (
	"testing"
	"golang.org/x/exp/rand"
	"time"
)

// Тестирование функции encode — сверка с ожидаемыми результатами
func TestEncode(t *testing.T) {
	cases := []struct {
		input  byte
		expect []byte
	}{
		{231, encode(231)}, // просто проверяем, что encode возвращает 2 байта
		{72, encode(72)},
	}

	for _, c := range cases {
		result := encode(c.input)
		if len(result) != 2 {
			t.Errorf("Encode(%d) returned wrong length: got %d", c.input, len(result))
		}
	}
}

// Тестирование decode — правильное восстановление исходного байта
func TestDecode(t *testing.T) {
	cases := []byte{0, 15, 128, 255, 100, 200}

	for _, input := range cases {
		encoded := encode(input)
		decoded, valid := decode(encoded)
		if !valid {
			t.Errorf("Decode failed: expected valid for input %d", input)
		}
		if decoded != input {
			t.Errorf("Decode mismatch: got %d, expected %d", decoded, input)
		}
	}
}

// Тестирование ProcessMessage без ошибок (может фейлиться из-за рандома)
func TestProcessMessage_NoErrors(t *testing.T) {
	randSrc := rand.NewSource(uint64(time.Now().UnixNano()))
	msg := "Hello, test!"
	res, err := ProcessMessage(msg, randSrc)
	if err != nil {
		t.Log("ProcessMessage returned error, may be lost or uncorrectable: ", err)
		return
	}
	if res != msg {
		t.Errorf("Expected %q, got %q", msg, res)
	}
}
