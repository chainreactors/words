package words

import (
	"bufio"
	"github.com/chainreactors/words/mask"
	"github.com/chainreactors/words/rule"
	"os"
	"strings"
)

type WordFunc func(string) []string

var CustomWords [][]string

func NewWorder(list []string) *Worder {
	worder := &Worder{
		token: 0,
		ch:    make(chan string),
		C:     make(chan string),
	}
	go func() {
		for _, l := range list {
			worder.ch <- l
		}
		close(worder.ch)
	}()
	return worder
}

func NewWorderWithChan(ch chan string) *Worder {
	worder := &Worder{
		token: 0,
		ch:    make(chan string),
		C:     make(chan string),
	}
	go func() {
		for w := range ch {
			worder.ch <- w
		}
		close(worder.ch)
	}()
	return worder
}

func NewWorderWithDsl(word string, params [][]string, keywords map[string][]string) (*Worder, error) {
	worder := &Worder{
		token: 0,
		ch:    make(chan string),
		C:     make(chan string),
	}

	var err error
	var input chan string
	if params != nil {
		input, err = mask.RunToStream(word, params, keywords)
	} else {
		input, err = mask.RunToStream(word, CustomWords, keywords)
	}
	if err != nil {
		return nil, err
	}

	worder.ch = input
	return worder, nil
}

func NewWorderWithFile(file *os.File) *Worder {
	worder := &Worder{
		token:   0,
		scanner: bufio.NewScanner(file),
		ch:      make(chan string),
		C:       make(chan string),
	}
	go func() {
		for worder.scanner.Scan() {
			worder.ch <- strings.TrimSpace(worder.scanner.Text())
		}
		close(worder.ch)
	}()

	return worder
}

type Worder struct {
	ch      chan string
	C       chan string
	token   int
	Rules   []rule.Expression
	scanner *bufio.Scanner
	Fns     []WordFunc
	Closed  bool
}

func (word *Worder) SetRules(rules string, filter string) {
	word.Rules = rule.Compile(rules, filter).Expressions
}

func (word *Worder) AddFunction(f WordFunc) {
	word.Fns = append(word.Fns, f)
}

func (word *Worder) EvalFunctions(w string) []string {
	var ss []string
	for _, f := range word.Fns {
		if ss != nil {
			for i, s := range ss {
				newss := f(s)
				if len(newss) == 1 {
					ss[i] = newss[0]
				} else if len(newss) > 1 {
					ss = append(ss[:i], newss...)
					ss = append(ss, ss[i+1:]...)
				} else {
					ss = append(ss[:i], ss[i+1:]...)
				}
			}
		} else {
			ss = f(w)
		}
	}
	return ss
}

func (word *Worder) Run() {
	go func() {
		for w := range word.ch {
			word.token++
			if w == "" {
				continue
			}
			if word.Rules != nil {
				for r := range rule.RunAsStream(word.Rules, w) {
					if word.Fns != nil {
						if ws := word.EvalFunctions(w); ws != nil {
							for _, i := range word.EvalFunctions(w) {
								word.C <- i
							}
						} else {
							word.C <- "" // 表示skip
						}
					} else {
						word.C <- r
					}
				}
			} else {
				if word.Fns != nil {
					if ws := word.EvalFunctions(w); ws != nil {
						for _, i := range word.EvalFunctions(w) {
							word.C <- i
						}
					} else {
						word.C <- "" // 表示skip
					}
				} else {
					word.C <- w
				}
			}
		}
		close(word.C)
		word.Closed = true
	}()
}

func (word *Worder) All() []string {
	var ws []string
	for w := range word.C {
		ws = append(ws, w)
	}
	return ws
}
