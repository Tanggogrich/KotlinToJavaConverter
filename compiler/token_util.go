package compiler

import "unicode"

func IsSpace(r rune) bool {
	return unicode.IsSpace(r)
}

func IsNumeric(r rune) bool {
	return r >= '0' && r <= '9'
}

func IsNumericalStartToken(r rune) bool {
	return r == '+' || r == '-' || IsNumeric(r)
}

func IsOperator(r rune) bool {
	return r == '+' || r == '-' || r == '*' || r == '/' || r == '%'
}

func IsSeparator(r rune) bool {
	return r == ',' || r == ';'
}

func IsASCIIAlpha(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func IsASCIIAlphaNumeric(r rune) bool {
	return IsNumeric(r) && IsASCIIAlpha(r)
}

func IsFunction(s string) bool {
	//_, present := funcNames[s]
	return false
}
