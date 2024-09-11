package model

import (
	"math"
	"sort"
	"sync"

	"github.com/zoturen/seearch/pkg/lexer"
)

type Document struct {
	Link        string
	Description string
	Text        string
	TermFreq    map[string]int
	TermCount   int
}

func (doc *Document) ComputeTF(term string) float64 {
	return float64(doc.TermFreq[term]) / float64(doc.TermCount)
}

type IndexModel struct {
	Documents    map[string]Document // [link]document
	DocumentFreq map[string]int      // [term]freq
	mutex        sync.RWMutex
}

func (indexModel *IndexModel) ComputeIDF(term string) float64 {
	n := float64(len(indexModel.Documents))
	return math.Log10(n / float64(indexModel.DocumentFreq[term]))
}

func (indexModel *IndexModel) Search(entry string) []KV {
	tokens := lexer.NewLexer(entry).Tokens()
	result := map[string]float64{}
	for _, doc := range indexModel.Documents {
		var rank float64 = 0
		for _, token := range tokens {
			rank = rank + (doc.ComputeTF(token) * indexModel.ComputeIDF(token))

		}
		if rank > 0 {
			result[doc.Link] = rank
		}
	}
	return sortMapByValue(result)
}

type KV struct {
	Key   string
	Value float64
}

func sortMapByValue(m map[string]float64) []KV {
	ss := make([]KV, 0, len(m))
	for k, v := range m {
		ss = append(ss, KV{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	return ss
}

func (indexModel *IndexModel) AddDocument(doc Document) {
	indexModel.mutex.Lock()
	defer indexModel.mutex.Unlock()

	for t := range doc.TermFreq {
		if c, ok := indexModel.DocumentFreq[t]; ok {
			indexModel.DocumentFreq[t] = c + 1
		} else {
			indexModel.DocumentFreq[t] = 1
		}
	}

	indexModel.Documents[doc.Link] = doc
}

func CreateSet(slice []string) map[string]struct{} {
	set := make(map[string]struct{})
	for _, s := range slice {
		set[s] = struct{}{}
	}
	return set
}

func ContainsSet(set map[string]struct{}, str string) bool {
	_, ok := set[str]
	return ok
}

func AppendToSet(set map[string]struct{}, str string) {
	set[str] = struct{}{}
}
