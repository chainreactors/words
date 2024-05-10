package words

import (
	"fmt"
	"testing"
)

var rules = ":\nc\nu\nC\n##append numbers 1...5\n$1\n$2\n$3\n$4\n$5\n$6\n$7\n$8\n$9\n$0\n$1 $2 $3\n$1 $2 $3 $4\n$1 $2 $3 $4 $5\n$1 $2 $3 $4 $5 $6\n##current year 2018-2022\n$2 $0 $1 $8\n$2 $0 $1 $9\n$2 $0 $2 $0\n$2 $0 $2 $1\n$2 $0 $2 $2\n#years or month\n$0 $1\n$0 $2\n$0 $3\n$0 $4\n$0 $5\n$0 $6\n$0 $7\n$0 $8\n$0 $9\n$1 $0\n$1 $1\n$1 $2\n$1 $3\n$1 $4\n$1 $5\n$1 $6\n$1 $7\n$1 $8\n$1 $9\n$2 $0\n$2 $1\n$2 $2\n$2 $3\n$2 $4\n$2 $5\n$2 $6\n$2 $7\n$2 $8\n$2 $9\n$3 $0\n$3 $1\n##append sepcial chars\n$!\n$@\n$#\n$$\n$! $@\n$! $@ $#\n$! $@ $# $$\n##special chars + numbers\n$1 $2 $3 $!\n$! $1 $2 $3\n$1 $@ !#\n$! $@ 1#\n##special chars + years\n$2 $0 $1 $8 $!\n$2 $0 $1 $9 $!\n$2 $0 $2 $0 $!\n$2 $0 $2 $1 $!\n$2 $0 $2 $2 $!\n$! $2 $0 $1 $8\n$! $2 $0 $1 $9\n$! $2 $0 $2 $0\n$! $2 $0 $2 $1\n$! $2 $0 $2 $2\n$2 $0 $1 $8 $! $@ $#\n$2 $0 $1 $9 $! $@ $#\n$2 $0 $2 $0 $! $@ $#\n$2 $0 $2 $1 $! $@ $#\n$2 $0 $2 $2 $! $@ $#\n$0 $1 $! \n$0 $2 $!\n$0 $3 $!\n$0 $4 $!\n$0 $5 $!\n$0 $6 $!\n$0 $7 $!\n$0 $8 $!\n$0 $9 $!\n$1 $0 $!\n$1 $1 $!\n$1 $2 $!\n$1 $3 $!\n$1 $4 $!\n$1 $5 $!\n$1 $6 $!\n$1 $7 $!\n$1 $8 $!\n$1 $9 $!\n$2 $0 $!\n$2 $1 $!\n$2 $2 $!\n$2 $3 $!\n$2 $4 $!\n$2 $5 $!\n$2 $6 $!\n$2 $7 $!\n$2 $8 $!\n$2 $9 $!\n$3 $0 $!\n$3 $1 $!\n#all above cap\nc $1\nc $2\nc $3\nc $4\nc $5\nc $6\nc $7\nc $8\nc $9\nc $0\nc $1 $2 $3\nc $1 $2 $3 $4\nc $1 $2 $3 $4 $5\nc $1 $2 $3 $4 $5 $6\nc $2 $0 $1 $8\nc $2 $0 $1 $9\nc $2 $0 $2 $0\nc $2 $0 $2 $1\nc $2 $0 $2 $2\nc $!\nc $@\nc $#\nc $$\nc $! $@\nc $! $@ $#\nc $! $@ $# $$\nc $1 $2 $3 $!\nc $! $1 $2 $3\nc $1 $@ !#\nc $! $@ 1#\nc $2 $0 $1 $8 $!\nc $2 $0 $1 $9 $!\nc $2 $0 $2 $0 $!\nc $2 $0 $2 $1 $!\nc $2 $0 $2 $2 $!\nc $! $2 $0 $1 $8\nc $! $2 $0 $1 $9\nc $! $2 $0 $2 $0\nc $! $2 $0 $2 $1\nc $! $2 $0 $2 $2\nc $2 $0 $1 $8 $! $@ $#\nc $2 $0 $1 $9 $! $@ $#\nc $2 $0 $2 $0 $! $@ $#\nc $2 $0 $2 $1 $! $@ $#\nc $2 $0 $2 $2 $! $@ $#\nc $0 $1 $! \nc $0 $2 $!\nc $0 $3 $!\nc $0 $4 $!\nc $0 $5 $!\nc $0 $6 $!\nc $0 $7 $!\nc $0 $8 $!\nc $0 $9 $!\nc $1 $0 $!\nc $1 $1 $!\nc $1 $2 $!\nc $1 $3 $!\nc $1 $4 $!\nc $1 $5 $!\nc $1 $6 $!\nc $1 $7 $!\nc $1 $8 $!\nc $1 $9 $!\nc $2 $0 $!\nc $2 $1 $!\nc $2 $2 $!\nc $2 $3 $!\nc $2 $4 $!\nc $2 $5 $!\nc $2 $6 $!\nc $2 $7 $!\nc $2 $8 $!\nc $2 $9 $!\nc $3 $0 $!\nc $3 $1 $!\nc $0 $1\nc $0 $2\nc $0 $3\nc $0 $4\nc $0 $5\nc $0 $6\nc $0 $7\nc $0 $8\nc $0 $9\nc $1 $0\nc $1 $1\nc $1 $2\nc $1 $3\nc $1 $4\nc $1 $5\nc $1 $6\nc $1 $7\nc $1 $8\nc $1 $9\nc $2 $0\nc $2 $1\nc $2 $2\nc $2 $3\nc $2 $4\nc $2 $5\nc $2 $6\nc $2 $7\nc $2 $8\nc $2 $9\nc $3 $0\nc $3 $1"

func TestNewWorder(t *testing.T) {
	worder, _ := NewWorderWithDsl("tz-{?u#2}", nil, nil)
	//worder.SetRules(rules, "<5")
	worder.Fns = append(worder.Fns, func(w string) []string {
		return []string{w + "a", w + "b", w + "c"}
	})
	worder.Fns = append(worder.Fns, func(w string) []string {
		return []string{w + "a"}
	})

	worder.Run()
	for w := range worder.C {
		fmt.Println(w)
	}
}
