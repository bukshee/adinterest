package bitfield

import (
	"testing"
)

func Test1(t *testing.T) {
	bf := New(5)
	bf.Set(0)
	bf.Set(2)
	bf.Clear(4)
	if bf.OnesCount() != 2 {
		t.Error("Should return 2")
	}
	val, _ := bf.Get(2)
	if val != true {
		t.Error("Should be true")
	}
	val, _ = bf.Get(3)
	if val != false {
		t.Error("Should be false")
	}

	bf = New(180)
	bf.Set(110)
	bf.Set(113)
	bf.Set(11)
	bf.Set(2)
	bf.Clear(2)
	if bf.OnesCount() != 3 {
		t.Error("Should be 3")
	}
}
