package mysort

import (
	"cmp"
	"errors"
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
)

/*
=== Утилита sort ===

Отсортировать строки (man sort)
Основное

# Поддержать ключи

-k — указание колонки для сортировки
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

# Дополнительное

# Поддержать ключи

-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учётом суффиксов

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/
var (
	ErrInvalidMonth  = errors.New("invalid month format")
	ErrInvalidSuffix = errors.New("invalid suffix")
)

type ErrNoColumn struct {
	column int
	line   int
}

func (e ErrNoColumn) Error() string {
	return fmt.Sprintf("there is no column %d at line %d", e.column, e.line)
}

type ErrNotSorted struct {
	lineNotInOrder int
}

func (e ErrNotSorted) Error() string {
	return fmt.Sprintf(" line #%d not in right order", e.lineNotInOrder)
}

type ErrLineDuplicate struct {
	numberLine int
	line       string
}

func (e ErrLineDuplicate) Error() string {
	return fmt.Sprintf("line %d: \"%s\" is duplicate", e.numberLine, e.line)
}

type SortType uint8

const (
	_ SortType = iota
	Alphabetical
	ByMonth
	Numeric
	HumanNumberic
)

var sortTypeNames = []string{
	"unknown",
	"alphabetical",
	"by month",
	"numeric",
	"human numeric",
}

func (s SortType) String() string {
	return sortTypeNames[s]
}

type SortOptions struct {
	SortType             SortType
	ReverseOrder         bool
	Unique               bool
	IgnoreTrailingSpaces bool
	CheckSorted          bool
	ByColumn             bool
	ColumnNumber         int
	Delim                string
}

const (
	isSortedMsg = "is sorted!"
)

func SortLinesInString(s string, opt SortOptions) (string, error) {
	originalLines := strings.Split(s, "\n")

	// Оставляяем только уникальные строки
	// и, если надо, проверяем нет ли дубликатов
	if opt.Unique {
		linesMap := make(map[string]struct{})
		for _, line := range originalLines {
			linesMap[line] = struct{}{}
		}
		newLines := make([]string, 0, len(linesMap))
		for i, line := range originalLines {
			if _, ok := linesMap[line]; ok {
				newLines = append(newLines, line)
				delete(linesMap, line)
			} else if opt.CheckSorted {
				return "", ErrLineDuplicate{
					numberLine: i,
					line:       line,
				}
			}
		}

		originalLines = newLines
	}

	// Значения по которым можно сортировать - не обязательно строка целиком
	// это могут быть значения столбцов.
	// Если сортировка должна производится не в алфавитном порядке,
	// в дальнейшем мы получим нужные велечины из строк.
	var stringValuesToSort []string

	if opt.ByColumn {
		// если нужно сортировать по столбцу
		sortColumns := make([]string, 0, len(originalLines))
		for i, line := range originalLines {
			columns := strings.Split(line, opt.Delim)
			if len(columns) <= opt.ColumnNumber {
				return "", ErrNoColumn{
					column: opt.ColumnNumber,
					line:   i,
				}
			}

			sortColumns = append(sortColumns, columns[opt.ColumnNumber])
		}

		stringValuesToSort = sortColumns
	} else {
		// если нужно сортировать строки целиком
		stringValuesToSort = slices.Clone(originalLines)
	}

	// делаем сортировку по выбранному типу
	// к сожалению, не придумал, как отрефакторить в более DRY-код
	switch opt.SortType {
	case Alphabetical:
		if opt.IgnoreTrailingSpaces {
			trimmedStringsToSort := make([]string, 0, len(stringValuesToSort))
			for _, str := range stringValuesToSort {
				trimmedStringsToSort = append(trimmedStringsToSort, strings.TrimRight(str, " "))
			}
			stringValuesToSort = trimmedStringsToSort
		}

		if opt.CheckSorted {
			return checkSorted(stringValuesToSort, opt.ReverseOrder)
		}

		sortSlicePair(stringValuesToSort, originalLines)

	case ByMonth:
		months, err := makeReferenceSlice(parseMonth, stringValuesToSort)
		if err != nil {
			return "", err
		}

		if opt.CheckSorted {
			return checkSorted(months, opt.ReverseOrder)
		}

		sortSlicePair(months, originalLines)

	case Numeric:
		numbers, err := makeReferenceSlice(strconv.Atoi, stringValuesToSort)
		if err != nil {
			return "", err
		}

		if opt.CheckSorted {
			return checkSorted(numbers, opt.ReverseOrder)
		}

		sortSlicePair(numbers, originalLines)

	case HumanNumberic:
		numbers, err := makeReferenceSlice(parseNumberWithSuffix, stringValuesToSort)
		if err != nil {
			return "", err
		}

		if opt.CheckSorted {
			return checkSorted(numbers, opt.ReverseOrder)
		}

		sortSlicePair(numbers, originalLines)
	}

	if opt.ReverseOrder {
		slices.Reverse(originalLines)
	}

	return strings.Join(originalLines, "\n"), nil
}

func checkSorted[T cmp.Ordered](linesToSort []T, reverseOrder bool) (string, error) {
	if reverseOrder {
		for i := 1; i < len(linesToSort); i++ {
			if linesToSort[i-1] < linesToSort[i] {
				return "", ErrNotSorted{
					lineNotInOrder: i,
				}
			}
		}

		return isSortedMsg, nil
	}

	for i := 1; i < len(linesToSort); i++ {
		if linesToSort[i-1] > linesToSort[i] {
			return "", ErrNotSorted{
				lineNotInOrder: i,
			}
		}
	}

	return isSortedMsg, nil

}

func sortSlicePair[R cmp.Ordered, M any](reference []R, mainSlice []M) {

	type tempSort struct {
		r R
		m M
	}
	ts := make([]tempSort, 0, len(reference))

	for i, r := range reference {
		ts = append(ts, tempSort{
			r: r,
			m: mainSlice[i],
		})
	}

	sort.Slice(ts, func(i, j int) bool {
		return ts[i].r < ts[j].r
	})

	for i, unit := range ts {
		mainSlice[i] = unit.m
	}
}

func makeReferenceSlice[T any](convFunc func(s string) (T, error), linesToSort []string) ([]T, error) {
	referenceSlice := make([]T, 0, len(linesToSort))
	for _, numStr := range linesToSort {
		num, err := convFunc(numStr)
		if err != nil {
			return nil, err
		}
		referenceSlice = append(referenceSlice, num)
	}
	return referenceSlice, nil
}

func parseNumberWithSuffix(s string) (float64, error) {

	number, err := strconv.Atoi(s)
	if err == nil {
		return float64(number), nil
	}

	number, err = strconv.Atoi(s[:len(s)-1])
	if err != nil {
		return 0, err
	}

	suffix := s[len(s)-1:]
	switch suffix {
	case "n":
		return float64(number) / 1_000_000_000, nil

	case "u":
		return float64(number) / 1_000_000, nil

	case "m":
		return float64(number) / 1_000, nil

	case "k", "K":
		return float64(number) * 1_000, nil

	case "M":
		return float64(number) * 1_000_000, nil

	case "G":
		return float64(number) * 1_000_000_000, nil

	case "T":
		return float64(number) * 1_000_000_000_000, nil
	}

	return 0, ErrInvalidSuffix

}

func parseMonth(s string) (time.Month, error) {
	switch strings.ToLower(s) {
	case "jan", "january":
		return time.January, nil
	case "feb", "february":
		return time.February, nil
	case "mar", "march":
		return time.March, nil
	case "apr", "april":
		return time.April, nil
	case "may":
		return time.May, nil
	case "jun", "june":
		return time.June, nil
	case "jul", "july":
		return time.July, nil
	case "aug", "august":
		return time.August, nil
	case "sep", "september":
		return time.September, nil
	case "oct", "october":
		return time.October, nil
	case "nov", "november":
		return time.November, nil
	case "dec", "december":
		return time.December, nil
	}

	return 0, ErrInvalidMonth
}

func main() {

}
