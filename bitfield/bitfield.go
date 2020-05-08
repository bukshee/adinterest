package bitfield

import (
	"errors"
	"math/bits"
)

// BitField type utilizing the power of 64bit CPUs
type BitField struct {
	size int
	data []uint64
}

// New returns a zeroed (all false) bit-field that can store size elements
func New(size int) BitField {
	var bf BitField
	bf.size = size
	count := (size + 64) / 64
	bf.data = make([]uint64, count)
	return bf
}

func (bf BitField) posVerify(pos int) error {
	if pos < 0 || pos > bf.size {
		return errors.New("Wrong position")
	}
	return nil
}

func (bf BitField) posToOffset(pos int) (index int, offset int, err error) {
	err = bf.posVerify(pos)
	if err != nil {
		return 0, 0, err
	}
	index = pos / 64
	offset = pos % 64
	return index, offset, nil
}

// Set sets a bit to 1 at position pos inside the bit-field
func (bf BitField) Set(pos int) error {
	index, offset, err := bf.posToOffset(pos)
	if err != nil {
		return err
	}
	bf.data[index] |= (1 << uint64(offset))
	return nil
}

// Clear clears the bit at position pos (sets to 0) inside the bit-field
func (bf BitField) Clear(pos int) error {
	index, offset, err := bf.posToOffset(pos)
	if err != nil {
		return err
	}
	bf.data[index] &= ^(1 << uint64(offset))
	return nil
}

// Get returns the bit (as a boolean) at position pos
func (bf BitField) Get(pos int) (bool, error) {
	index, offset, err := bf.posToOffset(pos)
	if err != nil {
		return false, err
	}
	data := bf.data[index] & (1 << uint64(offset))
	return data != 0, nil
}

// Size returns the number of bits the bit-field holds
func (bf BitField) Size() int {
	return bf.size
}

// OnesCount returns the number of bits set
func (bf BitField) OnesCount() int {
	count := 0
	for i := 0; i < len(bf.data); i++ {
		count += bits.OnesCount64(bf.data[i])
	}
	return count
}
