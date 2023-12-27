package mask

import "strings"

func NewGenerator(characterSet []string, repeat int, greedy bool) *GENERATOR {
	length := len(characterSet)
	var maxCount int
	tmp := make([]int, repeat+1)
	tmp[0] = 1
	for i := 1; i <= repeat; i++ {
		tmp[i] = tmp[i-1] * length
		if greedy {
			maxCount += tmp[i]
		} else {
			maxCount = tmp[i]
		}
	}

	gen := &GENERATOR{
		characterSet: characterSet,
		maxRepeat:    repeat,
		MaxCount:     maxCount,
		Streamer:     make(chan string),
		greedy:       greedy,
	}

	if gen.greedy {
		gen.Streamer = gen.GreedyProduct()
	} else {
		gen.Streamer = gen.Product()
	}
	return gen
}

func NewGeneratorSingle(s string) *GENERATOR {
	gen := &GENERATOR{
		characterSet: []string{s},
		MaxCount:     1,
	}
	gen.Streamer = gen.Product()
	return gen
}

type GENERATOR struct {
	close        bool
	greedy       bool
	characterSet []string
	Streamer     chan string
	Count        int
	MaxCount     int
	maxRepeat    int
}

func Product(a chan string, b []string) chan string {
	ch := make(chan string)
	go func() {
		for i := range a {
			b := wrapStream(b)
			for j := range b {
				ch <- i + j
			}
		}
		close(ch)
	}()
	return ch
}

func ProductGenerator(a chan string, b *GENERATOR) chan string {
	ch := make(chan string)
	go func() {
		for i := range a {
			for j := range b.Streamer {
				ch <- i + j
			}
			b.Reset()
		}
		close(ch)
	}()
	return ch
}

func wrapStream(s []string) chan string {
	ch := make(chan string)
	go func() {
		for _, i := range s {
			ch <- i
		}
		close(ch)
	}()
	return ch
}

func (g *GENERATOR) repeat(ss chan string, cur, max int) chan string {
	if max == 1 {
		return ss
	}

	if cur < max {
		return g.repeat(Product(ss, g.characterSet), cur+1, max)
	}
	return ss
}

func (g *GENERATOR) Inspect() string {
	return strings.Join(g.All(), "\n")
}

func (g *GENERATOR) Type() ObjectType { return GENERATOR_OBJ }

func (g *GENERATOR) Product() chan string {
	return g.repeat(wrapStream(g.characterSet), 1, g.maxRepeat)
}

func (g *GENERATOR) GreedyProduct() chan string {
	ch := make(chan string)
	go func() {
		for i := 1; i <= g.maxRepeat; i++ {
			for s := range g.repeat(wrapStream(g.characterSet), 1, i) {
				ch <- s
			}
		}
		close(ch)
	}()
	return ch
}

func (g *GENERATOR) Cross(other *GENERATOR) {
	g.MaxCount = g.MaxCount * other.MaxCount
	g.Streamer = ProductGenerator(g.Streamer, other)
}

func (g *GENERATOR) CrossString(ss []string) {
	g.MaxCount = g.MaxCount * len(ss)
	g.Streamer = Product(g.Streamer, ss)
}

func (g *GENERATOR) Reset() {
	if g.greedy {
		g.Streamer = g.GreedyProduct()
	} else {
		g.Streamer = g.Product()
	}
}

func (g *GENERATOR) All() []string {
	ss := make([]string, g.MaxCount)
	i := 0
	for s := range g.Streamer {
		ss[i] = s
		i++
	}
	return ss
}
