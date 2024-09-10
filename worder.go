package words

import (
	"bufio"
	"errors"
	"github.com/chainreactors/words/mask"
	"github.com/chainreactors/words/rule"
	"os"
	"strings"
)

type WordFunc func(string) []string

var CustomWords [][]string
var (
	ErrNilInputChannel = errors.New("input channel is nil")
)

func NewWorder(word *Worder) (*Worder, error) {
	if word.Input == nil {
		return nil, ErrNilInputChannel
	}
	word.Output = make(chan string)
	return word, nil
}

func NewWorderWithList(list []string) *Worder {
	worder := &Worder{
		Input:  make(chan string),
		Output: make(chan string),
	}
	go func() {
		for _, l := range list {
			worder.Input <- l
		}
		close(worder.Input)
	}()
	return worder
}

func NewWorderWithChan(ch chan string) *Worder {
	worder := &Worder{
		Input:  make(chan string),
		Output: make(chan string),
	}
	go func() {
		for w := range ch {
			worder.Input <- w
		}
		close(worder.Input)
	}()
	return worder
}

func NewWorderWithDsl(word string, params [][]string, keywords map[string][]string) (*Worder, error) {
	worder := &Worder{
		Input:  make(chan string),
		Output: make(chan string),
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

	worder.Input = input
	return worder, nil
}

func NewWorderWithFile(file *os.File) *Worder {
	worder := &Worder{
		Input:  make(chan string),
		Output: make(chan string),
	}

	scanner := bufio.NewScanner(file)

	go func() {
		for scanner.Scan() {
			worder.Input <- strings.TrimSpace(scanner.Text())
		}
		close(worder.Input)
	}()

	return worder
}

type Worder struct {
	Input  chan string
	Output chan string
	Rules  []rule.Expression
	Fns    []WordFunc
	Closed bool
	token  int
	count  int
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
		for w := range word.Input {
			word.count++
			if w == "" {
				continue
			}
			if word.Rules != nil {
				for r := range rule.RunAsStream(word.Rules, w) {
					if word.Fns != nil {
						if ws := word.EvalFunctions(w); ws != nil {
							for _, i := range word.EvalFunctions(w) {
								word.count++
								word.Output <- i
							}
						} else {
							word.Output <- "" // 表示skip
						}
					} else {
						word.count++
						word.Output <- r
					}
				}
			} else {
				if word.Fns != nil {
					if ws := word.EvalFunctions(w); ws != nil {
						for _, i := range word.EvalFunctions(w) {
							word.count++
							word.Output <- i
						}
					} else {
						word.Output <- "" // 表示skip
					}
				} else {
					word.count++
					word.Output <- w
				}
			}
		}
		close(word.Output)
	}()
}

func (word *Worder) All() []string {
	var ws []string
	for w := range word.Output {
		ws = append(ws, w)
	}
	return ws
}

func (word *Worder) Count() int {
	return word.count
}
