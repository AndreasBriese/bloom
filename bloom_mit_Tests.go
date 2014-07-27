package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/reddragon/bloomfilter"
	"github.com/vova616/xxhash"
	"hash"
	// "hash/crc64"
	"hash/fnv"
	"log"
	"os"
	"time"
)

// var prime32_2 = big.NewInt(2881947347)
var stt time.Time

func IF(bools *[]bool) bool {
	for _, v := range *bools {
		if v {
			return true
		}
	}
	return false
}

func COMP(bools *[]bool) bool {
	r := true
	for _, v := range *bools {
		r = r && v
	}
	return r
}

func SWITCH(bools *[]bool) bool {
	for _, v := range *bools {
		switch v {
		case true:
			return true
		}
	}
	return false
}

func UInt32ToByteSlice(ui32 uint32) []byte {
	return []byte{byte(ui32 >> 24), byte(ui32 >> 16), byte(ui32 >> 8), byte(ui32)}
}

func UInt64ToByteSlice(ui64 uint64) []byte {
	return []byte{byte(ui64 >> 56), byte(ui64 >> 48), byte(ui64 >> 40), byte(ui64 >> 32), byte(ui64 >> 24), byte(ui64 >> 16), byte(ui64 >> 8), byte(ui64)}
}

func swap64(ui64 uint64) (r uint64) {
	return ui64<<32>>56 + 256*ui64<<40>>56 + 65536*ui64<<48>>56 + 16777216*ui64<<56>>56
}

func main() {
	// N := 3000000 // 1 Milliarde
	// fmt.Println("IF")
	// for i := 3; i < N; i *= 10 {
	// 	bools := make([]bool, N)
	// 	bools[i] = true
	// 	stt = time.Now()
	// 	IF(&bools)
	// 	dur := time.Since(stt).Nanoseconds()
	// 	fmt.Println(i, dur, "ns (", dur/int64(i), "ns/op )")
	// }
	// fmt.Println("COMP")
	// for i := 3; i < N; i *= 10 {
	// 	bools := make([]bool, N)
	// 	bools[i] = true
	// 	stt = time.Now()
	// 	COMP(&bools)
	// 	dur := time.Since(stt).Nanoseconds()
	// 	fmt.Println(i, dur, "ns (", dur/int64(i), "ns/op )")
	// }
	// fmt.Println("SWITCH")
	// for i := 3; i < N; i *= 10 {
	// 	bools := make([]bool, N)
	// 	bools[i] = true
	// 	stt = time.Now()
	// 	SWITCH(&bools)
	// 	dur := time.Since(stt).Nanoseconds()
	// 	fmt.Println(i, dur, "ns (", dur/int64(i), "ns/op )")
	// }

	// var bigBool = make([]bool, 80)
	// for i := 0; i < len(bigBool); i++ {
	// 	bigBool[i] = rand.Intn(2) == 1
	// }

	// e := ToBytes(bigBool)
	// buf := bytes.NewBuffer(e)
	// eb := ToBool(buf)

	// st := time.Now()
	// for i, b := range eb {
	// 	if bigBool[i] != b {
	// 		fmt.Print(i, "F ")
	// 	}
	// }
	// fmt.Println("toCheck: ", time.Since(st))

	// st := time.Now()
	// for i, _ := range e {
	// 	if i > len(e)-10 {
	// 		break
	// 	}
	// 	hashr.Reset()
	// 	hashr.Write(e[i : i+9])
	// 	_ = hashr.Sum32()
	// 	hashr.Reset()
	// 	hashr.Write(e[i : i+9])
	// 	h := hashr.Sum32()
	// 	fmt.Println(h)

	// }
	// fmt.Println("toHash: ", time.Since(st).Nanoseconds()/int64(2*len(e)), "ns/op  N=", (2 * len(e)))
	// c = "abcdefghijklmnopqrstuvwxyzüöäABCDEFGHIJKLMNOPQRSTUVWXYZÜÖÄ1234567890!§$%&/()"

	file, err := os.Open("words.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	n := 10000 //10**7
	wordlist := make([]string, n)
	wordlist2 := make([]string, n)
	for i, _ := range wordlist {
		if scanner.Scan() {
			wordlist[i] = scanner.Text()
		}
		if scanner.Scan() {
			wordlist2[i] = scanner.Text()
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// for i := 0; i < n; i++ {
	// 	if reader.Scan() {
	// 		word := reader.Bytes()
	// 		for ii := 0; ii < i; ii++ {
	// 			if string(word) == string(wordlist[ii]) {
	// 				fmt.Println(i, string(word), string(wordlist[ii]))
	// 				i--
	// 				continue
	// 			}
	// 		}
	// 		wordlist[i] = word
	// 	} else {
	// 		break
	// 	}
	// }

	// list := []uint32{}
	// c := 0
	// for i := 0; i < n; i++ {
	// 	in := false
	// 	// h32 := crc32.NewIEEE()
	// 	h32 := fnv.New32a()
	// 	h32.Write([]byte(wordlist[i]))
	// 	sum := h32.Sum32()
	// 	for _, v := range list {
	// 		if sum == v {
	// 			fmt.Println(v, wordlist[i])
	// 			c++
	// 			in = true
	// 		}
	// 	}
	// 	if !in {
	// 		list = append(list, sum)
	// 	}
	// }
	// fmt.Println("...", c, len(list))

	st := time.Now()
	for i := 0; i < n; i++ {
		djb2([]byte(wordlist[i]))
	}
	fmt.Println("djb2   ", time.Since(st).Nanoseconds()/int64(n), "ns/op  N=", n)

	st = time.Now()
	for i := 0; i < n; i++ {
		djbCDB([]byte(wordlist[i]))
	}
	fmt.Println("djbCDB ", time.Since(st).Nanoseconds()/int64(n), "ns/op  N=", n)

	st = time.Now()
	for i := 0; i < n; i++ {
		smdb([]byte(wordlist[i]))
	}
	fmt.Println("smdb   ", time.Since(st).Nanoseconds()/int64(n), "ns/op  N=", n)

	h32 := fnv.New32a()
	st = time.Now()
	for i := 0; i < n; i++ {
		h32.Write([]byte(wordlist[i]))
		h32.Sum32()
		h32.Reset()
	}
	fmt.Println("fnv32a ", time.Since(st).Nanoseconds()/int64(n), "ns/op  N=", n)

	h64 := fnv.New64a()
	st = time.Now()
	for i := 0; i < n; i++ {
		h64.Write([]byte(wordlist[i]))
		h64.Sum64()
		h64.Reset()
	}
	fmt.Println("fnv64a ", time.Since(st).Nanoseconds()/int64(n), "ns/op  N=", n)

	xx32 := xxhash.New(0)
	st = time.Now()
	for i := 0; i < n; i++ {
		xx32.Write([]byte(wordlist[i]))
		xx32.Sum32()
		xx32.Reset()
	}
	fmt.Println("xxHash32", time.Since(st).Nanoseconds()/int64(n), "ns/op  N=", n)

	fmt.Println()

	timeit("djb2", &wordlist, djb2)
	timeit("djbCDB", &wordlist, djbCDB)
	timeit("smdb", &wordlist, smdb)

	// count := 0
	// st := time.Now()
	// for _, w := range wordlist {
	// 	if !rbf.Check([]byte(w)) {
	// 		rbf.Add([]byte(w))
	// 	} else {
	// 		count++
	// 	}
	// }
	// fmt.Println((*rbf).M)
	// fmt.Println("doubletts:", count, "\nBloomfilter Has&Add: ", time.Since(st).Nanoseconds()/int64(n), "ns/op  N=", n)
	// rbf = bloomfilter.NewBloomFilter(3, 1)

	fmt.Println("Reddragon")
	rbf := bloomfilter.NewBloomFilter(3, 50*n)
	count := 0
	countfalse := 0
	st = time.Now()
	for _, w := range wordlist {
		if !rbf.Check([]byte(w)) {
			rbf.Add([]byte(w))
		} else {
			count++
		}
	}
	for _, w := range wordlist2 {
		if rbf.Check([]byte(w)) {
			countfalse++
		}
	}
	fmt.Println("doubletts:", count, "\nfalse positive", countfalse, "\nBloomfilter Has&Add: ", time.Since(st).Nanoseconds()/int64(n), "ns/op  N=", n, "\n")
	rbf = bloomfilter.NewBloomFilter(3, 1)

	fmt.Println("<mySmdb>")
	bbf := New(50*n, 3)
	count = 0
	countfalse = 0
	st = time.Now()
	for _, w := range wordlist {
		if !bbf.Has([]byte(w)) {
			bbf.Add([]byte(w))
		} else {
			count++
		}
	}
	for _, w := range wordlist2 {
		if bbf.Has([]byte(w)) {
			countfalse++
		}
	}
	fmt.Println("doubletts:", count, "\nfalse positive", countfalse, "\nBloomfilter Has&Add: ", time.Since(st).Nanoseconds()/int64(n), "ns/op  N=", n, "\n")
	bbf = Empty()

	fmt.Println("\n\nAbri")
	bf := New(50*n, 3)
	count = 0
	st = time.Now()
	for _, w := range wordlist {
		if !bf.Has([]byte(w)) {
			bf.Add([]byte(w))
		} else {
			count++
		}
	}
	fmt.Println("doubletts:", count, "\nBloomfilter Has&Add: ", time.Since(st).Nanoseconds()/int64(n), "ns/op  N=", n)

	st = time.Now()
	jm := bf.JSONMarshal()
	fmt.Println("Bloomfilter JSONMarshal: ", time.Since(st).Seconds(), "s/op ")
	// fmt.Println(jm)
	// os.Stdout.Write(jm)

	st = time.Now()
	bf = Empty()
	bf.JSONUnmarshal(jm)
	fmt.Println("Bloomfilter JSONUnmarshal: ", time.Since(st).Seconds(), "s/op ")

	st = time.Now()
	for _, w := range wordlist {
		if !bf.Has([]byte(w)) {
			fmt.Println(string(w))
		}
	}
	fmt.Println("Bloomfilter Has: ", time.Since(st).Nanoseconds()/int64(n), "ns/op  N=", n)

	count = 0
	st = time.Now()
	for _, w := range wordlist2 {
		if bf.Has([]byte(w)) {
			count++
		}
	}
	fmt.Println("wrong positive", count, "\nBloomfilter Has: ", time.Since(st).Nanoseconds()/int64(n), "ns/op  N=", n)

}

func timeit(nm string, wl *[]string, fn func([]byte) uint64) {
	wordlist := *wl
	n := len(*wl)
	bf := New(n, 1)
	count := 0
	st := time.Now()
	for _, w := range wordlist {
		hash64 := fn([]byte(w))
		hsh := []byte{uint8(hash64 >> 32), uint8(hash64 >> 16 & 0xff), uint8(hash64 >> 8 & 0xff), uint8(hash64 & 0xff)}
		if !bf.Has(hsh) {
			bf.Add(hsh)
		} else {
			count++
		}
	}
	fmt.Println("fn <", nm, ">\ndoubletts:", count, "\nBloomfilter Has&Add: ", time.Since(st).Nanoseconds()/int64(n), "ns/op  N=", n, "\n")
	// bf = Empty()
}

func getSize(ui64 uint64) (size uint64, exponent uint64) {
	size = 1
	for size < ui64 {
		size <<= 1
		exponent++
	}
	return size, exponent
}

type bloom struct {
	boolSet []bool
	// noEntries  uint64
	sizeExp uint64
	setLocs uint64
	// clippings  []uint64
	// h32        hash.Hash32
	h64 hash.Hash64
}

// func (bf bloom) djb2(f []byte) uint64 {
// 	hash := 5381
// 	for _, c := range f {
// 		hash = hash*33 ^ int(c)
// 	}
// 	return uint64(hash)
// }

func (bf bloom) smdb(f []byte) uint64 {
	hash := uint32(0)
	for _, c := range f {
		hash = uint32(c) + (hash << 6) + (hash << 16) - hash
	}
	return uint64(hash)
}

type bloomJSONImExport struct {
	FilterSet []byte
	SetLocs   uint64
}

func New(entries, locs int) (bloomfilter bloom) {
	size, exponent := getSize(uint64(entries))

	bloomfilter = (bloom{
		boolSet: make([]bool, size),
		// noEntries:  uint64(size),
		sizeExp: exponent,
		setLocs: uint64(locs),
		// clippings:  clips,
		// h32:        fnv.New32a(),
		// h64:        crc64.New(crc64.MakeTable(crc64.ECMA)),
		h64: fnv.New64a(),
	})

	return bloomfilter
}

func Empty() (bloomfilter bloom) {
	return New(1, 1)
}

// func (bf bloom) Indices(entry []byte) (indices []uint64) {
// 	indices = make([]uint64, bf.hashRounds)
// 	bs_Len := big.NewInt(int64(len(bf.boolSet)))
// 	bs_len := bs_Len.Uint64()
// 	base := uint64(0)
// 	bf.hashr.Reset()
// 	bf.hashr.Write(entry)
// 	hashval := uint64(bf.hashr.Sum32())
// 	for i := uint8(0); i < bf.hashRounds; i++ {
// 		base = (base + hashval) % bs_len
// 		indices[i] = base
// 		baseInt := big.NewInt(int64(base))
// 		next := big.NewInt(0).Exp(big.NewInt(2), baseInt, prime32_2)
// 		hashval = next.Mod(next, bs_Len).Uint64()
// 	}
// 	return indices
// }

// func (bf bloom) Indices(entry []byte) (indices []uint64) {
// 	indices = make([]uint64, bf.hashRounds*3)
// 	bf_len := uint64(len(bf.boolSet))
// 	for i := uint8(0); i < bf.hashRounds; i++ {
// 		bf.h32.Reset()
// 		bf.h32.Write(entry)
// 		hash32 := uint64(bf.h32.Sum32())
// 		indices[4*i] = hash32 % bf_len
// 		bf.h64.Reset()
// 		bf.h64.Write(entry)
// 		hash64 := bf.h64.Sum64()
// 		indices[4*i+1] = (hash64&((1<<32)-1) + hash32) % bf_len
// 		indices[4*i+2] = (hash64>>32 + hash32) % bf_len
// 		indices[4*i+3] = hash64 % bf_len
// 		if bf.hashRounds > 1 {
// 			entry = []byte{uint8(hash32 >> 16), uint8(hash32 >> 8 & 0xff), uint8(hash32 & 0xff), uint8(hash64 >> 32), uint8(hash64 >> 16 & 0xff), uint8(hash64 >> 8 & 0xff), uint8(hash64 & 0xff)}
// 		}
// 	}
// 	return indices
// }

// func (bf bloom) Indices(entry []byte) (uint64, uint64) {
// 	// indices = make([]uint64, bf.hashRounds*3)
// 	// indices = make([]uint64, 3)
// 	//for i := uint8(0); i < bf.hashRounds; i++ {
// 	// bf.h32.Reset()
// 	// bf.h32.Write(entry)
// 	// hash32 := uint64(bf.h32.Sum32())
// 	// indices[4*i] = hash32 % bf_len
// 	// bf.h64.Reset()
// 	// bf.h64.Write(entry)
// 	bf.h64.Reset()
// 	bf.h64.Write(entry)
// 	hash64 := bf.h64.Sum64()
// 	// hash64 := bf.smdb(entry)
// 	l := hash64 << 32 >> 32
// 	h := hash64 >> 32

// 	// indices[0] = hash64_l % bf.noEntries
// 	// hash64_l += hash64 >> 32
// 	// indices[1] = hash64_l % bf.noEntries
// 	// indices[2] = (hash64_l ^ hash64>>32) % bf.noEntries
// 	// if bf.hashRounds > 1 {
// 	// 	entry = []byte{uint8(hash64 >> 32), uint8(hash64 >> 16 & 0xff), uint8(hash64 >> 8 & 0xff), uint8(hash64 & 0xff)}
// 	// }
// 	//}
// 	return l, h
// }

func (bf *bloom) Add(entry []byte) {
	// bf.h64.Reset()
	// bf.h64.Write(entry)
	// hash64 := bf.h64.Sum64()
	hash64 := bf.smdb(entry)
	// hash64 := bf.djb2(entry)
	l := hash64 << (bf.sizeExp >> 1) >> (bf.sizeExp >> 1)
	h := hash64 >> (bf.sizeExp >> 1)

	// xx32 := xxhash.New(5381)
	// xx32.Write(entry)
	// l := uint64(xx32.Sum32())
	// h := swap64(l)

	for i := uint64(0); i < bf.setLocs; i++ {
		bf.boolSet[(h+i*l)&(1<<bf.sizeExp-1)] = true
	}
}

func (bf bloom) Has(entry []byte) bool {
	// bf.h64.Reset()
	// bf.h64.Write(entry)
	// hash64 := bf.h64.Sum64()
	hash64 := bf.smdb(entry)
	// // hash64 := bf.djb2(entry)
	l := hash64 << (bf.sizeExp >> 1) >> (bf.sizeExp >> 1)
	h := hash64 >> (bf.sizeExp >> 1)

	// xx32 := xxhash.New(5381)
	// xx32.Write(entry)
	// l := uint64(xx32.Sum32())
	// h := swap64(l)

	for i := uint64(0); i < bf.setLocs; i++ {
		switch bf.boolSet[(h+i*l)&(1<<bf.sizeExp-1)] {
		case false:
			return false
		}
	}
	return true
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
	// bs_len := uint64(len(bf.boolSet))
	// clips := make([]uint64, hashrounds)
	// for i := range clips {
	// 	clips[i] = uint64(float64(bs_len) / (1 + float64(i)*3.1415))
	// }
	_, exponent := getSize(uint64(len(bf.boolSet)))
	(*bf).sizeExp = exponent
	//(*bf).noEntries = size
}

func (bf bloom) ToBytes() []byte {
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
