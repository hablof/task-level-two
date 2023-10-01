package main

/*
=== Задача на распаковку ===

Создать Go функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы / руны, например:
	- "a4bc2d5e" => "aaaabccddddde"
	- "abcd" => "abcd"
	- "45" => "" (некорректная строка)
	- "" => ""
Дополнительное задание: поддержка escape - последовательностей
	- qwe\4\5 => qwe45 (*)
	- qwe\45 => qwe44444 (*)
	- qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка функция должна возвращать ошибку. Написать unit-тесты.

Функция должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

import (
	"log"
	"regexp"
	"strconv"
	"strings"
)

type multiplyType uint8

const (
	_ multiplyType = iota
	multiplyLetter
	multiplyEscapeDigit
	multiplySlash
)

type multiply struct {
	startIndex   int
	endIndex     int
	changedPiece string
	// multiplyType multiplyType
}

var (
	validationRegex = regexp.MustCompile(`^\d+`)

	multiplyAllTypes = regexp.MustCompile(`[^0-9\\]\d+|\\\d\d*|\\\\\d+`)

	multiplyLetterRegex      = regexp.MustCompile(`[^0-9\\]\d+`)
	multiplyEscapeDigitRegex = regexp.MustCompile(`\\\d\d*`)
	multiplySlashRegex       = regexp.MustCompile(`\\\\\d+`)
)

func findUnpackSequenses(s string) string {
	if validationRegex.MatchString(s) {
		return ""
	}

	indexes := multiplyAllTypes.FindAllIndex([]byte(s), -1)
	multiplys := make([]multiply, 0, len(indexes))
	for _, indexPair := range indexes {
		piece := s[indexPair[0]:indexPair[1]]
		newMultiply := multiply{
			startIndex: indexPair[0],
			endIndex:   indexPair[1],
		}

		switch {
		case multiplyLetterRegex.MatchString(piece):

			pieceRunes := []rune(piece)
			symbol := string(pieceRunes[0])
			count, err := strconv.Atoi(string(pieceRunes[1:]))
			if err != nil {
				log.Printf("failed to parse %s: %v\n", piece, err)
			}

			newMultiply.changedPiece = strings.Repeat(symbol, count)

		case multiplyEscapeDigitRegex.MatchString(piece):
			symbol := piece[1:2]
			if len(piece) <= 2 {
				newMultiply.changedPiece = symbol
				break
			}

			count, err := strconv.Atoi(string(piece[2:]))
			if err != nil {
				log.Printf("failed to parse %s: %v\n", piece, err)
			}

			newMultiply.changedPiece = strings.Repeat(symbol, count)

		case multiplySlashRegex.MatchString(piece):
			count, err := strconv.Atoi(string(piece[2:]))
			if err != nil {
				log.Printf("failed to parse %s: %v\n", piece, err)
			}

			newMultiply.changedPiece = strings.Repeat("\\", count)
		}

		multiplys = append(multiplys, newMultiply)
	}

	totalLength := len(s)
	for _, mul := range multiplys {
		totalLength -= (mul.endIndex - mul.startIndex)
		totalLength += len(mul.changedPiece)
	}

	sb := strings.Builder{}
	sb.Grow(totalLength)
	lastIndexOldString := 0
	for _, mul := range multiplys {
		sb.WriteString(s[lastIndexOldString:mul.startIndex])
		sb.WriteString(mul.changedPiece)
		lastIndexOldString = mul.endIndex
	}
	sb.WriteString(s[lastIndexOldString:])

	return sb.String()
}
