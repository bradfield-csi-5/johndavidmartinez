package xkcdcore

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

func cacheDir(name string) string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("Bug in test getting wd: %w", err)
	}
	return fmt.Sprintf("%s/testcache/%s", wd, name)
}

func mkOrGetTestStorageDir(name string) (string, error) {
	cleanTestStorageDir(name)

	err := os.MkdirAll(cacheDir(name), 0777)
	if err != nil && !os.IsExist(err) {
		log.Fatal("Bug in test making cacheDir: %w", err)
	}
	return cacheDir(name), nil
}

func cleanTestStorageDir(name string) {
	err := os.RemoveAll(cacheDir(name))
	if err != nil && !os.IsNotExist(err) {
		log.Fatal("Bug in test cleaning testStorageDir: %w", err)
	}
}

func cacheExists(name string, comicIdx string) bool {
	_, err := os.Stat(cacheDir(name) + "/" + comicIdx)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		log.Fatal("Bug in test stat on cache dir failed: %w", err)
	}
	return true
}

func TestHappyPath(t *testing.T) {
	testCacheDir, err := mkOrGetTestStorageDir(t.Name())

	if err != nil {
		log.Fatal("Bug in Test: %w", err)
	}
	data, err := os.ReadFile("./test/xkcdComic1.json")
	if err != nil {
		log.Fatal("Bug in Test: %w", err)
	}
	// make client return a copy of this data
	dataBody := make([]byte, len(data))
	copy(dataBody, data)

	// branching in this function
	// to simulate failures etc
	client := NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		equals(t, req.URL.String(), fmt.Sprintf(GET_COMIC_FORMAT, 1))
		return &http.Response{
			StatusCode: http.StatusOK,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBuffer(dataBody)),
		}
	})

	// First call should not hit cache
	api := GetComicAPI{client, 3, 0}
	cl, err := newComicLoaderBuilder().withClient(&api).withCacheDirectory(testCacheDir).build()
	// Assert cache does not exist before call
	assert(t, !cacheExists(t.Name(), "1"))
	comicInfo, err := cl.getComicInfo(1)
	ok(t, err)
	equals(t, comicInfo.Title, "Barrel - Part 1")
	// Assert cache exists after call
	// Sleep for 1 second since caching is done async
	time.Sleep(1 * time.Second)
	assert(t, cacheExists(t.Name(), "1"))
	// Second call should hit cache
	comicInfo, err = cl.getComicInfo(1)
	ok(t, err)
	equals(t, comicInfo.Title, "Barrel - Part 1")
}

// Test flaky client
// IT should get the comic and store it in cache directory
func TestFlakyClient(t *testing.T) {
	testCacheDir, err := mkOrGetTestStorageDir(t.Name())

	if err != nil {
		log.Fatal("Bug in Test: %w", err)
	}
	data, err := os.ReadFile("./test/xkcdComic1.json")
	if err != nil {
		log.Fatal("Bug in Test: %w", err)
	}
	// make client return a copy of this data
	dataBody := make([]byte, len(data))
	copy(dataBody, data)

	timesCalled := 0
	// branching in this function
	// to simulate failures etc
	client := NewTestClient(func(req *http.Request) *http.Response {
		timesCalled += 1
		// Test request parameters
		equals(t, req.URL.String(), fmt.Sprintf(GET_COMIC_FORMAT, 1))

		if timesCalled == 1 {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				// Send response to be tested
				Body: ioutil.NopCloser(bytes.NewBufferString("err")),
			}
		} else {
			return &http.Response{
				StatusCode: http.StatusOK,
				// Send response to be tested
				Body: ioutil.NopCloser(bytes.NewBuffer(dataBody)),
			}
		}
	})

	// First call should not hit cache
	api := GetComicAPI{client, 3, 0}
	cl, err := newComicLoaderBuilder().withClient(&api).withCacheDirectory(testCacheDir).build()
	// Assert cache does not exist before call
	assert(t, !cacheExists(t.Name(), "1"))
	comicInfo, err := cl.getComicInfo(1)
	ok(t, err)
	equals(t, comicInfo.Title, "Barrel - Part 1")
	// Assert cache exists after call
	// Sleep for 1 second since caching is done async
	time.Sleep(1 * time.Second)
	assert(t, cacheExists(t.Name(), "1"))
	// Second call should hit cache
	comicInfo, err = cl.getComicInfo(1)
	ok(t, err)
	equals(t, comicInfo.Title, "Barrel - Part 1")
}

// Test failed client
// It should return an error
func TestFailedClient(t *testing.T) {
	testCacheDir, err := mkOrGetTestStorageDir(t.Name())

	if err != nil {
		log.Fatal("Bug in Test: %w", err)
	}
	data, err := os.ReadFile("./test/xkcdComic1.json")
	if err != nil {
		log.Fatal("Bug in Test: %w", err)
	}
	// make client return a copy of this data
	dataBody := make([]byte, len(data))
	copy(dataBody, data)

	// branching in this function
	// to simulate failures etc
	client := NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		equals(t, req.URL.String(), fmt.Sprintf(GET_COMIC_FORMAT, 1))

		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBufferString("err")),
		}
	})

	// First call should not hit cache
	api := GetComicAPI{client, 3, 0}
	cl, err := newComicLoaderBuilder().withClient(&api).withCacheDirectory(testCacheDir).build()
	_, err = cl.getComicInfo(1)
	assert(t, err != nil)
}

// Test corrupt cache file
// It should ignore the corrupted file and re-write it
func TestCorruptCache(t *testing.T) {
	testCacheDir, err := mkOrGetTestStorageDir(t.Name())

	if err != nil {
		log.Fatal("Bug in Test: %w", err)
	}
	data, err := os.ReadFile("./test/xkcdComic1.json")
	if err != nil {
		log.Fatal("Bug in Test: %w", err)
	}
	// make client return a copy of this data
	dataBody := make([]byte, len(data))
	copy(dataBody, data)

	// branching in this function
	// to simulate failures etc
	client := NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		equals(t, req.URL.String(), fmt.Sprintf(GET_COMIC_FORMAT, 1))
		return &http.Response{
			StatusCode: http.StatusOK,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBuffer(dataBody)),
		}
	})

	// Corrupt cache
	writeData := []byte("{junk: \"o")
	err = ioutil.WriteFile(cacheDir(t.Name())+"/"+"1", writeData, 0644)
	if err != nil {
		log.Fatal("Bug in test corrupting cache %w", err)
	}

	api := GetComicAPI{client, 3, 0}
	cl, err := newComicLoaderBuilder().withClient(&api).withCacheDirectory(testCacheDir).build()
	comicInfo, err := cl.getComicInfo(1)
	ok(t, err)

	equals(t, comicInfo.Title, "Barrel - Part 1")
	// Assert cache is not corrupt after call
	// Sleep for 1 second since caching is done async
	time.Sleep(2 * time.Second)
	assert(t, cacheExists(t.Name(), "1"))
	cachedComicInfo := &ComicInfo{}
	cachedData, err := os.ReadFile(cacheDir(t.Name()) + "/" + "1")
	if err != nil {
		log.Fatal("Bug in test reading cache: %w", err)
	}
	err = json.Unmarshal(cachedData, cachedComicInfo)
	ok(t, err)
	equals(t, cachedComicInfo.Title, "Barrel - Part 1")
}
