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
