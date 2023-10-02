package main

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
import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	mysort "github.com/hablof/task-level-two/develop/dev03/sort"

	flag "github.com/spf13/pflag"
)

var (
	ErrFileUnspecified = errors.New("filename unspecified")
)

type ErrSortTypeConflict struct {
	type1 mysort.SortType
	type2 mysort.SortType
}

// Error implements error.
func (e ErrSortTypeConflict) Error() string {
	return fmt.Sprintf("cannot use %s and %s sort at the same time", e.type1.String(), e.type2.String())
}

func getSortType(numeric, monthSort, humanNumeric bool) (mysort.SortType, error) {
	var sortType mysort.SortType

	if numeric {
		sortType = mysort.Numeric
	}

	if monthSort {
		if sortType != 0 {
			return 0, ErrSortTypeConflict{
				type1: sortType,
				type2: mysort.ByMonth,
			}
		} else {
			sortType = mysort.ByMonth
		}
	}

	if humanNumeric {
		if sortType != 0 {
			return 0, ErrSortTypeConflict{
				type1: sortType,
				type2: mysort.HumanNumberic,
			}
		} else {
			sortType = mysort.HumanNumberic
		}
	}

	if sortType == 0 {
		return mysort.Alphabetical, nil
	}

	return sortType, nil
}

func parseFlags(progname string, args []string) (config *mysort.SortOptions, filename string, output string, err error) {
	flags := flag.NewFlagSet(progname, flag.ContinueOnError)
	buf := bytes.Buffer{}
	flags.SetOutput(&buf)

	numeric := flags.BoolP("numeric", "n", false, "Sorts values as a numbers")
	reverse := flags.BoolP("reverse", "r", false, "Sorts values in reverse order")
	unique := flags.BoolP("unique", "u", false, "Suppresses all but one in each set of lines having equal keys")
	monthSort := flags.BoolP("month-sort", "M", false, "Compares as months")
	ignoreTrailingSpaces := flags.BoolP("ignore-trailing-spaces", "b", false, "Ignores trailing spaces")
	checkSorted := flags.BoolP("check-sorted", "c", false, "Check for sorted input")
	humanNumeric := flags.BoolP("human-numeric-sort", "h", false, "Compare human readable numbers (e.g., 2K 1G)")

	var column int
	flags.IntVarP(&column, "key", "k", -1, "Sets column to sort by")
	delim := flags.StringP("field-separator", "t", " ", "Uses specified column-separator")

	err = flags.Parse(args)
	if err != nil {
		return nil, "", buf.String(), err
	}

	sortType, err := getSortType(*numeric, *monthSort, *humanNumeric)
	if err != nil {
		return nil, "", "", err
	}

	if len(flags.Args()) == 0 {
		return nil, "", "", ErrFileUnspecified
	}

	cfg := mysort.SortOptions{
		SortType:             sortType,
		ReverseOrder:         *reverse,
		Unique:               *unique,
		IgnoreTrailingSpaces: *ignoreTrailingSpaces,
		CheckSorted:          *checkSorted,
		ByColumn:             column >= 0,
		ColumnNumber:         column,
		Delim:                *delim,
	}

	filename = flags.Args()[0]

	return &cfg, filename, "", nil
}

func main() {
	sortOpts, filename, output, err := parseFlags(os.Args[0], os.Args[1:])
	if err != nil {
		fmt.Println(err)
		if output != "" {
			fmt.Println(output)
		}

		os.Exit(1)
	}

	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	sortedLines, err := mysort.SortLinesInString(string(b), *sortOpts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(sortedLines)
}
