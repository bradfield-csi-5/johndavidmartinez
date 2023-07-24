package xkcdcore

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
)

type IndexMethod int64

const (
	LeveinstenComparison IndexMethod = iota
	PrefixTree
	Trigram
	LeveinstenAutomata
)

const (
	indexVersion = 1
)

type ComicIndex struct {
	indexDirectory string
	loader         *ComicLoader
	indexMethod    IndexMethod
}

type ComicIndexObject struct {
	TitleIndex   map[string]map[int]int `json:"titleIndex"`
	BodyIndex    map[string]map[int]int `json:"bodyIndex"`
	IndexVersion int                    `json:"indexVersion"`
}

func newComicIndexObject() ComicIndexObject {
	obj := ComicIndexObject{}
	titleIndex := make(map[string]map[int]int)
	bodyIndex := make(map[string]map[int]int)
	obj.TitleIndex = titleIndex
	obj.BodyIndex = bodyIndex
	obj.IndexVersion = indexVersion
	return obj
}

func (obj *ComicIndexObject) addTitleWord(idx int, word string) {
	if obj.TitleIndex[word] == nil {
		obj.TitleIndex[word] = make(map[int]int)
	}
	obj.TitleIndex[word][idx]++
}

func (obj *ComicIndexObject) addBodyWord(idx int, word string) {
	if obj.BodyIndex[word] == nil {
		obj.BodyIndex[word] = make(map[int]int)
	}
	obj.BodyIndex[word][idx]++
}

func (obj *ComicIndexObject) addComic(comic *ComicInfo) {
	for _, term := range normalizeWordsToDowncaseASCII(comic.Title) {
		obj.addTitleWord(comic.Id, term)
	}
	for _, term := range normalizeWordsToDowncaseASCII(comic.Transcript) {
		obj.addBodyWord(comic.Id, term)
	}
}

func normalizeWordsToDowncaseASCII(sentence string) []string {
	words := make([]string, 0, 0)
	scanner := bufio.NewScanner(strings.NewReader(sentence))
	scanner.Split(downcaseScanASCIIwords)
	for scanner.Scan() {
		term := scanner.Text()
		if shouldIndexWord(term) {
			words = append(words, term)
		}
	}
	return words
}

type ComicGetResult struct {
	comicInfo *ComicInfo
	success   bool
}

func newComicIndex(indexDirectory string, loader *ComicLoader) *ComicIndex {
	return &ComicIndex{indexDirectory, loader, LeveinstenComparison}
}

func normalize(word string) string {
	return strings.ToLower(word)
}

func downcaseASCIILetter(b byte) byte {
	if b >= 65 && b <= 90 {
		return b + 32
	}
	return b
}

func isDowncasedASCIILetter(b byte) bool {
	return b >= 97 && b <= 122
}

func shouldIndexWord(word string) bool {
	if len(word) <= 2 {
		return false
	}
	if word == "the" || word == "and" || word == "alt" {
		return false
	}
	return true
}

// Scanner that takes an array of bytes and return only the downcased ASCII words
// XKCD comics are in english so we can get away with ignoring non-ASCII characters
// Does not handle punctuationed words like "don't"
func downcaseScanASCIIwords(data []byte, atEOF bool) (advance int, token []byte, err error) {
	start := 0
	// Find the start of the next word
	for ; start < len(data); start++ {
		downcased := downcaseASCIILetter(data[start])
		if isDowncasedASCIILetter(downcased) {
			break
		}
	}

	// travel bytes until a non-ASCII letter is found
	// downcasing letters as we go
	for i := start; i < len(data); i++ {
		data[i] = downcaseASCIILetter(data[i])
		if !isDowncasedASCIILetter(data[i]) {
			return i + 1, data[start:i], nil
		}
	}

	return len(data) + 1, data[start:], nil
}

// returns number of comics indexed
func (ci *ComicIndex) index() error {
	comicResults := make([]ComicGetResult, 0, 0)
	resultChannel := make(chan ComicGetResult)
	done := make(chan struct{})
	comicIdx := 1
	var wg sync.WaitGroup

	// Get all the comics
	go func() {
		for {
			select {
			case result := <-resultChannel:
				comicResults = append(comicResults, result)
			case <-done:
				return
			}
		}
	}()

	var errCnt int32
	for {
		if comicIdx%10 == 0 {
			wg.Wait()
			if atomic.LoadInt32(&errCnt) > 10 {
				done <- struct{}{}
				defer close(done)
				defer close(resultChannel)
				break
			}
		}

		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			comic, err := ci.loader.getComicInfo(idx)
			success := err == nil
			if !success {
				atomic.AddInt32(&errCnt, 1)
			}
			resultChannel <- ComicGetResult{comic, success}
		}(comicIdx)
		comicIdx++
	}

	// Create inverted index of comics
	obj := newComicIndexObject()
	for _, result := range comicResults {
		if result.success {
			obj.addComic(result.comicInfo)
		}
	}

	data, err := json.MarshalIndent(obj, "", " ")
	if err != nil {
		return fmt.Errorf("Error Creating Index: %w", err)
	}
	err = os.MkdirAll(ci.indexDirectory, 0777)
	if err != nil {
		return fmt.Errorf("Error making directory for Index: %w", err)
	}
	err = ioutil.WriteFile(ci.indexDirectory+"/"+"index", data, 0644)
	if err != nil {
		return fmt.Errorf("Error writing index: %w", err)
	}

	return nil
}

type SearchResult struct {
	comicId int
	weight  int
}

type SearchResults struct {
	results map[int]SearchResult
}

func (sr *SearchResults) toSortedArray() []int {
	results := make([]SearchResult, 0, 0)
	for _, v := range sr.results {
		results = append(results, v)
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].weight > results[i].weight
	})
	ids := make([]int, 0, 0)
	for _, v := range results {
		ids = append(ids, v.comicId)
	}
	return ids
}

// Fine for small lists
func filterStrings(a []string, b []string) []string {
	filtered := make([]string, 0, 0)
	index := make(map[string]int)
	for _, s := range b {
		index[s]++
	}
	for _, s := range a {
		if _, ok := index[s]; !ok {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

func (ci *ComicIndex) search(searchSentence string, limit int) ([]int, error) {
	// perform same processing on search terms as we do on indexing
	terms := normalizeWordsToDowncaseASCII(searchSentence)
	searchResults := SearchResults{}
	searchResults.results = make(map[int]SearchResult)
	obj := ComicIndexObject{}
	data, err := ioutil.ReadFile(ci.indexDirectory + "/" + "index")
	if err != nil {
		return nil, fmt.Errorf("Error reading index: %w", err)
	}
	err = json.Unmarshal(data, &obj)
	if err != nil {
		return nil, fmt.Errorf("Error reading index: %w", err)
	}

	// Try exact matching first
	foundExactTerms := make([]string, 0, 0)
	for _, term := range terms {
		if titleTermMap, ok := obj.TitleIndex[term]; ok {
			foundExactTerms = append(foundExactTerms, term)
			for id, _ := range titleTermMap {
				result, ok := searchResults.results[id]
				if !ok {
					searchResults.results[id] = SearchResult{id, 0}
					result = searchResults.results[id]
				}
				result.weight++
			}
		}
		if bodyTermMap, ok := obj.BodyIndex[term]; ok {
			for id, _ := range bodyTermMap {
				result, ok := searchResults.results[id]
				if !ok {
					searchResults.results[id] = SearchResult{id, 0}
					result = searchResults.results[id]
				}
				result.weight++
			}
		}
	}
	// Filter out terms that had exact matches. No need to search those further.
	terms = filterStrings(terms, foundExactTerms)
	if len(terms) == 0 {
		return searchResults.toSortedArray()[:limit], nil
	}

	switch ci.indexMethod {
	case LeveinstenComparison:
		return searchResults.toSortedArray(), nil
	default:
		return nil, fmt.Errorf("Bad index method")
	}
}
