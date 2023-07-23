package xkcdcore

import (
	"strings"
	"sync"
	"sync/atomic"
)

type ComicIndex struct {
	indexDirectory string
	loader         *ComicLoader
}

type ComicIndexObject struct {
	titleIndex map[string]map[int]int
	bodyIndex  map[string]map[int]int
}

func newComicIndexObject() ComicIndexObject {
	obj := ComicIndexObject{}
	titleIndex := make(map[string]map[int]int)
	bodyIndex := make(map[string]map[int]int)
	obj.titleIndex = titleIndex
	obj.bodyIndex = bodyIndex
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

type ComicGetResult struct {
	comicInfo *ComicInfo
	success   bool
}

func newComicIndex(indexDirectory string, loader *ComicLoader) *ComicIndex {
	return &ComicIndex{indexDirectory, loader}
}

func normalize(word string) string {
	return strings.ToLower(word)
}

func (ci *ComicIndex) index() error {
	comicResults := make([]ComicGetResult, 0, 0)
	resultChannel := make(chan ComicGetResult)
	done := make(chan struct{})
	comicIdx := 1
	var wg sync.WaitGroup

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
	// for i := 0; i < len(comicResults); i++ {
	// 	if comicResults[i].success {
	// 		fmt.Println(comicResults[i].comicInfo.Title)
	// 	}
	// }
	// obj := newComicIndexObject()
	for _, result := range comicResults {
		if result.success {
			// iterate words in title
			// iterate words in body
		}
	}

	return nil
}

func (ci *ComicIndex) search(terms string) ([]int, error) {
	comicIds := make([]int, 1)
	return comicIds, nil
}
