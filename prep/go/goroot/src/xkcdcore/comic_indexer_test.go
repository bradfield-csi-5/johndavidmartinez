package xkcdcore

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	err = indexer.index()
	ok(t, err)
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
