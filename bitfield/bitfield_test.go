package bitfield

import (
	"testing"
)

func Test1(t *testing.T) {
	bf := New(5, 0, 2)
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

	bf = New(180, 110, 113, 11)
	bf.Set(2)
	bf.Clear(2)
	if bf.OnesCount() != 3 {
		t.Error("Should be 3")
	}
}

func Test2(t *testing.T) {
	bf1 := New(11)
	bf2 := New(12)
	res, err := bf1.And(bf2)
	if err == nil {
		t.Error("Should be false")
	}
	bf2 = New(11, 1, 10)
	bf1.Set(1)
	bf1.Set(7)
	res, err = bf1.And(bf2)
	if err != nil || res.OnesCount() != 1 {
		t.Error("Should be 1")
	}
	bf2.Set(7)
	res, err = bf1.And(bf2)
	if err != nil || !res.Equal(bf1) {
		t.Error("Should be 2")
	}
	bf2.Clear(10)
	if !bf2.Equal(bf1) {
		t.Error("Should be equal")
	}
}
