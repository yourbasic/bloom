package bloom

import (
	"bytes"
	"encoding/gob"
)

const binVer0 = 0

// BUG(korthaj): Filters cannot be moved between little-endian
// and big-endian machines.

// Data to be included in binary representation of Filter.
type marshalFilter struct {
	Version int
	Data    []uint64
	Lookups int
	Count   int64
}

// MarshalBinary returns a binary representation of the filter.
//
// This method implements the encoding.BinaryMarshaler interface.
// The packages encoding/gob, encoding/json, and encoding/xml
// all check for this interface.
//
// Filters cannot be moved between little-endian and big-endian machines.
func (f *Filter) MarshalBinary() ([]byte, error) {
	mf := marshalFilter{
		Version: binVer0,
		Data:    f.data,
		Lookups: f.lookups,
		Count:   f.count,
	}
	var b bytes.Buffer
	err := gob.NewEncoder(&b).Encode(mf)
	return b.Bytes(), err
}

// UnmarshalBinary imports binary data created by MarshalBinary
// into an empty filter. If the filter is not empty, all previous
// entries are overwritten.
func (f *Filter) UnmarshalBinary(data []byte) error {
	var mf marshalFilter
	err := gob.NewDecoder(bytes.NewBuffer(data)).Decode(&mf)
	if err != nil {
		return err
	}
	f.data = mf.Data
	f.lookups = mf.Lookups
	f.count = mf.Count
	return nil
}
