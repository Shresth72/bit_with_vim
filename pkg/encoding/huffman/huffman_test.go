package huffman_test

import (
	"testing"

	"github.com/Shresth72/ascii/pkg/encoding/buffer"
	"github.com/Shresth72/ascii/pkg/encoding/huffman"
	"github.com/Shresth72/ascii/pkg/encoding/utils"
	"github.com/stretchr/testify/require"
)

func getFreq() buffer.Frequency {
	freq := buffer.NewFrequency()
	freq.Freq(utils.New8BitIterator([]byte{
		'A', 'A', 'A',
		'B', 'B',
		'C', 'D',
	}))
	return freq
}

func TestHuffman(t *testing.T) {
	freq := getFreq()
	encodeLen := byte(huffman.HUFFMAN_ENCODE_LENGTH)
	data := huffman.CalculateHuffman(freq)

	require.Equal(t, []byte{
		0, 0, 0, encodeLen, 0, encodeLen * 2,
		0, 'A', 0, 0, 0, 0, // 0
		0, 0, 0, encodeLen * 3, 0, encodeLen * 4,
		0, 'B', 0, 0, 0, 0, // 10
		0, 0, 0, encodeLen * 5, 0, encodeLen * 6,
		0, 'D', 0, 0, 0, 0, // 110
		0, 'C', 0, 0, 0, 0, // 111
	}, data.DecodingTree)
}
