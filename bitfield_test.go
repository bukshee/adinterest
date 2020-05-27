package bitfield

import (
	"testing"
)

func Test1(t *testing.T) {
	if New(0).Len() != 0 || New(-2).Len() != 0 {
		t.Error("should be 0")
	}
	if New(65).SetAll().OnesCount() != 65 {
		t.Error("should be 65")
	}
	if New(3).Equal(New(4)) {
		t.Error("should be false")
	}
	if New(3).Set(0).Set(-1).Not().OnesCount() != 1 {
		t.Error("should be 1")
	}
	bf := New(129).Set(0).Set(-1).Clear(123).Clear(-3).Not().Not()
	if bf.OnesCount() != 2 {
		t.Error("should be 2")
	}
	if !bf.Get(0) || !bf.Get(-bf.Len()) || !bf.Get(bf.Len()) {
		t.Error("should be true")
	}
	if bf.And(New(129).Set(0).Set(1)).OnesCount() != 1 {
		t.Error("should be 1")
	}

	if bf.And(New(121)).OnesCount() != 1 {
		t.Error("should be 1")
	}
	if bf.Or(New(121)).OnesCount() != 1 {
		t.Error("should be 1")
	}
	if bf.Xor(New(121)).OnesCount() != 1 {
		t.Error("should be 1")
	}

	bf.Set(73).Set(-2).ClearAll().Set(-1)
	if !bf.Equal(New(129).Set(128)) {
		t.Error("should be equal")
	}
	if bf.Equal(New(129).Not()) {
		t.Error("should be not equal")
	}
	if bf.Get(127) {
		t.Error("should be false")
	}
	if !bf.Get(-1) {
		t.Error("should be true")
	}

	if !New(4).Flip(-1).Equal(New(4).Set(-1)) {
		t.Error("should be equal")
	}
	if !New(4).Flip(-1).Flip(-1).Equal(New(4)) {
		t.Error("should be equal")
	}

	bf2 := bf.Clone()
	if !bf.Equal(bf2) || bf.OnesCount() != bf2.OnesCount() || bf.Len() != bf2.Len() {
		t.Error("should be equal")
	}
	if bf2.Xor(bf2).OnesCount() != 0 {
		t.Error("should be 0")
	}

	bf2 = bf.Clone()
	if !bf2.Set(11).Or(bf).Get(11) {
		t.Error("should be true")
	}
}

func Test2(t *testing.T) {
	if !New(27).Set(0).Set(-1).Resize(65).Get(26) {
		t.Error("should be true")
	}
	if New(65).Set(-1).Resize(45).OnesCount() != 0 {
		t.Error("should be 0")
	}
	if New(65).SetAll().Resize(40).OnesCount() != 40 {
		t.Error("should be 40")
	}
	dest := New(65).Set(4).Resize(0)
	if dest.OnesCount() != 1 || !dest.Get(4) {
		t.Error("should be the same bitfield")
	}

	dest = New(44)
	if New(65).BitCopy(dest) {
		t.Error("should be false")
	}
	dest = New(65)
	if !New(65).Set(-1).BitCopy(dest) {
		t.Error("should be true")
	}
	if !dest.Get(64) || dest.OnesCount() != 1 {
		t.Error("not exact BitCopy")
	}
}

func TestPrivate1(t *testing.T) {
	want := [...][2]int{
		{67, -2}, {121, 121}, {3, -10}, {0, 2}, {0, 4},
	}
	need := [...]int{
		65, 0, 2, 0, 0,
	}

	for i, w := range want {
		res := New(w[0]).posNormalize(w[1])
		if res == need[i] {
			continue
		}
		t.Errorf("New(%d).posNormalize(%d) should map to %d. Got: %d", w[0], w[1], need[i], res)
	}
}

func TestPrivate2(t *testing.T) {
	want := [...][2]int{
		{65, 64}, {3, -1}, {65, -1},
	}
	need := [...][2]int{
		{1, 0}, {0, 2}, {1, 0},
	}
	for i, w := range want {
		ix, p := New(w[0]).posToOffset(w[1])
		if ix == need[i][0] && p == need[i][1] {
			continue
		}
		t.Errorf("New(%d).posToOffset(%d) should map to [%d,%d]. Got: [%d,%d]",
			w[0], w[1], need[i][0], need[i][1], ix, p)
	}
}
