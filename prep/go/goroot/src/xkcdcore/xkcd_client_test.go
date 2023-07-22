package xkcdcore

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestGetOK(t *testing.T) {
	data, err := os.ReadFile("./test/xkcdComic1.json")
	if err != nil {
		log.Fatal("Bug in Test")
	}
	// make client return a copy of this data
	dataBody := make([]byte, len(data))
	copy(dataBody, data)

	timesCalled := 0
	expectedTimesCalled := 1

	// branching in this function
	// to simulate failures etc
	client := NewTestClient(func(req *http.Request) *http.Response {
		timesCalled += 1
		// Test request parameters
		equals(t, req.URL.String(), fmt.Sprintf(GET_COMIC_FORMAT, 1))
		return &http.Response{
			StatusCode: http.StatusOK,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBuffer(dataBody)),
		}
	})

	// retry 3 times sleep 0 (fast test)
	api := GetComicAPI{client, 3, 0}
	body, err := api.getComic(1)
	ok(t, err)
	equals(t, data, body)
	equals(t, timesCalled, expectedTimesCalled)
}

// test it does not retry
func TestGetNOTFOUND(t *testing.T) {
	data, err := os.ReadFile("./test/xkcdComic1.json")
	if err != nil {
		log.Fatal("Bug in Test")
	}
	// make client return a copy of this data
	dataBody := make([]byte, len(data))
	copy(dataBody, data)

	timesCalled := 0
	expectedTimesCalled := 1

	client := NewTestClient(func(req *http.Request) *http.Response {
		timesCalled += 1
		// Test request parameters
		equals(t, req.URL.String(), fmt.Sprintf(GET_COMIC_FORMAT, 404))
		return &http.Response{
			StatusCode: http.StatusNotFound,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBufferString(`NOTFOUND`)),
		}
	})

	// retry 3 times sleep 0 (fast test)
	api := GetComicAPI{client, 3, 0}
	_, err = api.getComic(404)
	assert(t, strings.Contains(err.Error(), "get comic: not found"))
	equals(t, timesCalled, expectedTimesCalled)
}

func TestTransientGetFailure(t *testing.T) {
	data, err := os.ReadFile("./test/xkcdComic1.json")
	if err != nil {
		log.Fatal("Bug in Test")
	}
	// make client return a copy of this data
	dataBody := make([]byte, len(data))
	copy(dataBody, data)

	timesCalled := 0
	expectedTimesCalled := 2

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

	// retry 3 times sleep 0 (fast test)
	api := GetComicAPI{client, 3, 0}
	body, err := api.getComic(1)
	ok(t, err)
	equals(t, data, body)
	equals(t, timesCalled, expectedTimesCalled)
}

func TestFatalGetFailure(t *testing.T) {
	data, err := os.ReadFile("./test/xkcdComic1.json")
	if err != nil {
		log.Fatal("Bug in Test")
	}
	// make client return a copy of this data
	dataBody := make([]byte, len(data))
	copy(dataBody, data)

	timesCalled := 0
	expectedTimesCalled := 3

	// branching in this function
	// to simulate failures etc
	client := NewTestClient(func(req *http.Request) *http.Response {
		timesCalled += 1
		// Test request parameters
		equals(t, req.URL.String(), fmt.Sprintf(GET_COMIC_FORMAT, 1))

		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBufferString("err")),
		}
	})

	// retry 3 times sleep 0 (fast test)
	api := GetComicAPI{client, 3, 0}
	_, err = api.getComic(1)
	assert(t, err != nil)
	equals(t, timesCalled, expectedTimesCalled)
}

func assert(tb testing.TB, val bool) {
	if !val {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: assert failed: %v\033[39m\n\n", filepath.Base(file), line, val)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}
