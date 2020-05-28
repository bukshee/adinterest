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
	bf := New(129).SetMul(0, -1).ClearMul(123, -3).Not().Not()
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
	if !New(27).SetMul(0, -1).Resize(65).Get(26) {
		t.Error("should be true")
	}
	if New(65).Set(-1).Resize(45).OnesCount() != 0 {
		t.Error("should be 0")
	}
	if New(65).SetAll().Resize(40).OnesCount() != 40 {
		t.Error("should be 40")
	}

	dest := New(44)
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
	given := [...][2]int{
		{67, -2}, {121, 121}, {3, -10}, {0, 2}, {0, 4},
	}
	expected := [...]int{
		65, 0, 2, 0, 0,
	}

	for i, w := range given {
		res := New(w[0]).posNormalize(w[1])
		if res == expected[i] {
			continue
		}
		t.Errorf("New(%d).posNormalize(%d) should map to %d. Got: %d", w[0], w[1], expected[i], res)
	}
}

func TestPrivate2(t *testing.T) {
	given := [...][2]int{
		{65, 64}, {3, -1}, {65, -1},
	}
	expected := [...][2]int{
		{1, 0}, {0, 2}, {1, 0},
	}
	for i, w := range given {
		ix, p := New(w[0]).posToOffset(w[1])
		if ix == expected[i][0] && p == expected[i][1] {
			continue
		}
		t.Errorf("New(%d).posToOffset(%d) should map to [%d,%d]. Got: [%d,%d]",
			w[0], w[1], expected[i][0], expected[i][1], ix, p)
	}
}

func TestShift(t *testing.T) {
	given := [...][2]int{
		{100, 54}, {129, -3}, {6, 2}, {6, 0}, {0, 5}, {193, 192},
		{6, 6}, {6, -6}, {193, -192},
		{193, -1}, {191, -100},
	}
	expected := [...]int{
		46, 126, 4, 6, 0, 1,
		0, 0, 1,
		192, 91,
	}

	for i, in := range given {
		c := New(in[0]).SetAll().Shift(in[1]).OnesCount()
		if c != expected[i] {
			t.Errorf("%d: Shift(%d) had %d bits set, however %d was expected", in[0], in[1], c, expected[i])
		}
	}
}

func TestMid(t *testing.T) {
	if !New(5).SetAll().Mid(1, 1).Equal(New(1).Set(0)) {
		t.Error("should be equal")
	}
	if New(121).SetAll().Mid(-3, 3).OnesCount() != 3 {
		t.Error("should be 3")
	}

	if !New(65).Set(3).Left(3).Equal(New(3)) {
		t.Error("should be equal")
	}

	a := New(60).SetAll().Right(10)
	b := New(10).SetAll()

	if !a.Equal(b) {
		t.Error("should be equal")
	}

	a = New(10).SetAll().Right(11)
	a = New(10).SetAll().Left(11)
	if !a.Equal(b) {
		t.Error("should be equal")
	}
	a = New(10).SetAll().Mid(3, -1)
	if !New(0).Equal(a) {
		t.Error("should be equal")
	}
}

func TestAppend(t *testing.T) {
	// trivial cases
	a := New(0).Append(New(0))
	if !a.Equal(New(0)) {
		t.Error("should be equal")
	}
	a = New(0).Append(New(1))
	if !a.Equal(New(1)) {
		t.Error("should be equal")
	}

	// real cases
	a = New(10).SetAll().Append(New(3))
	if a.Len() != 13 || a.OnesCount() != 10 || a.Right(3).OnesCount() != 0 {
		t.Error("Append is wrong")
	}
}

func TestRotate(t *testing.T) {

	a := New(65).Set(63).Rotate(1)
	if !a.Equal(New(65).Set(64)) {
		t.Error("should be equal")
	}

	a = New(65).Set(0).Rotate(-1)
	if !a.Equal(New(65).Set(64)) {
		t.Error("should be equal")
	}

	const len = 163
	a = New(len).Set(0)
	for i := -len * 2; i < len*2; i++ {
		r := a.Clone().Rotate(i)
		if !r.Equal(New(len).Set(i)) {
			t.Errorf("@%d rotate failed", i)
		}
	}
}
