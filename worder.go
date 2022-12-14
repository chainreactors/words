package words

import (
	"bufio"
	"github.com/chainreactors/words/mask"
	"github.com/chainreactors/words/rule"
	"os"
	"strings"
)

var CustomWords [][]string

func NewWorder(word string, dicts [][]string, fns []func(string) string) (*Worder, error) {
	worder := &Worder{
		token: 0,
		ch:    make(chan string),
		C:     make(chan string),
		Fns:   fns,
	}

	var err error
	var input chan string
	if dicts != nil {
		input, err = mask.RunToStream(word, dicts)
	} else {
		input, err = mask.RunToStream(word, CustomWords)
	}

	if err != nil {
		return nil, err
	}

	go func() {
		for w := range input {
			worder.ch <- strings.TrimSpace(w)
		}
		worder.Close()
	}()
	return worder, nil
}

func NewWorderWithFns(wordlist []string, fns []func(string) string) *Worder {
	worder := &Worder{
		token: 0,
		ch:    make(chan string),
		C:     make(chan string),
		Fns:   fns,
	}

	go func() {
		for _, w := range wordlist {
			worder.ch <- strings.TrimSpace(w)
		}
		worder.Close()
	}()

	return worder
}

func NewWorderWithFile(file *os.File, fns []func(string) string) *Worder {
	worder := &Worder{
		token:   0,
		scanner: bufio.NewScanner(file),
		ch:      make(chan string),
		C:       make(chan string),
		Fns:     fns,
	}
	go func() {
		for worder.scanner.Scan() {
			worder.ch <- strings.TrimSpace(worder.scanner.Text())
		}
		worder.Close()
	}()

	return worder
}

func NewWorderWithDSL(dsl string, params [][]string) (*Worder, error) {
	worder := &Worder{
		token: 0,
		C:     make(chan string),
	}

	var ch chan string
	var err error
	if params != nil {
		ch, err = mask.RunToStream(dsl, params)
	} else {
		ch, err = mask.RunToStream(dsl, CustomWords)
	}

	if err != nil {
		return nil, err
	}
	worder.ch = ch
	return worder, nil
}

type Worder struct {
	ch      chan string
	C       chan string
	token   int
	Rules   []rule.Expression
	scanner *bufio.Scanner
	Fns     []func(string) string
	Closed  bool
}

func (word *Worder) CompileRules(rules string, filter string) {
	word.Rules = rule.Compile(rules, filter)
}

func (word *Worder) Run() {
	go func() {
		for w := range word.ch {
			word.token++
			if w == "" {
				continue
			}
			for _, fn := range word.Fns {
				w = fn(w)
			}
			word.C <- w
		}
		close(word.C)
	}()
}

func (word *Worder) RunWithRules() {
	go func() {
		for w := range word.ch {
			word.token++
			if w == "" {
				continue
			}

			if word.Rules != nil {
				for r := range rule.RunAsStream(word.Rules, w) {
					for _, fn := range word.Fns {
						r = fn(r)
					}
					word.C <- r
				}
			} else {
				for _, fn := range word.Fns {
					w = fn(w)
				}
				word.C <- w
			}
		}
		close(word.C)
	}()
}
func (word *Worder) All() []string {
	var ws []string
	for w := range word.C {
		ws = append(ws, w)
	}
	return ws
}

func (word *Worder) Close() {
	word.Closed = true
	close(word.ch)
}
