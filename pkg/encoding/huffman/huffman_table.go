package huffman

import (
	"container/heap"
	"errors"
	"fmt"

	"github.com/Shresth72/ascii/pkg/assert"
	"github.com/Shresth72/ascii/pkg/encoding/buffer"
	"github.com/Shresth72/ascii/pkg/encoding/utils"
)

const HUFFMAN_ENCODE_LENGTH = 6

var HuffmanTooLarge = errors.New("huffman tree is too large")

type HuffmanEncodingTable struct {
	Bits   []byte
	Len    int
	BitMap map[int][]byte
}

func NewHuffmanTable() *HuffmanEncodingTable {
	return &HuffmanEncodingTable{
		Bits:   make([]byte, 24, 24),
		Len:    0,
		BitMap: make(map[int][]byte),
	}
}

func (h *HuffmanEncodingTable) Left() {
	h.Bits[h.Len] = 0
	h.Len++
}

func (h *HuffmanEncodingTable) Right() {
	h.Bits[h.Len] = 1
	h.Len++
}

func (h *HuffmanEncodingTable) Pop() {
	h.Len--
}

func (h *HuffmanEncodingTable) Encode(value int) {
	encodingValue := make([]byte, h.Len, h.Len)
	copy(encodingValue, h.Bits)
	h.BitMap[value] = encodingValue
}

func (h *HuffmanEncodingTable) String() string {
	out := fmt.Sprintf("encoding table(%d): ", h.Len)
	for i := range h.Len {
		out += fmt.Sprintf("%d", h.Bits[i])
	}
	out += "\n"

	for k, v := range h.BitMap {
		out += fmt.Sprintf("  %d => ", k)
		for _, bit := range v {
			out += fmt.Sprintf("%d", bit)
		}
		out += "\n"
	}

	return out
}

func CalculateHuffman(freq buffer.Frequency) *Huffman {
	nodes := make(PriorityQueue, freq.Length(), freq.Length())
	for i, p := range freq.Points {
		nodes[i] = fromFreq(p)
	}
	heap.Init(&nodes)

	count := 1
	for len(nodes) > 1 {
		a := heap.Pop(&nodes).(*huffmanNode)
		b := heap.Pop(&nodes).(*huffmanNode)
		heap.Push(&nodes, join(a, b))
		count += 2
	}

	head := heap.Pop(&nodes).(*huffmanNode)
	fmt.Printf("\n%s\n", head.debug(0))

	size := count * HUFFMAN_ENCODE_LENGTH
	encoding := make([]byte, size, size)
	table := NewHuffmanTable()

	encodeTree(head, table, encoding, 0)

	return &Huffman{
		DecodingTree:  encoding,
		EncodingTable: table.BitMap,
	}
}

func encodeTree(node *huffmanNode, table *HuffmanEncodingTable, data []byte, idx int) int {
	if node == nil {
		return idx
	}

	assert.Assert(idx+5 < len(data), "idx will exceed the bounds of the huffman array during encoding")
	leftIdx := idx + HUFFMAN_ENCODE_LENGTH

	utils.Write16(data, idx, node.value)
	utils.Write16(data, idx+2, leftIdx)

	table.Left()
	rightIdx := encodeTree(node.left, table, data, leftIdx)
	table.Pop()

	utils.Write16(data, idx+4, rightIdx)

	table.Right()
	doneIdx := encodeTree(node.right, table, data, rightIdx)
	table.Pop()

	if node.isLeaf() {
		utils.Write16(data, idx+2, 0)
		utils.Write16(data, idx+4, 0)
		table.Encode(node.value)
	}

	return doneIdx
}
