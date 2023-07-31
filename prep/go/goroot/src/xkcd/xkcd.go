package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"xkcdcore"
)

const (
	INDEX_CMD    = "index"
	SEARCH_CMD   = "search"
	GET_CMD      = "get"
	XKCD_URL_FMT = "https://xkcd.com/%d/"
)

func usage() {
	fmt.Println("xkcd command - search for xkcd comics")
	fmt.Println("\tGenerate a index:          xkcd index")
	fmt.Println("\tGet comic by id:           xkcd get 3")
	fmt.Println("\tSearch for comic by terms: xkcd search \"blown apart\"")
}

type XKCDComicView struct {
	Id                int    `json:"id"`
	Title             string `json:"title"`
	TranscriptSummary string `json:"transcript_summary"`
	Link              string `json:"link"`
}

func toView(comic xkcdcore.ComicInfo) XKCDComicView {
	var view XKCDComicView
	view.Id = comic.Id
	view.Title = comic.Title
	// Some comics don't have transcripts
	// TODO OCR comics to extract transcripts??
	if len(comic.Transcript) >= 100 {
		view.TranscriptSummary = comic.Transcript[:100]
	} else {
		view.TranscriptSummary = comic.Transcript
	}
	view.TranscriptSummary += "..."
	view.Link = fmt.Sprintf(XKCD_URL_FMT, comic.Id)
	return view
}

func jsonify(comic xkcdcore.ComicInfo) (string, error) {
	data, err := json.MarshalIndent(toView(comic), "", " ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(0)
	}
	cmd := os.Args[1]

	switch cmd {
	case INDEX_CMD:
		indexer, err := xkcdcore.NewComicIndexBuilder().Build()
		if err != nil {
			fmt.Printf("Failed to index! %v\n", err)
			os.Exit(1)
		}
		err = indexer.Index()
		if err != nil {
			fmt.Printf("Failed to index! %v\n", err)
		}
		fmt.Println("Index successful!")
		os.Exit(0)
	case SEARCH_CMD:
		indexer, err := xkcdcore.NewComicIndexBuilder().Build()
		if err != nil {
			fmt.Printf("Failed to load index to search! %v\n", err)
			os.Exit(1)
		}
		ids, err := indexer.Search(strings.Join(os.Args[1:], " "), 3)
		if err != nil {
			fmt.Printf("Error searching index: %v", err)
		}
		cl, err := xkcdcore.NewComicLoaderBuilder().Build()
		if err != nil {
			fmt.Printf("Unexpected error! %v\n", err)
		}
		comics := make([]string, 0, 0)
		for i := 0; i < len(ids); i++ {
			info, err := cl.GetComicInfo(ids[i])
			if err != nil {
				fmt.Println("Error getting comic: %v\n", err)
				os.Exit(1)
			}
			s, err := jsonify(*info)
			if err != nil {
				fmt.Printf("Unexpected error! %v\n", err)
			}
			comics = append(comics, s)
		}
		fmt.Println(comics)
	case GET_CMD:
		if len(os.Args) != 3 {
			usage()
			os.Exit(0)
		}
		idx, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("err: get expects an integer")
			os.Exit(1)
		}
		cl, err := xkcdcore.NewComicLoaderBuilder().Build()
		if err != nil {
			fmt.Printf("Unexpected error! %v\n", err)
		}
		info, err := cl.GetComicInfo(idx)
		if err != nil {
			fmt.Printf("Error getting comic: %v\n", err)
			os.Exit(1)
		}

		s, err := jsonify(*info)
		if err != nil {
			fmt.Printf("Unexpected error! %v\n", err)
		}
		fmt.Println(s)
		os.Exit(0)
	default:
		usage()
		os.Exit(0)
	}
}
