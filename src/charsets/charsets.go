package charsets

import "strings"

type RandInt interface {
	Int() int
}

type Charset interface {
	TestRune(char rune) bool
	RandomRune(r RandInt) rune
	RandomString(r RandInt, size int) string

	String() string
}

func NewCharsetFromASCII(offset, size int) Charset {
	return charsetASCII{offset: offset, size: size}
}

type charsetASCII struct {
	offset int
	size   int
}

func (c charsetASCII) TestRune(char rune) bool {
	return int(char) >= c.offset && int(char) < c.offset+c.size
}

func (c charsetASCII) RandomRune(r RandInt) rune {
	num := c.offset + r.Int()%(c.size-1)
	return rune(num)
}

func (c charsetASCII) RandomString(r RandInt, size int) string {
	builder := strings.Builder{}
	for i := 0; i < size; i++ {
		builder.WriteRune(c.RandomRune(r))
	}
	return builder.String()
}

func (c charsetASCII) String() string {
	builder := strings.Builder{}
	for i := 0; i < c.size; i++ {
		builder.WriteRune(rune(c.offset + i))
	}
	return builder.String()
}

func NewCharsetFromString(s string) Charset {
	charsArray := make([]rune, len(s))
	charsMap := make(map[rune]bool, len(s))
	for i, v := range s {
		charsArray[i] = v
		charsMap[v] = true
	}

	return charsetFromString{
		charsArray: charsArray,
		charsMap:   charsMap,
	}
}

type charsetFromString struct {
	charsMap   map[rune]bool
	charsArray []rune
}

func (c charsetFromString) TestRune(char rune) bool {
	return c.charsMap[char]
}

func (c charsetFromString) RandomRune(r RandInt) rune {
	num := r.Int() % (len(c.charsArray) - 1)
	return c.charsArray[num]
}

func (c charsetFromString) RandomString(r RandInt, size int) string {
	builder := strings.Builder{}
	for i := 0; i < size; i++ {
		builder.WriteRune(c.RandomRune(r))
	}
	return builder.String()
}

func (c charsetFromString) String() string {
	builder := strings.Builder{}
	for _, v := range c.charsArray {
		builder.WriteRune(v)
	}
	return builder.String()
}

func NewCharsetUnion(opts ...Charset) Charset {
	charsets := []Charset{}
	return charsetUnion{
		charsets: append(charsets, opts...),
	}
}

type charsetUnion struct {
	charsets []Charset
}

func (c charsetUnion) TestRune(char rune) bool {
	for _, charset := range c.charsets {
		if charset.TestRune(char) {
			return true
		}
	}
	return false
}

func (c charsetUnion) RandomRune(r RandInt) rune {
	index := r.Int() % (len(c.charsets) - 1)
	charset := c.charsets[index]

	return charset.RandomRune(r)
}

func (c charsetUnion) RandomString(r RandInt, size int) string {
	builder := strings.Builder{}
	for i := 0; i < size; i++ {
		index := r.Int() % (len(c.charsets) - 1)
		charset := c.charsets[index]
		builder.WriteRune(charset.RandomRune(r))
	}
	return builder.String()
}

func (c charsetUnion) String() string {
	return ""
}
