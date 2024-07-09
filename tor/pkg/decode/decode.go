package decode

import (
	"fmt"
	"io"
	"strconv"
)

func DecodeBencode(bencodedString string, start int) (interface{}, int, error) {
if start == len(bencodedString) {
    return nil, start, io.ErrUnexpectedEOF
  }

  i := start
  switch {
  case bencodedString[i] == 'l':
    return DecodeList(bencodedString, i)
  case bencodedString[i] == 'i':
    return DecodeInteger(bencodedString, i)
  case bencodedString[i] >= '0' && bencodedString[i] <= '9':
    return DecodeString(bencodedString, i)
  default:
    return nil, start, fmt.Errorf("unsupported bencode: %q", bencodedString[i])
  }
}

func DecodeString(b string, st int) (interface{}, int, error) {
  colonIndex := st
  for ; colonIndex < len(b) && b[colonIndex] != ':'; colonIndex++ {}

  if colonIndex == len(b) {
    return nil, st, fmt.Errorf("invalid bencoded string")
  }

  length, err := strconv.Atoi(b[st:colonIndex])
  if err != nil {
    return nil, st, err 
  }

  start := colonIndex + 1 
  end := start + length
  if end > len(b) {
    return nil, st, fmt.Errorf("length mismatch: expected at least %d, got %d", length, len(b)-start)
  }

  return b[start:end], end, nil
}

func DecodeInteger(b string, st int) (interface{}, int, error) {
  end := st + 1 
  for ; end < len(b) && b[end] != 'e'; end++ {} 
  
  if end == len(b) {
    return nil, st, fmt.Errorf("invalid bencoded integer")
  }

  decodedInt, err := strconv.Atoi(b[st + 1: end])
  if err != nil {
    return nil, st, err 
  }

  return decodedInt, end + 1, err
}

func DecodeList(b string, st int) (interface{}, int, error) {
  decodedList := make([]interface{}, 0)
  i := st + 1

  for ; i < len(b) && b[i] != 'e'; {
    decodedValue, newSt, err := DecodeBencode(b, i)
    if err != nil {
      return nil, st, err
    }

    decodedList = append(decodedList, decodedValue)
    i = newSt
  }

  if i == len(b) {
    return nil, st, fmt.Errorf("invalid bencoded list")
  }

  return decodedList, i + 1, nil
}
