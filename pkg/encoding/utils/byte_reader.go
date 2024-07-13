package utils

import "github.com/Shresth72/ascii/pkg/assert"

type ByteIteratorResult struct {
  Value int
  Done  bool
}

type ByteIterator interface {
  Next() ByteIteratorResult
}

// 16 Bit Reader
type SixteenBitIterator struct {
  buffer []byte 
  idx    int 
  res    ByteIteratorResult
}

func New16BitIterator(buf []byte) *SixteenBitIterator {
  assert.Assert(len(buf)&0x1 == 0, "buf size must be even")
  
  return &SixteenBitIterator{
    buffer: buf,
    idx: 0,
    res: ByteIteratorResult{
      Value: 0,
      Done: false,
    },
  }
}

func Read16(buf []byte, offset int) int {
  assert.Assert(len(buf) > offset + 1, "cannot read outside the buffer")

  hi := int(buf[offset])
  lo := int(buf[offset + 1])
  return hi<<8 | lo
}

func (i *SixteenBitIterator) Next() ByteIteratorResult {
  assert.Assert(!i.res.Done, "Next called on exhausted iterator")

  value := Read16(i.buffer, i.idx)
  i.idx += 2

  i.res.Done = i.idx == len(i.buffer)
  i.res.Value = value
  
  return i.res
}

// 8 Bit Reader
type EightBitIterator struct {
  buffer []byte 
  idx    int
  res    ByteIteratorResult
}

func New8BitIterator(buf []byte) *EightBitIterator {
  return &EightBitIterator{
    buffer: buf,
    idx: 0,
    res: ByteIteratorResult{
      Value: 0,
      Done: false,
    },
  }
}

func (i *EightBitIterator) Next() ByteIteratorResult {
  assert.Assert(!i.res.Done, "Next called on exhausted iterator")

  val := i.buffer[i.idx]
  i.idx++

  i.res.Done = i.idx == len(i.buffer)
  i.res.Value = int(val)

  return i.res
}
