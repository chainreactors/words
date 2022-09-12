package words

import (
	"math/rand"
	"time"
	"unsafe"
)

var (
	Lowercase    = "abcdefghijklmnopqrstuvwxyz"
	Uppercase    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Letter       = Lowercase + Uppercase
	Digit        = "0123456789"
	LowercaseHex = Digit + "abcdef"
	UppercaseHex = Digit + "ABCDEF"
	Hex          = Digit + "abcdefABCDEF"
	Punctuation  = "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
	Printable    = Letter + Digit + Punctuation
	Whitespace   = "\t\n\r\x0b\x0c"
	Custom       = ""
)

var MetawordMap = map[string]string{
	"a": Lowercase,
	"A": Uppercase,
	"w": Letter,
	"d": Digit,
	"h": LowercaseHex,
	"H": UppercaseHex,
	"x": Hex,

	"p": Punctuation,
	"P": Printable,
	"s": Whitespace,
	"c": Custom,
}

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
		if idx := int(cache & letterIdMask); idx < len(Letter) {
			b[i] = Letter[idx]
			i--
		}
		cache >>= letterIdBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}
