package filter

import (
	"fmt"
	"strconv"
	"strings"
)

// -A - "after" печатать +N строк после совпадения
// -B - "before" печатать +N строк до совпадения
// -C - "context" (A+B) печатать ±N строк вокруг совпадения
// -c - "count" (количество строк)
// -i - "ignore-case" (игнорировать регистр)
// -v - "invert" (вместо совпадения, исключать)
// -F - "fixed", точное совпадение со строкой, не паттерн
// -n - "line num", печатать номер строки

const (
	colorGreen = "\033[32m"
	colorReset = "\033[0m"
)

type FilterOptions struct {
	Pattern string

	LinesBefore int
	LinesAfter  int

	Count      bool
	IgnoreCase bool
	Invert     bool
	Fixed      bool
	WithNums   bool
}

type shouldOutput struct {
	flags       []bool
	linesBefore int
	linesAfter  int
}

func (s shouldOutput) markShouldOutput(idx int) {
	start := max(0, idx-s.linesBefore)
	stop := min(len(s.flags), idx+s.linesAfter+1)
	for i := start; i < stop; i++ {
		s.flags[i] = true
	}
}

// окрашивает найденную подстроку, сделано плохо -- окрашивает только первое вхождение
func formatMatchInString(originalString string, idx int, fragmentLength int) string {
	return originalString[:idx] + colorGreen + originalString[idx:idx+fragmentLength] + colorReset + originalString[idx+fragmentLength:]
}

func Filter(inputString string, opts FilterOptions) string {

	opts.LinesAfter = max(0, opts.LinesAfter)
	opts.LinesBefore = max(0, opts.LinesBefore)

	if opts.IgnoreCase {
		opts.Pattern = strings.ToLower(opts.Pattern)
	}

	lines := strings.Split(inputString, "\n")
	shouldOutput := shouldOutput{
		flags: make([]bool, len(lines)),
	}

	if !opts.Invert {
		shouldOutput.linesBefore = opts.LinesBefore
		shouldOutput.linesAfter = opts.LinesAfter
	}

	var deciderFn func(string) bool
	if opts.Fixed {
		deciderFn = func(s string) bool {
			return s == opts.Pattern
		}
	} else {
		deciderFn = func(s string) bool {
			return strings.Contains(s, opts.Pattern)
		}
	}

	for i, elem := range lines {

		checkStr := elem

		if opts.IgnoreCase {
			checkStr = strings.ToLower(elem)
		}

		// (match) XOR (invert)
		if deciderFn(checkStr) != opts.Invert {
			shouldOutput.markShouldOutput(i)

			if !opts.Invert {
				idx := strings.Index(checkStr, opts.Pattern)
				lines[i] = formatMatchInString(elem, idx, len(opts.Pattern))
			}
		}
	}

	count := 0
	for _, b := range shouldOutput.flags {
		if b {
			count++
		}
	}

	if opts.Count {
		return fmt.Sprintf("matched %d lines", count)
	}

	outStrings := make([]string, 0, count)
	for i, b := range shouldOutput.flags {
		if b {
			strToInsert := lines[i]
			if opts.WithNums {
				strToInsert = strconv.Itoa(i+1) + ":" + strToInsert
			}
			outStrings = append(outStrings, strToInsert)
		}
	}

	return strings.Join(outStrings, "\n")
}
