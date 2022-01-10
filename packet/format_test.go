package packet

import "testing"

func TestReadVarInt(t *testing.T) {
	b := []byte{}
	_, _, err := ReadVarInt(b)
	if err == nil {
		t.Errorf("error expected")
	}
	b0 := []byte{1}
	v0, n0, err0 := ReadVarInt(b0)
	if err0 != nil {
		t.Errorf("did not expect any error")
	}
	if n0 != 1 {
		t.Errorf("expected 1 byte, read %d\n", n0)
	}
	if v0 != 1 {
		t.Errorf("expected value 1, read %d\n", v0)
	}

	b1 := []byte{128, 1}
	v1, n1, err1 := ReadVarInt(b1)
	if err1 != nil {
		t.Errorf("did not expect any error")
	}
	if n1 != 2 {
		t.Errorf("expected 2 bytes, read %d\n", n1)
	}
	if v1 != 128 {
		t.Errorf("expected value 128, read %d\n", v1)
	}

	b2 := []byte{128, 128, 128, 1}
	v2, n2, err2 := ReadVarInt(b2)
	if err2 != nil {
		t.Errorf("did not expect any error")
	}
	if n2 != 4 {
		t.Errorf("expected 4 bytes, read %d\n", n1)
	}
	if v2 != 2097152 {
		t.Errorf("expected value 2097152, read %d\n", v2)
	}

}

func TestWriteVarInt(t *testing.T) {
	b0 := WriteVarInt(0)
	if len(b0) > 1 || b0[0] != 0 {
		t.Errorf("expect [0], received [%d]", b0[0])
	}
	a, b, c := ReadVarInt(b0)
	if c == nil && b != 1 && a != 0 {
		t.Errorf("expect no error, 2 bytes read, extracted 2090 received [%s,%d,%d]", c, b, a)
	}

	b1 := WriteVarInt(2097152)
	if len(b1) != 4 || b1[0] != 128 || b1[1] != 128 || b1[2] != 128 || b1[3] != 1 {
		t.Errorf("expect [128, 128, 128, 1], received [%d,%d,%d,%d]", b1[0], b1[1], b1[2], b1[3])
	}
	a, b, c = ReadVarInt(b1)
	if c == nil && b != 4 && a != 2097152 {
		t.Errorf("expect no error, 2 bytes read, extracted 2090 received [%s,%d,%d]", c, b, a)
	}

	b2 := WriteVarInt(2090)
	if len(b2) != 2 || b2[0] != 170 || b2[1] != 16 {
		t.Errorf("expect [128, 10], received [%d,%d]", b2[0], b2[1])
	}
	a, b, c = ReadVarInt(b2)
	if c == nil && b != 2 && a != 2090 {
		t.Errorf("expect no error, 2 bytes read, extracted 2090 received [%s,%d,%d]", c, b, a)
	}
}

func TestRead2BytesInt(t *testing.T) {
	i0 := []byte{0, 1}
	v0 := Read2BytesInt(i0, 0)
	if v0 != 1 {
		t.Errorf("expected 1 received %d", v0)
	}
	i1 := []byte{1, 1}
	v1 := Read2BytesInt(i1, 0)
	if v1 != 257 {
		t.Errorf("expected 257 received %d", v1)
	}
	i2 := []byte{0, 2, 1, 1}
	v2 := Read2BytesInt(i2, 2)
	if v2 != 257 {
		t.Errorf("expected 257 received %d", v2)
	}
	i3 := []byte{176, 89}
	v3 := Read2BytesInt(i3, 0)
	if v3 != 45145 {
		t.Errorf("expected 45145 received %d", v3)
	}
}

func TestWrite2BytesInt(t *testing.T) {
	i0 := 0
	v0 := Write2BytesInt(i0)
	if v0[0] != 0 || v0[1] != 0 {
		t.Errorf("expected [0, 0] received [%d, %d]", v0[0], v0[1])
	}

	i1 := 1024
	v1 := Write2BytesInt(i1)
	if v1[0] != 4 || v1[1] != 0 {
		t.Errorf("expected [4, 0] received [%d, %d]", v1[0], v1[1])
	}

	i2 := 1025
	v2 := Write2BytesInt(i2)
	if v2[0] != 4 || v2[1] != 1 {
		t.Errorf("expected [4, 1] received [%d, %d]", v2[0], v2[1])
	}
}

func TestRead4BytesInt(t *testing.T) {
	b := []byte{0, 0, 0, 1}
	v := Read4BytesInt(b)
	if v != 1 {
		t.Errorf("expected 1, read %d\n", v)
	}
	b1 := []byte{0, 0, 1, 0}
	v1 := Read4BytesInt(b1)
	if v1 != 256 {
		t.Errorf("expected 256, read %d\n", v1)
	}
	b2 := []byte{1, 0, 0, 1}
	v2 := Read4BytesInt(b2)
	if v2 != 16777217 {
		t.Errorf("expected 16777217, read %d\n", v2)
	}
}
