package words

import (
	"bufio"
	"github.com/chainreactors/words/mask"
	"os"
	"strings"
)

func NewWorder(wordlist []string) *Worder {
	worder := &Worder{
		token:    0,
		wordlist: wordlist,
		ch:       make(chan string),
		C:        make(chan string),
		//checkCh:     make(chan string),
	}

	go func() {
		for _, w := range worder.wordlist {
			worder.ch <- strings.TrimSpace(w)
		}
		worder.Close()
	}()

	worder.init()
	return worder
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

	worder.init()
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
	worder.init()
	return worder
}

type Worder struct {
	ch chan string
	C  chan string
	//checkCh     chan string
	token    int
	wordlist []string
	scanner  *bufio.Scanner
	Closed   bool
}

func (word *Worder) init() {
	go func() {
		for w := range word.ch {
			word.token++
			word.C <- w
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
