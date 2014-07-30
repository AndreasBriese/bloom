## bloom filter for go/golang
===

package implements a fast bloom filter with a []bool as 'bit'-set and JSONMarshal/JSONUnmarshal to store/reload the filter after casting it to []uint8. 

===

This bloom filter was developed to strengthen a website-log database and was tested and optimized for this log-entry mask: "2014/%02i/%02i %02i:%02i:%02i /info.html". 
Nonetheless bloom should work with any other form of entries. 

provides 32bit and 64bit boolsets for smaller or larger entry volumes. 
32bit - bloom' s hash function is Berkeley DB smdb hash (slightly modified to optimize for smaller bitsets len<=4096). smdb <--- http://www.cse.yorku.ca/~oz/hash.html
64bit - bloom' s hash function is go's generic FNV64a hash 

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

to work with the bloom filter
