package tests_d_test

import (
	"testing"

	"github.com/codecrafters-io/tor/pkg/decode"
	"github.com/stretchr/testify/assert"
)

func TestDecodeString(t *testing.T) {
  tests := []struct {
    input string
    expected interface{}
  }{
    {"5:hello", "hello"},
    {"9:hello1234", "hello1234"},
    {"3:abc", "abc"},
    {"0:", ""},
    {"5:hello123", "hello"},
  }

  for _, test := range tests {
    t.Run(test.input, func(t *testing.T) {
      got, _, err := decode.DecodeString(test.input, 0)
      assert.NoError(t, err, "DecodeString(%s) returned an error", test.input)
      assert.Equal(t, test.expected, got)
    })
  }
}

func TestDecodeInteger(t *testing.T) {
  tests := []struct {
    input string
    expected interface{}
  }{
    {"i1234e", 1234},
    {"i0e", 0},
    {"i-42e", -42},
    {"i123456789e", 123456789},
  }

  for _, test := range tests {
    t.Run(test.input, func(t *testing.T) {
      got, _, err := decode.DecodeInteger(test.input, 0)
      assert.NoError(t, err, "DecodeInteger(%s) returned an error", test.input)
      assert.Equal(t, test.expected, got)
    })
  }
}

func TestDecodeList(t *testing.T) {
  tests := []struct {
    input string
    expected interface{}
  }{
    {"l4:spam4:eggse", []interface{}{"spam", "eggs"}},
    {"l3:abci42ee", []interface{}{"abc", 42}},
    {"l5:helloi-52ee", []interface{}{"hello", -52}},
    {"le", []interface{}{}},
    {"l3:abcl4:list3:xyzee", []interface{}{"abc", []interface{}{"list", "xyz"}}},
  }

  for _, test := range tests {
    t.Run(test.input, func(t *testing.T) {
      got, _, err := decode.DecodeList(test.input, 0)
      assert.NoError(t, err, "DecodeList(%s) returned an error", test.input)
      assert.Equal(t, test.expected, got)
    })
  }
}

func TestDecodeBencode(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"5:hello", "hello"},
		{"i123e", 123},
		{"l4:spam4:eggse", []interface{}{"spam", "eggs"}},
		{"l3:abcl4:list3:xyzee", []interface{}{"abc", []interface{}{"list", "xyz"}}},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			got, _, err := decode.DecodeBencode(test.input, 0)
			assert.NoError(t, err, "DecodeBencode(%s) returned an error", test.input)
			assert.Equal(t, test.expected, got, "DecodeBencode(%s) = %v; want %v", test.input, got, test.expected)
		})
	}
}
