package crawler

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

type CrawlConfig struct {
	Site  string
	Depth int
}

type internalCrawlConfig struct {
	folder     string
	urlToCrawl *url.URL
}

func Crawl(cfg CrawlConfig) error {
	initUrl, err := url.Parse(cfg.Site)
	if err != nil {
		return err
	}

	folderName := initUrl.Host
	// if err := os.Mkdir(folderName, os.ModePerm); err != nil {
	// 	return err
	// }

	// mem := make(map[string]struct{})

	// icfg := internalCrawlConfig{
	// 	folder:     folderName,
	// 	urlToCrawl: initUrl,
	// }

	col := colly.NewCollector(
		colly.AllowedDomains(folderName),
		colly.MaxDepth(cfg.Depth),
	)

	col.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if strings.Contains(link, "#") ||
			strings.Contains(link, ".zip") ||
			strings.Contains(link, ".exe") ||
			strings.Contains(link, ".msi") ||
			strings.Contains(link, ".gz") ||
			strings.Contains(link, ".pkg") {
			return
		}
		// Print link
		// fmt.Printf("Link found: %s\n", link)
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		col.Visit(e.Request.AbsoluteURL(link))
	})

	col.OnResponse(func(r *colly.Response) {
		filename := folderName + strings.TrimRight(r.Request.URL.Path, "/") + ".html"
		pathFragments := strings.Split(filename, "/")
		dir := strings.Join(pathFragments[:len(pathFragments)-1], "/")
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			fmt.Println(err)
			fmt.Printf("dir spicified: %s\n", dir)
		}

		if _, err := os.Stat(filename); os.IsExist(err) {
			fmt.Println("Already exists: " + filename)
			return
		}

		if err := r.Save(filename); err != nil {
			fmt.Println(err)
			// return err
		} else {
			fmt.Println("File saved: " + filename)
		}
	})

	col.Visit(initUrl.String())
	col.Wait()

	return nil
}

func crawl(cfg internalCrawlConfig, mem map[string]struct{}) error {
	panic("unimlemented")
}
