package xkcdcore

import (
	"bufio"
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
	titleIndex   map[string]map[int]int
	bodyIndex    map[string]map[int]int
	indexVersion int
}

func newComicIndexObject() ComicIndexObject {
	obj := ComicIndexObject{}
	titleIndex := make(map[string]map[int]int)
	bodyIndex := make(map[string]map[int]int)
	obj.titleIndex = titleIndex
	obj.bodyIndex = bodyIndex
	obj.indexVersion = indexVersion
	return obj
}

func (obj *ComicIndexObject) addTitleWord(idx int, word string) {
	if obj.titleIndex[word] == nil {
		obj.titleIndex[word] = make(map[int]int)
	}
	obj.titleIndex[word][idx]++
}

func (obj *ComicIndexObject) addBodyWord(idx int, word string) {
	if obj.bodyIndex[word] == nil {
		obj.bodyIndex[word] = make(map[int]int)
	}
	obj.bodyIndex[word][idx]++
}

func (obj *ComicIndexObject) addComic(comic *ComicInfo) {
	scanner := bufio.NewScanner(strings.NewReader(comic.Title))
	scanner.Split(downcaseScanASCIIwords)
	for scanner.Scan() {
		term := scanner.Text()
		if shouldIndexWord(term) {
			obj.addTitleWord(comic.Id, term)
		}
	}
	scanner = bufio.NewScanner(strings.NewReader(comic.Transcript))
	scanner.Split(downcaseScanASCIIwords)
	for scanner.Scan() {
		term := scanner.Text()
		if shouldIndexWord(term) {
			obj.addBodyWord(comic.Id, term)
		}
	}
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

	return nil
}

func (ci *ComicIndex) search(terms string) ([]int, error) {
	comicIds := make([]int, 1)
	return comicIds, nil
}
