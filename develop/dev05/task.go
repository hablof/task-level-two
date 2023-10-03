package main

/*
=== Утилита grep ===

Реализовать утилиту фильтрации (man grep)

Поддержать флаги:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", печатать номер строки

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/hablof/task-level-two/develop/dev05/filter"

	flag "github.com/spf13/pflag"
)

var (
	ErrNotEnoughArgs = errors.New("not enough arguments, expected [pattern], [filename]")
)

func parseFlags(progname string, args []string) (config *filter.FilterOptions, filename string, output string, err error) {
	flags := flag.NewFlagSet(progname, flag.ContinueOnError)
	buf := bytes.Buffer{}
	flags.SetOutput(&buf)

	count := flags.BoolP("count", "c", false, "Returns the quantity of matched strings")
	ignoreCase := flags.BoolP("ignore-case", "i", false, "Ignores case")
	invert := flags.BoolP("invert", "v", false, "Invert the sense of matching, to select non-matching lines")
	fixed := flags.BoolP("fixed", "r", false, "Interpret patterns as fixed strings")
	withNums := flags.BoolP("line-num", "n", false, "Prefix each line of output with the 1-based line number within its input file")

	after := flags.IntP("after", "A", -1, "Print num lines of trailing context after matching lines")
	before := flags.IntP("before", "B", -1, "Print num lines of leading  context after matching lines")
	context := flags.IntP("context", "C", -1, "Print num lines of leading and trailing output context")

	err = flags.Parse(args)
	if err != nil {
		return nil, "", buf.String(), err
	}

	if len(flags.Args()) < 2 {
		return nil, "", "", ErrNotEnoughArgs
	}

	cfg := filter.FilterOptions{
		Pattern:     flags.Arg(0),
		LinesBefore: max(*before, *context),
		LinesAfter:  max(*after, *context),
		Count:       *count,
		IgnoreCase:  *ignoreCase,
		Invert:      *invert,
		Fixed:       *fixed,
		WithNums:    *withNums,
	}

	filename = flags.Arg(1)

	return &cfg, filename, "", nil
}

func main() {

	filterOptions, filename, output, err := parseFlags(os.Args[0], os.Args[1:])
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

	filteredLines := filter.Filter(string(b), *filterOptions)

	fmt.Println(filteredLines)
}
