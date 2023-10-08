package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/hablof/task-level-two/develop/dev09/crawler"

	flag "github.com/spf13/pflag"
)

/*
=== Утилита wget ===

Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

const (
	maxDepth = 10
)

var (
	ErrTooDeep    = errors.New("depth is too big")
	ErrInvalidURL = errors.New("invalid url")
)
var (
	validURL = regexp.MustCompile(`^https?:\/\/(?:www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b(?:[-a-zA-Z0-9()@:%_\+.~#?&\/=]*)$`)
)

func parseFlags(progname string, args []string) (cfg *crawler.CrawlConfig, output string, err error) {
	flags := flag.NewFlagSet(progname, flag.ContinueOnError)
	buf := bytes.Buffer{}
	flags.SetOutput(&buf)

	depth := flags.IntP("recursive", "r", 5, "Selects recursive crawl depth")

	err = flags.Parse(args)
	if err != nil {
		return nil, buf.String(), err
	}

	urlCandidate := strings.TrimSpace(flags.Arg(0))
	if !validURL.MatchString(urlCandidate) {
		return nil, "", ErrInvalidURL
	}

	if *depth > maxDepth {
		buf.WriteString(fmt.Sprintf("maximum depth is %d", maxDepth))
		return nil, buf.String(), ErrTooDeep
	}

	cfg = &crawler.CrawlConfig{
		Site:  urlCandidate,
		Depth: *depth,
	}

	return cfg, "", nil
}

func main() {

	cfg, output, err := parseFlags(os.Args[0], os.Args[1:])
	if err != nil {
		fmt.Println(output)
		fmt.Println(err)
		os.Exit(1)
	}

	if err := crawler.Crawl(*cfg); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("well done")
}
