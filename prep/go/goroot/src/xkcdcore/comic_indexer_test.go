package xkcdcore

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

// utility function that populates
// the first 100 comics
// TODO make run only once
// TODO make directory not relative
func populateTestComics(t *testing.T) {
	testCacheDir, err := mkOrGetTestStorageDir(t.Name())
	// Real client to populate comics from the internet
	cl, err := newComicLoaderBuilder().withCacheDirectory(testCacheDir).build()
	if err != nil {
		log.Fatal("Bug in test populating test comics %w", err)
	}
	for i := 1; i < 100; i++ {
		comicInfo, err := cl.getComicInfo(i)
		if err != nil {
			log.Fatal("Bug in test populating test comics %w", err)
		}
		writeData, err := json.MarshalIndent(comicInfo, "", " ")
		if err != nil {
			log.Fatal("Bug in test populating test comics %w", err)
		}
		err = ioutil.WriteFile("./testcomics/"+strconv.Itoa(i), writeData, 0644)
		if err != nil {
			log.Fatal("Bug in test populating test comics %w", err)
		}
	}
}

// xkcd url is formatted https://xkcd.com/1/info.0.json"
func parseIdFromRequestURL(link string) int {
	linkParts := strings.Split(link, "/")
	id, err := strconv.Atoi(linkParts[3])
	if err != nil {
		log.Fatal("Bug in test parsing link from id %w", err)
	}
	return id
}

func loadTestComicData(idx int) []byte {
	data, err := os.ReadFile("./testcomics/" + strconv.Itoa(idx))
	if err != nil {
		log.Fatal("Bug in Test loading ComicTestData: %w", err)
	}
	return data
}

func TestIndexHappyPath(t *testing.T) {
	// populateTestComics(t)
	testCacheDir, err := mkOrGetTestStorageDir(t.Name())
	if err != nil {
		log.Fatal("Bug in test happy path test storage dir %w", err)
	}
	// clean index
	err = os.RemoveAll("./testindex")
	if err != nil && !os.IsNotExist(err) {
		log.Fatal("Bug removing testindex directory")
	}

	client := NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		id := parseIdFromRequestURL(req.URL.String())
		if id >= 100 {
			return &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(bytes.NewBufferString("nada")),
			}
		}

		data := loadTestComicData(id)

		return &http.Response{
			StatusCode: http.StatusOK,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBuffer(data)),
		}
	})

	api := GetComicAPI{client, 3, 0}
	cl, err := newComicLoaderBuilder().withClient(&api).withCacheDirectory(testCacheDir).build()
	if err != nil {
		log.Fatal("Bug in test creating comic loader %w", err)
	}
	assert(t, !fileExists("./testindex"))
	indexer := newComicIndex("./testindex", cl)
	err = indexer.index()
	ok(t, err)
	assert(t, fileExists("./testindex"))
}

// Benchmark entire index path
func BenchmarkIndexHappyPath(b *testing.B) {
	testCacheDir, err := mkOrGetTestStorageDir(b.Name())

	client := NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		id := parseIdFromRequestURL(req.URL.String())
		if id >= 100 {
			return &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(bytes.NewBufferString("nada")),
			}
		}

		data := loadTestComicData(id)

		return &http.Response{
			StatusCode: http.StatusOK,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBuffer(data)),
		}
	})

	api := GetComicAPI{client, 3, 0}
	cl, err := newComicLoaderBuilder().withClient(&api).withCacheDirectory(testCacheDir).build()
	if err != nil {
		log.Fatal("Bug in test creating comic loader %w", err)
	}
	indexer := newComicIndex("./testindex", cl)
	for i := 0; i < b.N; i++ {
		indexer.index()
	}
}

func TestComicIndexObject(t *testing.T) {
	comicInfo := &ComicInfo{}
	data, err := os.ReadFile("./test/xkcdComic1.json")
	if err != nil {
		log.Fatal("Bug in Test")
	}
	err = json.Unmarshal(data, comicInfo)
	if err != nil {
		log.Fatal("Bug in test")
	}

	obj := newComicIndexObject()
	obj.addComic(comicInfo)

	expectedTitleIndex := map[string]map[int]int{
		"barrel": {1: 1},
		"part":   {1: 1},
	}
	eq := reflect.DeepEqual(expectedTitleIndex, obj.TitleIndex)
	assert(t, eq)
	expectedBodyIndex := map[string]map[int]int{
		"boy":      {1: 2},
		"sits":     {1: 1},
		"barrel":   {1: 2},
		"which":    {1: 1},
		"floating": {1: 1},
		"ocean":    {1: 1},
		"wonder":   {1: 1},
		"where":    {1: 1},
		"float":    {1: 1},
		"next":     {1: 1},
		"drifts":   {1: 1},
		"into":     {1: 1},
		"distance": {1: 1},
		"nothing":  {1: 1},
		"don":      {1: 1}, // right now the scanner cuts apostrophe'd words
		"all":      {1: 1},
		"can":      {1: 1},
		"else":     {1: 1},
		"seen":     {1: 1},
	}

	eq = reflect.DeepEqual(expectedBodyIndex, obj.BodyIndex)
	assert(t, eq)
}

func TestSearchExactTerm(t *testing.T) {
	// populateTestComics(t)
	testCacheDir, err := mkOrGetTestStorageDir(t.Name())
	if err != nil {
		log.Fatal("Bug in test happy path test storage dir %w", err)
	}
	// clean index
	err = os.RemoveAll("./testindex")
	if err != nil && !os.IsNotExist(err) {
		log.Fatal("Bug removing testindex directory")
	}

	client := NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		id := parseIdFromRequestURL(req.URL.String())
		if id >= 100 {
			return &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(bytes.NewBufferString("nada")),
			}
		}

		data := loadTestComicData(id)

		return &http.Response{
			StatusCode: http.StatusOK,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBuffer(data)),
		}
	})

	api := GetComicAPI{client, 3, 0}
	cl, err := newComicLoaderBuilder().withClient(&api).withCacheDirectory(testCacheDir).build()
	if err != nil {
		log.Fatal("Bug in test creating comic loader %w", err)
	}
	assert(t, !fileExists("./testindex"))
	indexer := newComicIndex("./testindex", cl)
	err = indexer.index()
	ok(t, err)
	assert(t, fileExists("./testindex"))

	comicIds, err := indexer.search("Blown apart zipo", 1)
	ok(t, err)
	assert(t, contains(comicIds, 5))
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
