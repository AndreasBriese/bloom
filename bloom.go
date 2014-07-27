package bloom

import (
	"bufio"
	"bytes"
	"encoding/json"
	"hash/fnv"
	"log"
	"os"
)

// Alternative Hash Function
// http://cr.yp.to/cdb/cdb.txt
// The cdb hash function is ``h = ((h << 5) + h) ^ c'', with a starting
// hash of 5381.

// func (bf bloom) djb2(f []byte) uint64 {
// 	hash := 5381
// 	for _, c := range f {
// 		hash = (hash<<5 +hash) ^ int(c)
// 	}
// 	return uint64(hash)
// }

type bloom interface {
	Add
	Has
}

// getSize
// helperfunction to set boolset-length (=size of database)
func getSize(ui64 uint64) (size uint64, exponent uint64) {
	size = 1
	for size < ui64 {
		size <<= 1
		exponent++
	}
	return size, exponent
}

type bloom64 struct {
	boolSet []bool
	sizeExp uint64
	setLocs uint64
}

func (bf bloom64) fnv64a(b []byte) (l, h uint64) {
	h64 := fnv.New64a()
	h64.Write(b)
	hash64 := h64.Sum64()
	l = hash64 << 32 >> 32
	h = hash64 >> 32
	return l, h
}

type bloom32 struct {
	boolSet []bool
	sizeExp uint64
	setLocs uint64
}

// <--- http://www.cse.yorku.ca/~oz/hash.html
// Berkeley DB Hash (32bit)
func (bf bloom32) smdb(b []byte) (l, h uint64) {
	hash := uint32(0)
	for _, c := range b {
		hash = uint32(c) + (hash << 6) + (hash << 16) - hash
	}
	l = hash << (bf.sizeExp >> 1) >> (bf.sizeExp >> 1)
	h = hash >> (bf.sizeExp >> 1)
	return uint64(l), uint64(h)
}

// New32
// returns a new bloom32/bloom64 bloomfilter
func New(entries, locs int) (bloomfilter bloom) {
	size, exponent := getSize(uint64(entries))
	if size > 1<<32 {
		size, exponent := getSize(uint64(entries))
		bloomfilter = (bloom64{
			boolSet: make([]bool, size),
			sizeExp: exponent,
			setLocs: uint64(locs),
		})
	}
	bloomfilter = (bloom32{
		boolSet: make([]bool, size),
		sizeExp: exponent,
		setLocs: uint64(locs)})
	return bloomfilter
}

// // New32
// // returns a new bloom64 filter
// func New64(entries, locs int) (bloomfilter bloom32) {
// 	size, exponent := getSize(uint64(entries))
// 	bloomfilter = (bloom64{
// 		boolSet: make([]bool, size),
// 		sizeExp: exponent,
// 		setLocs: uint64(locs),
// 	})
// }

// Empty32
// returns a new empty bloom32 filter
func Empty32() (bloomfilter bloom32) {
	return New(1, 1)
}

// Empty64
// returns a new empty bloom64 filter
func Empty64() (bloomfilter bloom64) {
	return New(1<<33, 1)
}

func (bf *bloom32) Add(entry []byte) {
	l, h := bf.smdb(entry)
	for i := uint64(0); i < bf.setLocs; i++ {
		bf.boolSet[(h+i*l)&(1<<bf.sizeExp-1)] = true
	}
}

func (bf bloom32) Has(entry []byte) bool {
	l, h := bf.smdb(entry)
	for i := uint64(0); i < bf.setLocs; i++ {
		switch bf.boolSet[(h+i*l)&(1<<bf.sizeExp-1)] {
		case false:
			return false
		}
	}
	return true
}

func (bf *bloom64) Add(entry []byte) {
	l, h := bf.fnv64a(entry)
	for i := uint64(0); i < bf.setLocs; i++ {
		bf.boolSet[(h+i*l)&(1<<bf.sizeExp-1)] = true
	}
}

func (bf bloom64) Has(entry []byte) bool {
	l, h := bf.fnv64a(entry)
	for i := uint64(0); i < bf.setLocs; i++ {
		switch bf.boolSet[(h+i*l)&(1<<bf.sizeExp-1)] {
		case false:
			return false
		}
	}
	return true
}

// func (bf bloom32) JSONMarshal() []byte {
// 	bloomImEx := bloomJSONImExport{}
// 	bloomImEx.SetLocs = bf.setLocs
// 	bloomImEx.FilterSet = bf.ToBytes()
// 	data, err := json.Marshal(bloomImEx)
// 	if err != nil {
// 		log.Fatal("json.Marshal failed: ", err)
// 	}
// 	return data
// }

// func (bf *bloom32) JSONUnmarshal(dbData []byte) {
// 	bloomImEx := bloomJSONImExport{}
// 	json.Unmarshal(dbData, &bloomImEx)
// 	buf := bytes.NewBuffer(bloomImEx.FilterSet)
// 	(*bf).ToBool(buf)
// 	(*bf).setLocs = bloomImEx.SetLocs
// 	_, exponent := getSize(uint64(len(bf.boolSet)))
// 	(*bf).sizeExp = exponent
// }

type bloomJSONImExport struct {
	FilterSet []byte
	SetLocs   uint64
}

func (bf bloom) JSONMarshal() []byte {
	bloomImEx := bloomJSONImExport{}
	bloomImEx.SetLocs = bf.setLocs
	bloomImEx.FilterSet = bf.ToBytes()
	data, err := json.Marshal(bloomImEx)
	if err != nil {
		log.Fatal("json.Marshal failed: ", err)
	}
	return data
}

func (bf *bloom) JSONUnmarshal(dbData []byte) {
	bloomImEx := bloomJSONImExport{}
	json.Unmarshal(dbData, &bloomImEx)
	buf := bytes.NewBuffer(bloomImEx.FilterSet)
	(*bf).ToBool(buf)
	(*bf).setLocs = bloomImEx.SetLocs
	_, exponent := getSize(uint64(len(bf.boolSet)))
	(*bf).sizeExp = exponent
}

func (bf *bloom) ToBytes() []byte {
	var (
		outBytes bytes.Buffer
		calc     = []uint8{128, 64, 32, 16, 8, 4, 2, 1}
		lc       = len(calc)
		l        = len(bf.boolSet)
		v        uint8
		bools    = make([]bool, 8)
	)

	if l%lc != 0 {
		log.Fatalf("Len([]bool) must be multiple of %v but is %v !", lc, l)
	}

	for i := 0; i < l; i += lc {
		bools = bf.boolSet[i : i+lc]
		v = uint8(0)
		for i, b := range bools {
			if b {
				v += calc[i]
			}
		}
		outBytes.WriteByte(v)
	}
	return outBytes.Bytes()
}

func (bf *bloom) ToBool(buf *bytes.Buffer) {
	var calc = []uint8{128, 64, 32, 16, 8, 4, 2, 1}
	hlpBool := make([]bool, (8 * (*buf).Len()))
	index := 0
	for {
		if r, err := buf.ReadByte(); err != nil {
			break
		} else {
			for _, c := range calc {
				if r/c > 0 {
					hlpBool[index] = true
					r -= c
				}
				index++
			}

		}
	}
	(*bf).boolSet = hlpBool[:index]
}

// func (bf *bloom64) ToBytes() []byte {
// 	var (
// 		outBytes bytes.Buffer
// 		calc     = []uint8{128, 64, 32, 16, 8, 4, 2, 1}
// 		lc       = len(calc)
// 		l        = len(bf.boolSet)
// 		v        uint8
// 		bools    = make([]bool, 8)
// 	)

// 	if l%lc != 0 {
// 		log.Fatalf("Len([]bool) must be multiple of %v but is %v !", lc, l)
// 	}

// 	for i := 0; i < l; i += lc {
// 		bools = bf.boolSet[i : i+lc]
// 		v = uint8(0)
// 		for i, b := range bools {
// 			if b {
// 				v += calc[i]
// 			}
// 		}
// 		outBytes.WriteByte(v)
// 	}
// 	return outBytes.Bytes()
// }

// func (bf *bloom64) ToBool(buf *bytes.Buffer) {
// 	var calc = []uint8{128, 64, 32, 16, 8, 4, 2, 1}
// 	hlpBool := make([]bool, (8 * (*buf).Len()))
// 	index := 0
// 	for {
// 		if r, err := buf.ReadByte(); err != nil {
// 			break
// 		} else {
// 			for _, c := range calc {
// 				if r/c > 0 {
// 					hlpBool[index] = true
// 					r -= c
// 				}
// 				index++
// 			}

// 		}
// 	}
// 	(*bf).boolSet = hlpBool[:index]
// }
