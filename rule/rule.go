package rule

import (
	"strings"
)

func toggleCase(s string) string {
	var tmp strings.Builder
	for _, c := range s {
		if c >= 97 && c <= 122 {
			c = c - 32
		} else if c >= 65 && c <= 90 {
			c = c + 32
		}
		tmp.WriteRune(c)
	}
	return tmp.String()
}

func toggleCasePosition(s string, pos int) string {
	return toggleCase(s[:pos]) + toggleCase(s[pos:pos+1]) + toggleCase(s[pos+1:])
}

// reverse the string
func reverse(s string) string {
	var tmp strings.Builder
	for i := len(s) - 1; i >= 0; i-- {
		tmp.WriteByte(s[i])
	}
	return tmp.String()
}

func ProcessFunction(word string, toks Tokens) string {
	if word == "" {
		return ""
	}
	funcName := toks.MustNext()
	raw := word
	length := len(word)
	switch funcName {
	case ":":
		return word
	case "l":
		return strings.ToLower(word)
	case "u":
		return strings.ToUpper(word)
	case "c":
		if length < 1 {
			return strings.ToUpper(word)
		}
		return strings.ToUpper(string(raw[0])) + strings.ToLower(word)[1:]
	case "C":
		if length < 1 {
			return strings.ToLower(word)
		}
		return strings.ToLower(string(raw[0])) + strings.ToUpper(word)[1:]
	case "t":
		return toggleCase(word)
	case "T":
		i := toks.MustNextInt()
		return toggleCasePosition(word, i)
	case "r":
		return reverse(word)
	case "d":
		return word + word
	case "p":
		t := toks.MustNextInt()
		for i := 0; i < t; i++ {
			word += raw
		}
		return word
	case "f":
		return word + reverse(word)
	case "{":
		if len(word) == 1 {
			return word
		}
		return word[1:] + word[0:1]
	case "}":
		if len(word) == 1 {
			return word
		}
		return word[length-1:] + word[0:length-1]
	case "$":
		s := toks.MustNext()
		return word + s
	case "^":
		s := toks.MustNext()
		return s + word
	case "[":
		return word[1:]
	case "]":
		if len(word) == 1 {
			return ""
		}
		return word[:length-1]
	case "D":
		i := toks.MustNextInt()
		if i > length-1 {
			return word
		} else if i == len(word) {
			return word[:i]
		} else {
			return word[:i] + word[i+1:]
		}
	case "x":
		n := toks.MustNextInt()
		m := toks.MustNextInt()
		if n < length {
			if m < length {
				return word[n:m]
			} else {
				return word[n:]
			}
		} else {
			return ""
		}
	case "O":
		n := toks.MustNextInt()
		m := toks.MustNextInt()
		if n < length {
			if m < length {
				return word[:n] + word[m:]
			} else {
				return word[:n]
			}
		}
		return word
	case "i":
		i := toks.MustNextInt()
		s := toks.MustNext()
		if i < length {
			return word[:i] + s + word[i:]
		} else {
			return word + s
		}
	case "o":
		i := toks.MustNextInt()
		s := toks.MustNext()
		if i < length {
			return word[:i] + s + word[i+1:]
		} else {
			return word
		}
	case "'":
		i := toks.MustNextInt()
		if i < length-1 {
			return word[i:]
		} else {
			return word
		}
	case "s":
		x := toks.MustNext()
		y := toks.MustNext()
		return strings.Replace(word, x, y, -1)
	case "@":
		x := toks.MustNext()
		return strings.Replace(word, x, "", -1)
	case "z":
		i := toks.MustNextInt()
		for j := 0; j < i; j++ {
			word = string(raw[0]) + word
		}
		return word
	case "Z":
		i := toks.MustNextInt()
		for j := 0; j < i; j++ {
			word = string(raw[length-1]) + word
		}
		return word
	case "k":
		if length == 1 {
			return word
		} else if length == 2 {
			return string(word[1]) + string(word[0])
		} else {
			return string(word[1]) + string(word[0]) + word[2:]
		}
	case "K":
		if length == 1 {
			return word
		} else if length == 2 {
			return string(word[1]) + string(word[0])
		} else {
			return word[:length-2] + string(word[length-2]) + string(word[length-1])
		}
	case "*":
		n := toks.MustNextInt()
		m := toks.MustNextInt()
		if n < length && m < length {
			bs := []byte(word)
			tmp := bs[n]
			bs[n] = bs[m]
			bs[m] = tmp
			return string(bs)
		}
		return word
	case "L":
		i := toks.MustNextInt()
		if i < length {
			bs := []byte(word)
			bs[i] = bs[i] << 1
			return string(bs)
		}
		return word
	case "R":
		i := toks.MustNextInt()
		if i < length {
			bs := []byte(word)
			bs[i] = bs[i] >> 1
			return string(bs)
		}
		return word
	case "+":
		i := toks.MustNextInt()
		if i < length {
			bs := []byte(word)
			bs[i] = bs[i] + 1
			return string(bs)
		}
		return word
	case "-":
		i := toks.MustNextInt()
		if i < length {
			bs := []byte(word)
			bs[i] = bs[i] - 1
			return string(bs)
		}
		return word
	case ".":
		i := toks.MustNextInt()
		if i < length {
			bs := []byte(word)
			bs[i] = bs[i+1]
			return string(bs)
		}
		return word
	case ",":
		i := toks.MustNextInt()
		if i < length {
			bs := []byte(word)
			bs[i] = bs[i-1]
			return string(bs)
		}
		return word
	case "y":
		i := toks.MustNextInt()
		if i < length {
			return word[:i] + word
		} else {
			return word + word
		}
	case "Y":
		i := toks.MustNextInt()
		if i < length {
			return word[length-i:] + word
		} else {
			return word + word
		}
	case "E":
		return strings.Title(strings.ToLower(word))
	case "e":
		i := toks.MustNextInt()
		word = strings.ToLower(word)
		if i < length {
			return word[:i] + strings.Title(word)
		} else {
			return word
		}
	default:
		return word
	}
}

// if func return false, will reject the word
func ProcessFilter(word string, toks Tokens) bool {
	if word == "" {
		return false
	}
	funcName := toks.MustNext()
	length := len(word)
	switch funcName {
	case "<":
		i := toks.MustNextInt()
		if length < i {
			return true
		}
		return false
	case ">":
		i := toks.MustNextInt()
		if length > i {
			return true
		}
		return false
	case "_":
		i := toks.MustNextInt()
		if length == i {
			return true
		}
		return false
	case "!":
		x := toks.MustNext()
		if strings.Contains(word, x) {
			return false
		}
		return true
	case "/":
		x := toks.MustNext()
		if strings.Contains(word, x) {
			return true
		}
		return false
	case "(":
		x := toks.MustNext()
		if strings.HasPrefix(word, x) {
			return true
		}
		return false
	case ")":
		x := toks.MustNext()
		if strings.HasSuffix(word, x) {
			return true
		}
		return false
	case "=":
		i := toks.MustNextInt()
		x := toks.MustNext()
		if i < length {
			return string(word[i]) == x
		}
		return false
	case "%":
		i := toks.MustNextInt()
		x := toks.MustNext()
		if strings.Count(word, x) < i {
			return false
		}
		return true
	default:
		return true
	}
}
