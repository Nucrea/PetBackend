package charsets

type CharsetType int

const (
	CharsetTypeAll CharsetType = iota
	CharsetTypeLettersLower
	CharsetTypeLettersUpper
	CharsetTypeLetters
	CharsetTypeNumeric
)

var (
	charsetNumeric      = NewCharsetFromASCII(0x30, 10)
	charsetLettersLower = NewCharsetFromASCII(0x41, 26)
	charsetLettersUpper = NewCharsetFromASCII(0x61, 26)
	charsetLetters      = NewCharsetUnion(charsetLettersLower, charsetLettersUpper)
	charsetAll          = NewCharsetUnion(charsetNumeric, charsetLettersLower, charsetLettersUpper)
)

func GetCharset(charsetType CharsetType) Charset {
	switch charsetType {
	case CharsetTypeNumeric:
		return charsetNumeric
	case CharsetTypeLettersLower:
		return charsetLettersLower
	case CharsetTypeLettersUpper:
		return charsetLettersLower
	case CharsetTypeLetters:
		return charsetLetters
	case CharsetTypeAll:
		return charsetAll
	default:
		return nil
	}
}
