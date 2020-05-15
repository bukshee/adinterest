package bitfield64

import (
	"testing"
)

func Test1(t *testing.T) {
	bf1 := New().SetAll().Clear(3)
	if bf1.OnesCount() != 63 {
		t.Error("should be 63")
	}
	if bf1.Not().OnesCount() != 1 {
		t.Error("should be 1")
	}
	if bf1.Get(3) {
		t.Error("should be false")
	}
	bf2 := New().Set(3).Set(1).Clear(5).Or(New().Set(4))

	if bf2.OnesCount() != 3 {
		t.Error("should be 3")
	}
	if bf1.Xor(bf1).OnesCount() != 0 {
		t.Error("should be zero")
	}
	if bf2.And(bf1).OnesCount() != 2 {
		t.Error("should be 2")
	}

	if New().SetAll().ClearAll().OnesCount() != 0 {
		t.Error("should be zero")
	}
}
