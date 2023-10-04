package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/hablof/task-level-two/develop/dev06/cut"

	flag "github.com/spf13/pflag"
)

/*
=== Утилита cut ===

, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

var (
	ErrNoFields = errors.New("fields unspecified")
)

func parseFlags(progname string, args []string) (config *cut.CutOptions, output string, err error) {
	flags := flag.NewFlagSet(progname, flag.ContinueOnError)
	buf := bytes.Buffer{}
	flags.SetOutput(&buf)

	separated := flags.BoolP("separated", "s", false, "Do not print lines not containing delimiters")
	fields := flags.IntSliceP("fields", "f", []int{}, "Select only these fields on each line")
	delimeter := flags.StringP("delimeter", "d", "\t", "Use character sprcified instead of a tab for the field delimiter")

	err = flags.Parse(args)
	if err != nil {
		return nil, buf.String(), err
	}

	if len(*fields) <= 0 {
		buf.WriteString("you must specify a list of fields, use -f flag")
		return nil, buf.String(), ErrNoFields
	}

	cfg := cut.CutOptions{
		OnlyDelimited: *separated,
		Delimiter:     *delimeter,
		Fields:        *fields,
	}

	return &cfg, "", nil
}

func main() {

	cutOptions, output, err := parseFlags(os.Args[0], os.Args[1:])
	if err != nil {
		fmt.Println(err)
		if output != "" {
			fmt.Println(output)
		}

		os.Exit(1)
	}

	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	out := cut.Cut(string(b), *cutOptions)
	fmt.Println(out)
}
