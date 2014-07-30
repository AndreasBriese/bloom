## bloom filter for go/golang
===

package implements a fast bloom filter with a []bool as 'bit'-set and JSONMarshal/JSONUnmarshal to store/reload the filter after casting it to []uint8. 

===

This bloom filter was developed to strengthen a website-log database and was tested and optimized for this log-entry mask: "2014/%02i/%02i %02i:%02i:%02i /info.html". 
Nonetheless bloom should work with any other form of entries. 

provides 32bit and 64bit boolsets for smaller or larger entry volumes. 

32bit - bloom: hash function is Berkeley DB smdb hash (slightly modified to optimize for smaller bitsets len<=4096). smdb <--- http://www.cse.yorku.ca/~oz/hash.html

64bit - bloom: hash function is go's generic FNV64a hash 

###install

```sh
go get github.com/AndreasBriese/bloom
```

###test
+ change to folder ../bloom 
+ create wordlist in file "words.txt" (you might use `python permut.py`)
+ run go test within the folder

```go
go test
```

If you've installed the GOCONVEY TDD-framework http://goconvey.co/ you can run the tests automatically.

### usage

after installation add

```go
import (
	...
	"github.com/AndreasBriese/bloom"
	...
	)
```

at your header. In the program use

```go
// create a bloom filter for 65536 items and 1 % wrong-positive ratio 
bf := bloom.New(float64(1<<16), float64(0.01))

// create a bloom filter with 650000 for 65536 items and 7 locs per hash
bf = bloom.New(float64(650000), float64(7))

// add one item
bf.Add([]byte("butter"))

// check if item is in the filter
isIn := bf.Has([]byte("butter"))
isNotIn := bf.Has([]byte("Butter"))

// convert to JSON ([]byte) 
// (the bitset (type []bool) is casted to []uint8 for compression)
Json := bf.JSONMarshal()

// restore a bloom filter from storage 
bf = JSONUnmarshal(Json)
```

to work with the bloom filter.

### why 'fast'? 

It's about 3 times faster than William Fitzgeralds real bitset bloom filter https://github.com/willf/bloom . 
This comes at the cost of 8times greater memory usage for the 'bit'-set which is basically a []bool slice (even if it could be stored in smaller []uint8 JSON-Object). 
	
	Bloom filter (filter size 524288, 7 hashlocs)
	github.com/AndreasBriese/bloom 'Add' 65536 items (10 repetitions): 6304684 ns (96 ns/op)
	github.com/AndreasBriese/bloom 'Has' 65536 items (10 repetitions): 6568663 ns (100 ns/op)
	github.com/willf/bloom 'Add' 65536 items (10 repetitions): 24367224 ns (371 ns/op)
	github.com/willf/bloom 'Test' 65536 items (10 repetitions): 21881142 ns (333 ns/op)
	github.com/dataence/bloom/standard 'Add' 65536 items (10 repetitions): 23041644 ns (351 ns/op)
	github.com/dataence/bloom/standard 'Check' 65536 items (10 repetitions): 19153133 ns (292 ns/op)
	github.com/cabello/bloom 'Add' 65536 items (10 repetitions): 131921507 ns (2012 ns/op)
	github.com/cabello/bloom 'Contains' 65536 items (10 repetitions): 131108962 ns (2000 ns/op)

(on MBPro15 OSX10.8.5 i7 2.4Ghz: bf.Add 91ns/op, bf.Has 102ns/op)


With 32bit bloom filters (bloom32) using smdb, bloom32 does hashing with only 2 bit shifts, one xor and one substraction per byte. smdb is about as fast as fnv64a but gives less collisions with the dataset (see mask above). bloom.New(float64(10 * 1<<16),float64(7)) populated with 1<<16 random items from the dataset (see above) and tested against the rest results in less than 0.05% collisions.   
