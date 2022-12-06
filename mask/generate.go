package mask

import "strings"

func NewGenerator(characterSet []string, repeat int, greedy bool) *GENERATOR {
	length := len(characterSet)
	maxCount := length

	for i := 1; i < repeat; i++ {
		maxCount = maxCount * length
	}

	gen := &GENERATOR{
		characterSet: characterSet,
		maxRepeat:    repeat,
		MaxCount:     maxCount,
		greedy:       greedy,
	}

	gen.Product()
	return gen
}

func NewGeneratorSingle(s string) *GENERATOR {
	gen := &GENERATOR{
		characterSet: []string{s},
		MaxCount:     1,
		Strings:      []string{s},
	}

	return gen
}

type GENERATOR struct {
	close        bool
	greedy       bool
	characterSet []string
	Strings      []string
	Count        int
	MaxCount     int
	maxRepeat    int
}

func Product(a, b []string) []string {
	ss := make([]string, len(a)*len(b))
	sum := 0
	for _, i := range a {
		for _, j := range b {
			ss[sum] = i + j
			sum++
		}
	}
	return ss
}

func (g *GENERATOR) repeat(ss []string, cur int) []string {
	if g.maxRepeat == 1 {
		return g.characterSet
	}

	if cur < g.maxRepeat {
		if g.greedy {
			g.Strings = append(g.Strings, ss...)
		}
		return g.repeat(Product(ss, g.characterSet), cur+1)
	} else {
		g.Strings = append(g.Strings, ss...)
	}
	return ss
}

func (g *GENERATOR) Inspect() string {
	return strings.Join(g.Strings, "\n")
}

func (g *GENERATOR) Type() ObjectType { return GENERATOR_OBJ }

func (g *GENERATOR) Product() []string {
	return g.repeat(g.characterSet, 1)
}

func (g *GENERATOR) Cross(other *GENERATOR) {
	g.Strings = Product(g.Strings, other.Strings)
}

func (g *GENERATOR) CrossString(ss []string) []string {
	return Product(g.Strings, ss)
}

func (g *GENERATOR) Stream() chan string {
	ch := make(chan string)
	go func() {
		for _, s := range g.Strings {
			g.Count++
			ch <- s
		}
		close(ch)
	}()
	return ch
}
