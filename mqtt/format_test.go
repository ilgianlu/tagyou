package mqtt

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

	b1 := WriteVarInt(2097152)
	if len(b1) != 4 || b1[3] != 128 || b1[2] != 128 || b1[1] != 128 || b1[0] != 1 {
		t.Errorf("expect [128, 128, 128, 1], received [%d,%d,%d,%d]", b1[3], b1[2], b1[1], b1[0])
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
