package format

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestReadVarIntTable is a table-driven test for ReadVarInt function
func TestReadVarIntFromBytesTable(t *testing.T) {
	tests := []struct {
		input   []byte
		want    []int
		wantErr bool
	}{
		{[]byte{}, []int{0, 0}, true},
		{[]byte{0x01}, []int{1, 1}, false},
		{[]byte{0x80, 0x01}, []int{128, 2}, false},
		{[]byte{0x80, 0x80, 0x80, 1}, []int{2097152, 4}, false},
		{[]byte{0x80, 0x80, 0x80, 0x80, 0x01}, []int{0, 0}, true},
		{[]byte{0xaa, 0x10}, []int{2090, 2}, false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("input=%v", tt.input), func(t *testing.T) {
			v, n, err := ReadVarIntFromBytes(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.want[0], v)
				assert.Equal(t, tt.want[1], n)
			}
		})
	}
}

// TestWriteVarIntTable is a table-driven test for WriteVarInt function
func TestWriteVarIntTable(t *testing.T) {
	tests := []struct {
		input int
		want  []byte
	}{
		{0, []byte{0}},
		{127, []byte{0x7F}},
		{128, []byte{0x80, 0x01}},
		{129, []byte{0x81, 0x01}},
		{16383, []byte{0xFF, 0x7F}},
		{16384, []byte{0x80, 0x80, 0x01}},
		{2097151, []byte{0xFF, 0xFF, 0x7F}},
		{2097152, []byte{0x80, 0x80, 0x80, 0x01}},
		{268435455, []byte{0xFF, 0xFF, 0xFF, 0x7F}},
		{2090, []byte{0xaa, 0x10}},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("input=%v", tt.input), func(t *testing.T) {
			got, _ := WriteVarInt(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestRead2BytesIntTable is a table-driven test for Read2BytesInt function
func TestRead2BytesIntTable(t *testing.T) {
	tests := []struct {
		input   []byte
		want    int
		wantErr bool
	}{
		{[]byte{0}, 0, true},
		{[]byte{0, 1}, 1, false},
		{[]byte{1, 1}, 257, false},
		{[]byte{176, 89}, 45145, false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("input=%v", tt.input), func(t *testing.T) {
			got, err := Read2BytesInt(tt.input, 0)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

// TestWrite2BytesIntTable is a table-driven test for Write2BytesInt function
func TestWrite2BytesIntTable(t *testing.T) {
	tests := []struct {
		input int
		want  []byte
	}{
		{0, []byte{0, 0}},
		{1024, []byte{4, 0}},
		{1025, []byte{4, 1}},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("input=%v", tt.input), func(t *testing.T) {
			got := Write2BytesInt(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestRead4BytesIntTable is a table-driven test for Read4BytesInt function
func TestRead4BytesIntTable(t *testing.T) {
	tests := []struct {
		input   []byte
		want    uint32
		wantErr bool
	}{
		{[]byte{0, 0}, 0, true},
		{[]byte{0, 0, 0, 1}, 1, false},
		{[]byte{0, 0, 1, 0}, 256, false},
		{[]byte{1, 0, 0, 1}, 16777217, false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("input=%v", tt.input), func(t *testing.T) {
			got, err := Read4BytesInt(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
