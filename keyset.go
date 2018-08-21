package keyset

import (
	"bytes"
	"encoding/binary"
	"sort"
)

// ByteIndex is a slice of unique, sorted keys([]byte) such as what an index points to
type ByteIndex [][]byte

// Add a key (ignoring duplicates)
func (v *ByteIndex) Add(key []byte) {
	i := sort.Search(len(*v), func(i int) bool {
		return bytes.Compare((*v)[i], key) >= 0
	})

	if i < len(*v) && bytes.Equal((*v)[i], key) {
		// already added
		return
	}

	*v = append(*v, nil)
	copy((*v)[i+1:], (*v)[i:])
	(*v)[i] = key
}

// AddUint64 to index
func (v *ByteIndex) AddUint64(id uint64) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, id)
	v.Add(b)
}

// Remove a key (if exists)
func (v *ByteIndex) Remove(key []byte) {
	i := sort.Search(len(*v), func(i int) bool {
		return bytes.Compare((*v)[i], key) >= 0
	})

	if i < len(*v) {
		copy((*v)[i:], (*v)[i+1:])
		(*v)[len(*v)-1] = nil
		*v = (*v)[:len(*v)-1]
	}
}

// RemoveUint64 from index
func (v *ByteIndex) RemoveUint64(id uint64) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, id)
	v.Remove(b)
}

// Contains returns true if key is found
func (v *ByteIndex) Contains(key []byte) bool {
	i := sort.Search(len(*v), func(i int) bool {
		return bytes.Compare((*v)[i], key) >= 0
	})

	return (i < len(*v) && bytes.Equal((*v)[i], key))
}

// ContainsUint64 returns true if key is found
func (v *ByteIndex) ContainsUint64(id uint64) bool {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, id)
	return v.Contains(b)
}

func (v *ByteIndex) MarshalToByte() (index []byte, err error) {
	// Create a byte array large enough, then fill it with each value
	index = make([]byte, len((*v))*8)
	for i, id := range *v {
		copy(id[0:], index[i:i+8])
		// index[i] = id[0]
		// index[i+1] = id[1]
		// index[i+2] = id[2]
		// index[i+3] = id[3]
		// index[i+4] = id[4]
		// index[i+5] = id[5]
		// index[i+6] = id[6]
		// index[i+7] = id[7]
	}
	return
}

func (v *ByteIndex) UnmarshalFromByte(index []byte) error {
	(*v) = make([][]byte, len(index)/8)
	for i := 0; i < len(index); i += 8 {
		// (*v) = append((*v), index[i:i+8])
		(*v)[i/8] = index[i : i+8]
	}
	return nil
}

// TODO Using append we can build from an arbitrary stream, the problem is that
// if the index is too large to fit in memory, we are wasting time anyway
// func (v *keyList) UnmarshalFromStream(index []byte) error {
// 	for i := 0; i < len(index); i += 8 {
// 		(*v) = append((*v), index[i:i+8])
// 	}
// 	return nil
// }
