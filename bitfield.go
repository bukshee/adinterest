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
		// we avoid returning error in order to be chainable
		return nil
	}
	ret := BitField{
		data: make(bitFieldData, 1+len/64),
		len:  len,
	}
	return &ret
}

// Len returns the number of bits the BitField holds
func (bf *BitField) Len() int {
	return bf.len
}

func (bf *BitField) posToOffset(pos int) (index int, offset int) {
	for pos < 0 {
		pos += bf.len
	}
	pos %= bf.len
	index = pos / 64
	offset = pos % 64
	return index, offset
}

// clearEnd zeroes the bits beyond Len(): the underlying BitField64
// allocates space in 64bit increments and Len() might be smaller than
// the space allocated: it needs to be kept zeroed at all times to be
// consistent
func (bf *BitField) clearEnd() *BitField {
	index, offset := bf.posToOffset(bf.Len() - 1)
	// point to after the last element:
	for i := offset + 1; i < 64; i++ {
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

// Copy returns a copy of the bitfield or nil if copy failed
func (bf *BitField) Copy() *BitField {
	bfNew := BitField{
		data: make(bitFieldData, len(bf.data), cap(bf.data)),
		len:  bf.len,
	}
	copy(bfNew.data, bf.data)
	return &bfNew
}