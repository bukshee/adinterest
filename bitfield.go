/*
Package bitfield is slice of bitfield64-s to make it possible to store more
than 64 bits. Most functions are chainable, positions outside the [0,len) range
will get the modulo treatment, so Get(len) will return the 0th bit, Get(-1) will
return the last bit: Get(len-1)
*/
package bitfield

import (
	bf64 "github.com/bukshee/bitfield64"
)

type bitFieldData []bf64.BitField64

// BitField is a slice of BitField64-s.
type BitField struct {
	data bitFieldData
	len  int
}

// New creates a slice of BitField64 and returns it. Returns nil if len<=0
func New(len int) *BitField {
	if len <= 0 {
		len = 0
	}
	ret := BitField{
		data: make(bitFieldData, 1+len/64),
		len:  len,
	}
	return &ret
}

// Resize resizes the bitfield to newLen in size.
// Returns a newly allocated one, leaves the original intact.
// If newLen < Len() bits are lost at the end.
// If newLen > Len() the newly added bits will be zeroed.
func (bf *BitField) Resize(newLen int) *BitField {
	if newLen < 0 {
		newLen = 0
	}
	ret := New(newLen)
	copy(ret.data, bf.data)
	if newLen < bf.len {
		ret.clearEnd()
	}
	return ret
}

// Len returns the number of bits the BitField holds
func (bf *BitField) Len() int {
	return bf.len
}

func (bf *BitField) posNormalize(pos int) int {
	if bf.len == 0 {
		return 0
	}
	for pos < 0 {
		pos += bf.len
	}
	pos %= bf.len
	return pos
}

func (bf *BitField) posToOffset(pos int) (index int, offset int) {
	pos = bf.posNormalize(pos)
	index = pos / 64
	offset = pos % 64
	return index, offset
}

// clearEnd zeroes the bits beyond Len(): the underlying BitField64
// allocates space in 64bit increments and Len() might be smaller than
// the space allocated: it needs to be kept zeroed at all times to be
// consistent
func (bf *BitField) clearEnd() *BitField {
	const n = 64
	index, offset := bf.Len()/n, bf.Len()%n
	// point to after the last element:
	for i := offset; i < n; i++ {
		bf.data[index] = bf.data[index].Clear(i)
	}
	return bf
}

// Set sets a bit to 1 at position pos inside the bit-field
func (bf *BitField) Set(pos int) *BitField {
	index, offset := bf.posToOffset(pos)
	bf.data[index] = bf.data[index].Set(offset)
	return bf
}

// SetMul sets multiple bits at once
func (bf *BitField) SetMul(pos ...int) *BitField {
	for _, p := range pos {
		bf.Set(p)
	}
	return bf
}

// SetAll sets all bits to 1
func (bf *BitField) SetAll() *BitField {
	for i := range bf.data {
		bf.data[i] = bf.data[i].SetAll()
	}
	return bf.clearEnd()
}

// Clear clears the bit at position pos (sets to 0) inside the bit-field
func (bf *BitField) Clear(pos int) *BitField {
	index, offset := bf.posToOffset(pos)
	bf.data[index] = bf.data[index].Clear(offset)
	return bf
}

// ClearMul clears multiple bits at once
func (bf *BitField) ClearMul(pos ...int) *BitField {
	for _, p := range pos {
		bf.Clear(p)
	}
	return bf
}

// ClearAll sets all bits to 1
func (bf *BitField) ClearAll() *BitField {
	for i := range bf.data {
		bf.data[i] = bf.data[i].ClearAll()
	}
	return bf
}

// Get returns the bit (as a boolean) at position pos
func (bf *BitField) Get(pos int) bool {
	index, offset := bf.posToOffset(pos)
	return bf.data[index].Get(offset)
}

// Flip inverts the bit at position pos
func (bf *BitField) Flip(pos int) *BitField {
	index, offset := bf.posToOffset(pos)
	bf.data[index] = bf.data[index].Flip(offset)
	return bf
}

// OnesCount returns the number of bits set
func (bf *BitField) OnesCount() int {
	count := 0
	for i := range bf.data {
		count += bf.data[i].OnesCount()
	}
	return count
}

// And does a binary AND with bfOther. Modifies the bitfield in place and returns it.
func (bf *BitField) And(bfOther *BitField) *BitField {
	if bf.len != bfOther.len {
		return bf
	}
	for i := range bf.data {
		bf.data[i] = bf.data[i].And(bfOther.data[i])
	}
	return bf
}

// Or does a binary OR with bfOther. Modifies the bitfield in place and returns it.
func (bf *BitField) Or(bfOther *BitField) *BitField {
	if bf.len != bfOther.len {
		return bf
	}
	for i := range bf.data {
		bf.data[i] = bf.data[i].Or(bfOther.data[i])
	}
	return bf
}

// Not does a binary NOT (inverts all bits). Modifies the bitfield in place and returns it.
func (bf *BitField) Not() *BitField {
	for i := range bf.data {
		bf.data[i] = bf.data[i].Not()
	}
	return bf.clearEnd()
}

// Xor does a binary XOR with bfOther. Modifies the bitfield in place and returns it.
func (bf *BitField) Xor(bfOther *BitField) *BitField {
	if bf.len != bfOther.len {
		return bf
	}
	for i := range bf.data {
		bf.data[i] = bf.data[i].Xor(bfOther.data[i])
	}
	return bf.clearEnd()
}

// Equal tells if two bitfields are equal or not
func (bf *BitField) Equal(bfOther *BitField) bool {
	if bf.len != bfOther.len {
		return false
	}
	for i := range bf.data {
		if bf.data[i] != bfOther.data[i] {
			return false
		}
	}
	return true
}

// Clone creates a copy of the bitfield and returns it
func (bf *BitField) Clone() *BitField {
	bfNew := BitField{
		data: make(bitFieldData, len(bf.data), cap(bf.data)),
		len:  bf.len,
	}
	copy(bfNew.data, bf.data)
	return &bfNew
}

// Copy deprecated, use Clone instead: just a rename
func (bf *BitField) Copy() *BitField {
	return bf.Clone()
}

// BitCopy copies the content of the bitfield to dest.
// Returns false if Len()-s differ, true otherwise.
func (bf *BitField) BitCopy(dest *BitField) bool {
	if bf.Len() != dest.Len() {
		return false
	}
	copy(dest.data, bf.data)
	return true
}

// Shift shifts the bitfield by count bits in place and returns it.
// If count is positive it shifts towards higher bit positions;
// If negative it shifts towards lower bit positions.
// Bits exiting at one end are discarded;
// bits entering at the other end are zeroed.
func (bf *BitField) Shift(count int) *BitField {
	if count == 0 {
		return bf
	}
	if count <= -bf.Len() || count >= bf.Len() {
		return bf.ClearAll()
	}

	const n = 64
	if count > 0 {
		ix, delta := count/n, count%n
		for i := len(bf.data) - 1; i >= 0; i-- {
			tmp := bf64.New()
			if i-ix >= 0 {
				tmp = bf.data[i-ix]
			}
			a, b := tmp.Shift2(delta)
			bf.data[i] = a
			if i+1 < len(bf.data) {
				bf.data[i+1] = bf.data[i+1].Or(b)
			}
		}
		bf.clearEnd()
	}
	if count < 0 {
		ix, delta := -count/n, -count%n
		for i := 0; i < len(bf.data); i++ {
			tmp := bf64.New()
			if i+ix < len(bf.data) {
				tmp = bf.data[i+ix]
			}
			a, b := tmp.Shift2(-delta)
			bf.data[i] = a
			if i > 0 {
				bf.data[i-1] = bf.data[i-1].Or(b)
			}
		}
	}
	return bf
}

// Mid returns counts bits from position pos as a new BitField
func (bf *BitField) Mid(pos, count int) *BitField {
	pos = bf.posNormalize(pos)
	if count < 0 {
		count = 0
	}
	return bf.Clone().Shift(-pos).Resize(count)
}

// Left returns count bits in the range of [0,count-1] as a new BitField
func (bf *BitField) Left(count int) *BitField {
	if count > bf.Len() {
		count = bf.Len()
	}
	return bf.Mid(0, count)
}

// Right returns count bits in the range of [63-count,63] as a new BitField
func (bf *BitField) Right(count int) *BitField {
	if count > bf.Len() {
		count = bf.Len()
	}
	return bf.Mid(bf.Len()-count, count)
}

// Append appends 'other' BitField to the end
// A newly created bitfield will be returned
func (bf *BitField) Append(other *BitField) *BitField {
	if other.Len() == 0 {
		return bf.Clone()
	}
	len := bf.Len()
	newLen := len + other.Len()
	ret := other.Resize(newLen).Shift(len)
	return ret.Or(bf.Resize(newLen))
}
