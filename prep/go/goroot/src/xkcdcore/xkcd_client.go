// XKCD client
// gets client from xkcd
// optionally uses disk to cache (comics don't change)
// defaults to xkcd url and cache true but can be overriden for testing

package xkcdcore

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const GET_COMIC_FORMAT = "https://xkcd.com/%d/info.0.json"

func defaultGetComicApi() *GetComicAPI {
	return &GetComicAPI{
		Client:     http.DefaultClient,
		maxRetry:   3,
		retrySleep: 3,
	}
}

type GetComicAPI struct {
	Client     *http.Client
	maxRetry   int
	retrySleep int
}

func (api *GetComicAPI) getComic(comicIdx int) ([]byte, error) {
	retryattempt := 0
	var resp *http.Response
	var err error
	for {
		if retryattempt == api.maxRetry {
			return nil, fmt.Errorf("get comic: Failed to fetch comic %v", comicIdx)
		}
		resp, err = api.Client.Get(fmt.Sprintf(GET_COMIC_FORMAT, comicIdx))
		if err == nil && (resp.StatusCode == http.StatusNotFound ||
			resp.StatusCode == http.StatusOK) {
			break
		}
		resp.Body.Close()
		retryattempt++
		time.Sleep(time.Duration(api.retrySleep) * time.Second)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("get comic: not found %v", comicIdx)
	}

	// Resp is ok
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

// Returns true if error returned is not found error
func isNotFound(err error) bool {
	return strings.Contains(err.Error(), "get comic: not found")
}

func Hello() string {
	return "HI"
}
