package words

import (
	"github.com/chainreactors/words/mask/parser"
	"math/rand"
	"time"
	"unsafe"
)

type Words struct {
	raw string
	ch  chan string
}

var src = rand.NewSource(time.Now().UnixNano())

const (
	// 6 bits to represent a letter index
	letterIdBits = 6
	// All 1-bits as many as letterIdBits
	letterIdMask = 1<<letterIdBits - 1
	letterIdMax  = 63 / letterIdBits
)

func RandPath() string {
	n := 16
	b := make([]byte, n)
	b[0] = byte(0x2f)
	// A rand.Int63() generates 63 random bits, enough for letterIdMax letters!
	for i, cache, remain := n-1, src.Int63(), letterIdMax; i >= 1; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdMax
		}
		if idx := int(cache & letterIdMask); idx < len(parser.Letter) {
			b[i] = parser.Letter[idx]
			i--
		}
		cache >>= letterIdBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}
