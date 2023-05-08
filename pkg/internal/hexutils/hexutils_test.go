package hexutils

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

type unmarshalTest struct {
	input   string
	want    []byte
	wantErr error // if set, decoding must fail
}

type marshalTest struct {
	input []byte
	want  string
}

var (
	decodeBytesTests = []unmarshalTest{
		// invalid
		{input: ``, wantErr: fmt.Errorf("empty hex string")},
		{input: `0`, wantErr: fmt.Errorf("failed to decode string [0]")},
		{input: `0x0`, wantErr: fmt.Errorf("failed to decode string [0]")},
		{input: `0x023`, wantErr: fmt.Errorf("failed to decode string [023]")},
		{input: `0xxx`, wantErr: fmt.Errorf("failed to decode string [xx]")},
		{input: `0x01zz01`, wantErr: fmt.Errorf("failed to decode string [01zz01]")},
		// valid
		{input: `0x`, want: []byte{}},
		{input: `0X`, want: []byte{}},
		{input: `0x02`, want: []byte{0x02}},
		{input: `0X02`, want: []byte{0x02}},
		{input: `0xffffffffff`, want: []byte{0xff, 0xff, 0xff, 0xff, 0xff}},
		{
			input: `0xffffffffffffffffffffffffffffffffffff`,
			want:  []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		},
		{input: `00`, want: []byte{0x00}},
		{input: `02`, want: []byte{0x02}},
		{input: `ffffffffff`, want: []byte{0xff, 0xff, 0xff, 0xff, 0xff}},
		{
			input: `ffffffffffffffffffffffffffffffffffff`,
			want:  []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		},
	}

	encodeBytesTests = []marshalTest{
		{[]byte{}, "0x"},
		{[]byte{0}, "0x00"},
		{[]byte{0, 0, 1, 2}, "0x00000102"},
	}
)

func TestDecode(t *testing.T) {
	for _, test := range decodeBytesTests {
		t.Run(test.input, func(t *testing.T) {
			dec, err := Decode(test.input)
			if !reflect.DeepEqual(err, test.wantErr) {
				t.Fatalf("unexpected error\nexpected: %v\nactual:   %v", test.wantErr, err)
			}
			if !bytes.Equal(test.want, dec) {
				t.Errorf("unexpected result\nexpected: %v\nactual:   %v", test.want, dec)
			}
		})
	}
}

func TestEncode(t *testing.T) {
	for _, test := range encodeBytesTests {
		enc := Encode(test.input)
		if enc != test.want {
			t.Errorf("input %x: wrong encoding %s", test.input, enc)
		}
	}
}
