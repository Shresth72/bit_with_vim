package huffman

import (
	"fmt"
	"strings"

	"github.com/Shresth72/ascii/pkg/encoding/buffer"
)

type huffmanNode struct {
	value int
	count int
	left  *huffmanNode
	right *huffmanNode
}

func (h *huffmanNode) isLeaf() bool {
	return h.left == nil && h.right == nil
}

func (h *huffmanNode) String() string {
	if h == nil {
		return "nil"
	}
	return fmt.Sprintf("node(%d): %d", h.count, h.value)
}

func (h *huffmanNode) debug(indent int) string {
	indentStr := strings.Repeat(" ", indent*2)
	if h == nil {
		return fmt.Sprintf("%s-> nil\n", indentStr)
	}

	return fmt.Sprintf("%s->%s\n", indentStr, h.String()) +
		h.left.debug(indent+1) +
		h.right.debug(indent+1)
}

func fromValue(value int) *huffmanNode {
	return &huffmanNode{
		value: value,
		count: 1,
		left:  nil,
		right: nil,
	}
}

func join(a, b *huffmanNode) *huffmanNode {
	return &huffmanNode{
		value: 0,
		count: a.count + b.count,
		left:  a,
		right: b,
	}
}

func fromFreq(freq *buffer.FreqPoint) *huffmanNode {
	return &huffmanNode{
		value: freq.Val,
		count: freq.Count,
		left:  nil,
		right: nil,
	}
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*huffmanNode

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].count < pq[j].count
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x any) {
	item := x.(*huffmanNode)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*pq = old[0 : n-1]
	return item
}
