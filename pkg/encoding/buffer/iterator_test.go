package buffer_test

import (
	"reflect"
	"testing"

	"github.com/Shresth72/ascii/pkg/encoding/buffer"
	"github.com/Shresth72/ascii/pkg/encoding/utils"
	"github.com/stretchr/testify/assert"
)

func ProcessString(str string) ([][]interface{}, map[string]int) {
	byteArray := []byte(str)
	pout := [][]interface{}{}

	iterator := utils.New16BitIterator(byteArray)

	for {
		result := iterator.Next()
		output := make([]interface{}, 3)

		buf16 := make([]byte, 2)
		utils.Write16(buf16, 0, result.Value)
		output[0] = result.Value
		output[1] = string(buf16)

		buf8 := make([]byte, 1)
		writer8 := utils.U8Writer{}
		writer8.Set(buf8)
		writer8.Write(result.Value)
		output[2] = string(buf8)

		pout = append(pout, output)
		if result.Done {
			break
		}
	}

	reader := utils.New8BitIterator(byteArray)
	freq := buffer.NewFrequency()
	freq.Freq(reader)

	frequencies := make(map[string]int)
	for _, v := range freq.PointMap {
		frequencies[string(rune(v.Val))] = v.Count
	}

	return pout, frequencies
}

func TestProcessingString(t *testing.T) {
	str := "aaabbcde"
	expected := [][]interface{}{
		{24929, "aa", "a"},
		{24930, "ab", "b"},
		{25187, "bc", "c"},
		{25701, "de", "e"},
	}

	expectedFreq := map[string]int{
		"a": 3,
		"b": 2,
		"c": 1,
		"d": 1,
		"e": 1,
	}

	results, frequencies := ProcessString(str)
  assert.EqualValues(t, expected, results)

	if !reflect.DeepEqual(frequencies, expectedFreq) {
		t.Errorf("frequencies = %v; want %v", frequencies, expectedFreq)
	}
}

