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

type charsetPart struct {
	Offset int
	Size   int
}

var charsets = map[Charset][]charsetPart{}

func NewRand() *RandomUtil {
	charsetLettersLower := charsetPart{ //CharsetLettersLower
		Offset: 0x41,
		Size:   26,
	}

	charsetLettersUpper := charsetPart{ //CharsetLettersUpper
		Offset: 0x61,
		Size:   26,
	}

	charsetNumeric := charsetPart{ //CharsetLettersNumeric
		Offset: 0x30,
		Size:   10,
	}

	charsets = map[Charset][]charsetPart{
		CharsetNumeric:      {charsetNumeric},
		CharsetLettersLower: {charsetLettersLower},
		CharsetLettersUpper: {charsetLettersUpper},
		CharsetLetters:      {charsetLettersLower, charsetLettersUpper},
		CharsetAll:          {charsetLettersLower, charsetLettersUpper, charsetNumeric},
	}

	return &RandomUtil{}
}

type RandomUtil struct{}

func (r *RandomUtil) RandomID(outputLenght int, charset Charset) string {
	src := rand.NewSource(time.Now().UnixMicro())
	randGen := rand.New(src)

	charsetData := charsets[charset]

	builder := strings.Builder{}
	for i := 0; i < outputLenght; i++ {
		charsetIdx := randGen.Int() % len(charsetData)
		charsetPart := charsetData[charsetIdx]

		byte := charsetPart.Offset + (randGen.Int() % charsetPart.Size)
		builder.WriteRune(rune(byte))
	}
	return builder.String()
}
