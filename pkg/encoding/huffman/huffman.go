package huffman

type Huffman struct {
	DecodingTree  []byte
	EncodingTable map[int][]byte
}
