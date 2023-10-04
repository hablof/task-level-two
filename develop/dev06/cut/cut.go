package cut

import "strings"

type CutOptions struct {
	OnlyDelimited bool
	Delimiter     string
	Fields        []int
}

func Cut(s string, opts CutOptions) string {
	lines := strings.Split(s, "\n")
	outputLines := make([]string, 0, len(lines))

	for _, line := range lines {
		hasDelimiter := strings.Contains(line, opts.Delimiter)
		if opts.OnlyDelimited && !hasDelimiter {
			continue
		}

		if !opts.OnlyDelimited && !hasDelimiter {
			outputLines = append(outputLines, line)
			continue
		}

		fields := strings.Split(line, opts.Delimiter)
		outFields := make([]string, 0, len(opts.Fields))
		for _, i := range opts.Fields {
			if i-1 < len(fields) && i-1 >= 0 {
				outFields = append(outFields, fields[i-1])
			}
		}

		outputLines = append(outputLines, strings.Join(outFields, opts.Delimiter))
	}

	return strings.Join(outputLines, "\n")
}
