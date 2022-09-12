package words

import (
	"bufio"
	"os"
	"strings"
)

var DefaultPeriod = 0

func NewWorder(wordlist []string) *Worder {
	worder := &Worder{
		checkPeriod: DefaultPeriod,
		token:       0,
		wordlist:    wordlist,
		ch:          make(chan string),
		C:           make(chan string),
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
		checkPeriod: DefaultPeriod,
		token:       0,
		scanner:     bufio.NewScanner(file),
		ch:          make(chan string),
		C:           make(chan string),
		//checkCh:     make(chan string, 10),
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

type Worder struct {
	ch chan string
	C  chan string
	//checkCh     chan string
	token       int
	wordlist    []string
	scanner     *bufio.Scanner
	checkPeriod int
	Closed      bool
}

func (word *Worder) SetPeriod(period int) {
	word.checkPeriod = period
}

func (word *Worder) init() {
	go func() {
		for path := range word.ch {
			word.token++
			if word.checkPeriod != 0 && (word.token%word.checkPeriod) == 0 {
				word.C <- RandPath()
			}
			word.C <- path
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
