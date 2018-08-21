package keyset

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestByteIndex(t *testing.T) {
	bi := &ByteIndex{}

	// Add several IDs (in []byte form)
	id := make([]byte, 8)
	rand.Read(id)
	bi.Add(id)

	id2 := make([]byte, 8)
	rand.Read(id2)
	bi.Add(id2)
	bi.Add(id2) // No-op

	id3 := make([]byte, 8)
	rand.Read(id3)
	bi.Add(id3)

	if !bi.Contains(id2) {
		t.Error("ID2 not saved")
	}

	bi.Remove(id2)

	if bi.Contains(id2) {
		t.Error("ID2 not removed")
	}
}

func TestByteIndexUint64(t *testing.T) {
	bi := &ByteIndex{}

	// Add several IDs
	id := uint64(rand.Int63())
	bi.AddUint64(id)

	id2 := uint64(rand.Int63())
	bi.AddUint64(id2)
	bi.AddUint64(id2) // No-op

	id3 := uint64(rand.Int63())
	bi.AddUint64(id3)

	if !bi.ContainsUint64(id2) {
		t.Error("ID2 not saved")
	}

	bi.RemoveUint64(id2)

	if bi.ContainsUint64(id2) {
		t.Error("ID2 not removed")
	}
}

func TestByteIndexMarshal(t *testing.T) {

	bi := &ByteIndex{}
	for i := 0; i < 10; i++ {
		id := make([]byte, 8)
		rand.Read(id)
		bi.Add(id)
	}

	data, err := bi.MarshalToByte()
	if err != nil {
		t.Error(err)
	}

	bi2 := &ByteIndex{}
	err = bi.UnmarshalFromByte(data)
	if err != nil {
		t.Error(err)
	}

	for _, b := range *bi {
		for _, b2 := range *bi2 {
			if !testEq(b, b2) {
				t.Error("Error encoding/decoding")
			}
		}
	}

}

func BenchmarkByteIndexMarshal(b *testing.B) {
	// Run two sizes, 1k and 100k
	for _, x := range []int{1000, 100000} {

		bi := &ByteIndex{}
		for i := 0; i < x; i++ {
			id := make([]byte, 8)
			rand.Read(id)
			bi.Add(id)
		}

		b.ResetTimer()

		b.Run(fmt.Sprintf("Byte%d", x), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				data, err := bi.MarshalToByte()
				if err != nil {
					b.Error(err)
				}

				bi = &ByteIndex{}
				err = bi.UnmarshalFromByte(data)
				if err != nil {
					b.Error(err)
				}

				if len(*bi) != x {
					b.Error("Error decoding")
				}
			}
		})
	}
}

func testEq(a, b []byte) bool {
	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
