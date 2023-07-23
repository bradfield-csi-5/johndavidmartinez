package xkcdcore

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

var COMIC_STORAGE_DIR string

type ComicInfo struct {
	Month      string `json:"month"`
	Id         int    `json:"num"`
	Link       string `json:"link"`
	Year       string `json:"year"`
	News       string `json:"new"`
	SafeTitle  string `json:"safe_title"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
	Img        string `json:"img"`
	Title      string `json:"title"`
	Day        string `json:"day"`
}

type ComicLoaderBuilder struct {
	client         *GetComicAPI
	cacheDirectory string
	comicIdx       int
}

func newComicLoaderBuilder() *ComicLoaderBuilder {
	return &ComicLoaderBuilder{}
}

func (b *ComicLoaderBuilder) withClient(client *GetComicAPI) *ComicLoaderBuilder {
	b.client = client
	return b
}

func (b *ComicLoaderBuilder) withCacheDirectory(cacheDirectory string) *ComicLoaderBuilder {
	b.cacheDirectory = cacheDirectory
	return b
}

func mkOrGetComicStorageDir() (string, error) {
	if COMIC_STORAGE_DIR != "" {
		return COMIC_STORAGE_DIR, nil
	}
	comicStorageDir := "/.local/xkcd_cmd/"
	homeDir, err := os.UserHomeDir()
	if err != nil && !os.IsExist(err) {
		return "", err
	}
	COMIC_STORAGE_DIR = homeDir + comicStorageDir
	err = os.Mkdir(COMIC_STORAGE_DIR, 0777)
	if err != nil && !os.IsExist(err) {
		return COMIC_STORAGE_DIR, err
	}
	return COMIC_STORAGE_DIR, nil
}

func (b *ComicLoaderBuilder) build() (*ComicLoader, error) {
	// Set defaults
	if b.client == nil {
		b.client = defaultGetComicApi()
	}
	if b.cacheDirectory == "" {
		defaultDir, err := mkOrGetComicStorageDir()
		if err != nil {
			return &ComicLoader{}, err
		}
		b.cacheDirectory = defaultDir
	}

	return &ComicLoader{
		comicInfo:      &ComicInfo{},
		client:         b.client,
		cacheDirectory: b.cacheDirectory,
		comicIdx:       b.comicIdx,
	}, nil
}

type ComicLoader struct {
	comicInfo      *ComicInfo
	client         *GetComicAPI
	cacheDirectory string
	comicIdx       int
}

func (cl *ComicLoader) getComicInfo(idx int) (*ComicInfo, error) {
	comicIdx := strconv.Itoa(idx)
	comicInfo := &ComicInfo{}
	data, err := os.ReadFile(cl.cacheDirectory + "/" + comicIdx)
	// If file is read from cache without error
	if err == nil {
		// Unmarshal cached file
		err := json.Unmarshal(data, comicInfo)
		if err != nil {
			// Sanity check title exists else fall through
			if comicInfo.Title != "" {
				return comicInfo, nil
			}
		}
		// File is not present or corrupted. Fall though to internet load
	}
	data, err = cl.client.getComic(idx)
	if err != nil {
		return comicInfo, fmt.Errorf("getComic Failed: %w", err)
	}
	err = json.Unmarshal(data, comicInfo)
	if err != nil {
		return comicInfo, fmt.Errorf("getComic Unmarshal failed: %w", err)
	}
	// Flush to disk cache
	go func() {
		writeData, err := json.MarshalIndent(comicInfo, "", " ")
		if err == nil {
			ioutil.WriteFile(cl.cacheDirectory+"/"+comicIdx, writeData, 0644)
		}
	}()
	return comicInfo, nil
}
