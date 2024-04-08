package bloomfilter

import (
	"fmt"
	"github.com/spaolacci/murmur3"
	"hash"
)

type BloomFilter interface {
	Add(record string) error
	AddList(records []string) error
	Contains(record string) (bool, error)
}

type bfilter struct {
	bs      *bitset
	hashers []hash.Hash64
	N       uint64
	K       int
}

func NewBloomFilter(n uint64, k int) BloomFilter {
	b := &bfilter{}
	b.N = n
	b.K = k
	b.bs = NewBitSet(n)
	b.hashers = make([]hash.Hash64, k)
	for i := 0; i < k; i++ {
		b.hashers[i] = murmur3.New64WithSeed(uint32(i + 1))
	}
	return b
}

func (b *bfilter) Add(record string) error {

	hArray := b.getHashArray(record)
	if hArray == nil {
		return fmt.Errorf("error while calculating hash of the record, %s", record)
	}

	for _, kthHash := range hArray {
		b.bs.Set(kthHash)
	}

	return nil
}

func (b *bfilter) AddList(records []string) error {
	for index, record := range records {
		if err := b.Add(record); err != nil {
			return fmt.Errorf("errored out while processing %d record", index)
		}
	}
	return nil
}

func (b *bfilter) Contains(record string) (bool, error) {
	hArray := b.getHashArray(record)
	if hArray == nil {
		return false, fmt.Errorf("error while calculating hash of the record, %s", record)
	}

	for _, kthHash := range hArray {
		if !b.bs.Get(kthHash) {
			return false, nil
		}
	}

	return true, nil
}

func (b *bfilter) getHashArray(record string) []uint64 {
	hArray := make([]uint64, 10)
	recordInBytes := []byte(record)
	for index := 0; index <= len(b.hashers); index++ {
		_, err := b.hashers[index].Write(recordInBytes)
		if err != nil {
			return nil
		}
		h := b.hashers[index].Sum64()
		h = h % b.N
		hArray[index] = h
	}
	return hArray
}

type bitset struct {
	cap     uint64
	bitList []uint64
}

func NewBitSet(n uint64) *bitset {
	b := &bitset{}
	b.cap = n
	bitListSize := n/64 + 1
	b.bitList = make([]uint64, bitListSize)
	return b
}

func (b *bitset) Set(n uint64) {
	bitIndex := n / 64
	bitRem := n % 64
	b.bitList[bitIndex] = b.bitList[bitIndex] | (uint64(1) << bitRem)
}

func (b *bitset) Get(n uint64) bool {
	bitIndex := n / 64
	bitRem := n % 64
	if b.bitList[bitIndex]&(uint64(1)<<bitRem) > 0 {
		return true
	}
	return false
}
