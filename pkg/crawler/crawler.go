package crawler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zoturen/seearch/pkg/model"

	"golang.org/x/net/html"
)

type WebCrawler struct {
	UnvisitedLinks []string
	VisitedLinks   map[string]struct{}
}

func NewCrawler() *WebCrawler {
	return &WebCrawler{
		VisitedLinks:   model.CreateSet([]string{}),
		UnvisitedLinks: []string{},
	}
}

type CrawlerResult struct {
	Document model.Document
	Link     string
}

func (c *WebCrawler) AppendVisitedLink(url string) {
	model.AppendToSet(c.VisitedLinks, url)
}

func (c *WebCrawler) visitSite(url string) model.Document {
	fmt.Printf("Trying to visit: --- %s\n", url)
	res, err := http.Get(url)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer res.Body.Close()

	z := html.NewTokenizer(res.Body)
	document := model.Document{}
	document.Link = url

	isScriptText := false

	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			if z.Err().Error() == "EOF" {
				break
			}
		}
		if isScriptText {
			isScriptText = false
			continue
		}

		if tt == html.StartTagToken {
			name, attr := z.TagName()
			if name[0] == 'a' && attr {
				key, val, _ := z.TagAttr()
				if string(key) == "href" {
					value := string(val)
					if strings.HasPrefix(value, "#") {
						continue
					}

					if !strings.HasPrefix(value, "http") {
						if strings.HasSuffix(url, "/") {
							value = strings.TrimPrefix(value, "/")
						}
						value = fmt.Sprintf("%s%s", url, value)
					}
					c.UnvisitedLinks = append(c.UnvisitedLinks, value)

				}
			}

			if string(name) == "script" {
				isScriptText = true
				continue
			}
		}

		if tt == html.TextToken {
			document.Text += string(strings.ToLower(string(z.Text())))
		}
	}

	return document
}

func (c *WebCrawler) CrawlWeb(startUrl string, depth int, result chan CrawlerResult) {

	c.UnvisitedLinks = append(c.UnvisitedLinks, startUrl)
	i := 0
	for {
		if i > depth {
			break
		}
		if len(c.UnvisitedLinks) > 0 {
			link := c.UnvisitedLinks[0]
			c.UnvisitedLinks = append(c.UnvisitedLinks[:0], c.UnvisitedLinks[1:]...)
			if !model.ContainsSet(c.VisitedLinks, link) {
				document := c.visitSite(link)
				c.AppendVisitedLink(link)
				crawlerResult := CrawlerResult{
					Document: document,
					Link:     link,
				}

				result <- crawlerResult
			}
		} else {
			break
		}
		i++
	}

}
