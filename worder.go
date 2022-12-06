package words

import (
	"bufio"
	"github.com/chainreactors/words/mask"
	"github.com/chainreactors/words/rule"
	"os"
	"strings"
)

func NewWorder(wordlist []string, fns []func(string) string) *Worder {
	worder := &Worder{
		token:    0,
		wordlist: wordlist,
		ch:       make(chan string),
		C:        make(chan string),
		Fns:      fns,
	}

	go func() {
		for _, w := range worder.wordlist {
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

func NewWorderWithDSL(dsl string) *Worder {
	worder := &Worder{
		token: 0,
		C:     make(chan string),
	}
	ch, err := mask.RunToStream(dsl)
	if err != nil {
		panic(err)
	}
	worder.ch = ch
	return worder
}

type Worder struct {
	ch       chan string
	C        chan string
	token    int
	Rules    []rule.Expression
	wordlist []string
	scanner  *bufio.Scanner
	Fns      []func(string) string
	Closed   bool
}

func (word *Worder) CompileRules(rules string) {
	word.Rules = rule.Compile(rules)
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
			for _, fn := range word.Fns {
				w = fn(w)
			}
			for r := range rule.RunAsStream(word.Rules, w) {
				word.C <- r
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
