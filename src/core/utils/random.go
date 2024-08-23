package utils

import (
	"math/rand"
	"strings"
	"time"
)

type Charset int

const (
	CharsetAll Charset = iota
	CharsetLettersLower
	CharsetLettersUpper
	CharsetLetters
	CharsetNumeric
)

type charsetBlock struct {
	Offset int
	Size   int
}

func NewRand() *RandomUtil {
	charsetLettersLower := charsetBlock{
		Offset: 0x41,
		Size:   26,
	}

	charsetLettersUpper := charsetBlock{
		Offset: 0x61,
		Size:   26,
	}

	charsetNumeric := charsetBlock{
		Offset: 0x30,
		Size:   10,
	}

	return &RandomUtil{
		charsets: map[Charset][]charsetBlock{
			CharsetNumeric:      {charsetNumeric},
			CharsetLettersLower: {charsetLettersLower},
			CharsetLettersUpper: {charsetLettersUpper},
			CharsetLetters:      {charsetLettersLower, charsetLettersUpper},
			CharsetAll:          {charsetLettersLower, charsetLettersUpper, charsetNumeric},
		},
	}
}

type RandomUtil struct {
	charsets map[Charset][]charsetBlock
}

func (r *RandomUtil) RandomID(outputLenght int, charset Charset) string {
	src := rand.NewSource(time.Now().UnixMicro())
	randGen := rand.New(src)

	charsetBlocks := r.charsets[charset]

	builder := strings.Builder{}
	for i := 0; i < outputLenght; i++ {
		charsetBlock := charsetBlocks[randGen.Int()%len(charsetBlocks)]

		byte := charsetBlock.Offset + (randGen.Int() % charsetBlock.Size)
		builder.WriteRune(rune(byte))
	}
	return builder.String()
}
