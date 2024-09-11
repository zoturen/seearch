package main

import (
	"log"

	"github.com/zoturen/seearch/pkg/crawler"
	"github.com/zoturen/seearch/pkg/lexer"
	"github.com/zoturen/seearch/pkg/model"
	"github.com/zoturen/seearch/pkg/server"
)

func main() {
	indexModel := &model.IndexModel{}

	indexModel.DocumentFreq = map[string]int{}
	indexModel.Documents = map[string]model.Document{}
	crawlerResult := make(chan crawler.CrawlerResult)

	webCrawler := crawler.NewCrawler()

	go func() {
		webCrawler.CrawlWeb("https://github.com/zoturen", 10, crawlerResult)
	}()

	go func() {
		for cr := range crawlerResult {
			doc := cr.Document

			tokens := lexer.NewLexer(doc.Text).Tokens()
			doc.TermFreq = map[string]int{}

			for _, token := range tokens {
				if c, ok := doc.TermFreq[token]; ok {
					doc.TermFreq[token] = c + 1
				} else {
					doc.TermFreq[token] = 1
				}
				doc.TermCount++
			}
			log.Println("Added a new doc into index model")
			indexModel.AddDocument(doc)
		}
	}()

	server.NewHttpServer(indexModel).Serve()

}
