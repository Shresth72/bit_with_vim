package utils

import (
	"errors"

	"github.com/Shresth72/ascii/pkg/assert"
)

var ByteWriterExceedsBuffer = errors.New("writer exceeds underlying bufer")

type ByteWriter interface {
  Write(value int) error
  Len() int
}

func SafeWrite16(buf []byte, offset, value int) bool {
  if len(buf) <= offset + 1 {
    return false
  }

  hi := (value & 0xFF00) >> 8
  lo := value & 0xFF 
  buf[offset] = byte(hi)
  buf[offset + 1] = byte(lo)

  return true
}

func Write16(buf []byte, offset, value int) {
  assert.Assert(len(buf) > offset + 1, "cannot write outside buffer") 

  hi := (value & 0xFF00) >> 8
  lo := value & 0xFF 
  buf[offset] = byte(hi)
  buf[offset + 1] = byte(lo)
}

// 8 Bit writer
type U8Writer struct {
  buf []byte
  offset int
}

func (b *U8Writer) Write(value int) error {
  if b.offset == len(b.buf) {
    return ByteWriterExceedsBuffer
  }

  b.buf[b.offset] = byte(value)
  b.offset++

  return nil
}

func (b *U8Writer) Set(data []byte) {
  assert.Assert(len(data) > 0, "buffer length should 1 or more")
  b.buf = data
  b.offset = 0
}

func (b *U8Writer) Len() int {
  return b.offset
}

// 16 Bit writer
type U16Writer struct {
  buf []byte
  offset int
}

func (b *U16Writer) Set(data []byte) {
  assert.Assert(len(data)&1 == 0 && len(data) > 0, "must be even sized buffer to create u16 writer")
  b.buf = data
  b.offset = 0
}

func (b *U16Writer) Write(value int) error {
  if SafeWrite16(b.buf, b.offset, value) {
    b.offset += 2
    return nil
  }

  return ByteWriterExceedsBuffer
}

func (b *U16Writer) Len() int {
  return b.offset
}
