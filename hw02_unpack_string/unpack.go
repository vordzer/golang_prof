package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	if len(s) == 0 {
		return "", nil
	}
	const esc rune = '\\'
	var curRune rune
	var res strings.Builder
	var (
		isEsc     = false
		runeReady = false
	)
	for _, r := range s {
		if isEsc {
			if runeReady {
				res.WriteRune(curRune)
			}
			isEsc = false
			curRune = r
			continue
		}
		switch {
		case r == esc:
			isEsc = true
			runeReady = true
		case unicode.IsDigit(r):
			if !runeReady {
				return "", ErrInvalidString
			}
			rep, _ := strconv.Atoi(string(r))
			res.WriteString(strings.Repeat(string(curRune), rep))
			runeReady = false
		default:
			if runeReady {
				res.WriteRune(curRune)
			}
			curRune = r
			runeReady = true
		}
	}
	if runeReady {
		res.WriteRune(curRune)
	}
	if isEsc {
		return "", ErrInvalidString
	}
	return res.String(), nil
}
