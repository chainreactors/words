package words

import (
	"bufio"
	"github.com/chainreactors/words/mask"
	"github.com/chainreactors/words/rule"
	"os"
	"strings"
)

var CustomWords [][]string

func NewWorder(word string, params [][]string, keywords map[string][]string) (*Worder, error) {
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
		worder.Close()
	}()

	return worder
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
	word.Rules = rule.Compile(rules, filter).Expressions
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
					if r == "" {
						continue
					}
					word.C <- r
				}
			} else {
				for _, fn := range word.Fns {
					w = fn(w)
				}
				if w == "" {
					continue
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
