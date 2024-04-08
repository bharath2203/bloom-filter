package bloomfilter

import "testing"

func Test_bitset(t *testing.T) {

	bs := NewBitSet(1000)

	if bs.Get(10) {
		t.Fatalf("bitset test failed")
	}

	bs.Set(100)
	if !bs.Get(100) {
		t.Fatalf("bitset test failed")
	}

	bs.Set(999)
	if !bs.Get(999) {
		t.Fatalf("bitset test failed")
	}

	bs.Set(63)
	if !bs.Get(63) {
		t.Fatalf("bitset test failed")
	}

}
