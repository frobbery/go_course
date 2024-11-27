package hw02unpackstring

import (
	"errors"
	"strconv"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	runes := []rune(s)
	newRunes := make([]rune, 0, len(runes))
	var prevRune rune
	prevRuneNum := false
	prevRuneShielded := false
	for runeInd, currRune := range runes {
		//nolint:nestif
		if currRune == '\\' && !prevRuneShielded {
			prevRuneShielded = true
		} else {
			if repeatNum, err := strconv.ParseInt(string(currRune), 10, 32); !prevRuneShielded && err == nil {
				if !prevRuneShielded && (prevRuneNum || prevRune == 0) {
					return "", ErrInvalidString
				}
				prevRuneNum = true
				for i := int64(0); i < repeatNum; i++ {
					newRunes = append(newRunes, prevRune)
				}
			} else {
				if prevRuneShielded && err != nil && currRune != '\\' {
					return "", ErrInvalidString
				}
				if !prevRuneNum && prevRune != 0 {
					newRunes = append(newRunes, prevRune)
				}
				if runeInd == len(runes)-1 {
					newRunes = append(newRunes, currRune)
				}
				prevRuneNum = false
			}
			prevRune = currRune
			prevRuneShielded = false
		}
	}
	return string(newRunes), nil
}
