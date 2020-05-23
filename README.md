![Coverage](https://img.shields.io/badge/coverage-100%25-green)

# bitfield
A library for working with bits.

## Description
Create and manage bits in a flexible way: Specify how many bits you want, address, set or clear individual bits by position, do binary operations like AND, OR, NOT, XOR.

## Implementation details
Package bitfield is slice of bitfield64-s to make it possible to store more
than 64 bits. Most functions are chainable, positions outside the [0,len) range
will get the modulo treatment, so Get(len) will return the 0th bit, Get(-1) will
return the last bit: Get(len-1)

See test file for usage.
