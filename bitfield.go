// Package bitfield is slice of bitfield64-s to make it possible to store more
// than 64 bits. Most functions are chainable
package bitfield

import (
	"errors"

	bf64 "./bitfield64"
)

// BitField is a slice of BitField64-s.
type BitField []bf64.BitField64

// New creates a slice of BitField64
func New(len int) *BitField {
	ret := make(BitField, len/64)
	return &ret
}

func (bf *BitField) posVerify(pos int) error {
	if pos < 0 || pos > len(*bf)*64 {
		return errors.New("wrong position")
	}
	return nil
}

func (bf *BitField) posToOffset(pos int) (index int, offset int, err error) {
	err = bf.posVerify(pos)
	if err != nil {
		return 0, 0, err
	}
	index = pos / 64
	offset = pos % 64
	return index, offset, nil
}

// Set sets a bit to 1 at position pos inside the bit-field
func (bf *BitField) Set(pos int) *BitField {
	index, offset, err := bf.posToOffset(pos)
	if err == nil {
		(*bf)[index] = (*bf)[index].Set(offset)
	}
	return bf
}

// SetAll sets all bits to 1
func (bf *BitField) SetAll() *BitField {
	for i := range *bf {
		(*bf)[i].SetAll()
	}
	return bf
}

// Clear clears the bit at position pos (sets to 0) inside the bit-field
func (bf *BitField) Clear(pos int) *BitField {
	index, offset, err := bf.posToOffset(pos)
	if err == nil {
		(*bf)[index] = (*bf)[index].Clear(offset)
	}
	return bf
}

// Get returns the bit (as a boolean) at position pos
func (bf *BitField) Get(pos int) (bool, error) {
	index, offset, err := bf.posToOffset(pos)
	if err != nil {
		return false, err
	}
	return (*bf)[index].Get(offset), nil
}

// OnesCount returns the number of bits set
func (bf *BitField) OnesCount() int {
	count := 0
	for i := range *bf {
		count += (*bf)[i].OnesCount()
	}
	return count
}

// And ANDs a bitfield to this one and returns the result as a new bitfield
func (bf *BitField) And(bfOther *BitField) *BitField {
	if len(*bf) != len(*bfOther) {
		return bf
	}
	for i := range *bf {
		(*bf)[i] = (*bf)[i].And((*bfOther)[i])
	}
	return bf
}

// Equal tells if two bitfields are equal or not
func (bf *BitField) Equal(bfOther *BitField) bool {
	if len(*bf) != len(*bfOther) {
		return false
	}
	for i := range *bf {
		if (*bf)[i] != (*bfOther)[i] {
			return false
		}
	}
	return true
}
