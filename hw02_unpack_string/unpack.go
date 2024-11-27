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
	var prevRuneNum bool = false
	for runeInd, currRune := range runes {
		if repeatNum, err := strconv.ParseInt(string(currRune), 10, 32); err == nil {
			if prevRuneNum || prevRune == 0 {
				return "", ErrInvalidString
			}
			prevRuneNum = true
			for i := int64(0); i < repeatNum; i++ {
				newRunes = append(newRunes, prevRune)
			}
		} else {
			if !prevRuneNum && prevRune != 0  {
				newRunes = append(newRunes, prevRune)
			}
			if runeInd == len(runes) - 1 {
				newRunes = append(newRunes, currRune)
			}
			prevRuneNum = false
		}
		prevRune = currRune
	}
	
	return string(newRunes), nil
}
