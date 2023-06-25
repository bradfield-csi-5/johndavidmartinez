package main

import (
        "os"
        "fmt"
        "strconv"
        "strings"
        "log"
        "encoding/json"
        "net/http"
        "io/ioutil"
        "sync"
        "unicode"
        "bufio"
)

type XKCDComic struct {
    Month string `json:"month"`
    Id int `json:"num"`
    Link string `json:"link"`
    Year string `json:"year"`
    News string `json:"new"`
    SafeTitle string `json:"safe_title"`
    Transcript string `json:"transcript"`
    Alt string `json:"alt"`
    Img string `json:"img"`
    Title string `json:"title"`
    Day string `json:"day"`
    Indexed bool 
}

type XKCDComicView struct {
    Id int `json:"id"`
    Title string `json:"title"`
    Transcript string `json:"transcript"`
    Link string `json:"link"`
}

type ComicIndex map[string]map[int]bool

func newComicIndex() ComicIndex {
    var index ComicIndex
    index = make(map[string]map[int]bool)
    return index
}

func addComicIndex(index ComicIndex, term string, id int) {
    if index[term] == nil {
        index[term] = make(map[int]bool)
    }
    index[term][id] = true
}


func toView(comic XKCDComic) XKCDComicView {
    var view XKCDComicView
    view.Id = comic.Id
    view.Title = comic.Title
    view.Transcript = comic.Transcript
    view.Link = fmt.Sprintf(XKCD_URL_FMT, comic.Id)
    return view
}


// First just write code to fetch and unmarshal
// then write code to lookup by id
// after that you can think more about indexing
var XKCD_URL_JSON_FMT string 
var XKCD_URL_FMT string 
var COMIC_STORAGE_FMT string 
var COMIC_INDEX_STORAGE string


// do something naive first match on words a word index
func mkComicStorageDir() string { 
    comicStorageDir := "/.local/xkcd_cmd/"
    homeDir, err := os.UserHomeDir()
    if err != nil {
        log.Fatal(err)
        os.Exit(1)
    }
    storageDir := strings.Join([]string{homeDir, comicStorageDir}, "")
    // TODO fix file permissions 
    err = os.Mkdir(storageDir, 0777)
    if err != nil && !os.IsExist(err) {
        log.Fatal(err)
        os.Exit(1)
    }
    storageDirFmt := strings.Join([]string{storageDir, "%d"}, "")
    return storageDirFmt
}

func main() {
    XKCD_URL_FMT = "https://xkcd.com/%d/"
    XKCD_URL_JSON_FMT = "https://xkcd.com/%d/info.0.json"
    COMIC_STORAGE_FMT = mkComicStorageDir()
    homeDir, err := os.UserHomeDir()
    if err != nil {
        log.Fatal(err)
        os.Exit(1)
    }
    COMIC_INDEX_STORAGE = homeDir + "/.local/xkcd_cmd/index"

    if len(os.Args) < 2 {
        fmt.Printf("Provide a comic by id or search terms\n")
        os.Exit(1)
    }

    comicId, err := strconv.Atoi(os.Args[1])
    if comicId == 404 {
        fmt.Printf("You should know better...\n")
        os.Exit(1)
    }

    // In practice we know 3000 is a better starting capacity
    // as xkcd has ~2790 comics but its nice to force
    // the resize logic to run for testing + learning purposes
    comics := make([]XKCDComic, 2000, 2000)
    // TODO make id lookup fast
    loadComics(&comics)

    if err == nil && comicId > 0 {
        // if requested by number print comic
        if comicId > len(comics) {
            fmt.Printf("Comic %v not found\n", comicId)
            os.Exit(0)
        }
        fmt.Printf("Looking up comic %d\n", comicId)
        printComics([]int{comicId}, &comics)
    } else {
        index := loadAndBuildIndex(&comics)
        ids := searchIndex(os.Args[1:], index)
        printComics(ids, &comics)
    }
}

func printComics(ids []int, comics *[]XKCDComic) {
    comicViews := make([]XKCDComicView, 0, 0)

    for _, v := range ids {
        view := toView((*comics)[v])
        comicViews = append(comicViews, view)
    }

    data, err := json.MarshalIndent(comicViews, "", "  ")
    if err != nil {
        log.Fatal(err)
        os.Exit(1)
    }
    fmt.Println(string(data))
}

func searchIndex(terms []string, index ComicIndex) []int {
    ids := make([]int, 0, 0)
    idsMap := make(map[int]bool)

    for _, term := range terms {
        for comicId, _ := range index[strings.ToLower(term)] {
            // put in map to dedup
            idsMap[comicId] = true
        }
    }
    // convert map into array
    for comicId, _ := range idsMap {
        ids = append(ids, comicId)
    }
    return ids
}


func loadAndBuildIndex(comics *[]XKCDComic) ComicIndex {
    indexData, err := os.ReadFile(COMIC_INDEX_STORAGE)
    if err != nil {
        if os.IsNotExist(err) {
            indexData = []byte("{}")
        } else {
            log.Fatal(err)
            os.Exit(1)
        }
    } 
    // map of map so ids don't duplicate
    index := newComicIndex()
    err = json.Unmarshal(indexData, &index)
    if err != nil {
        log.Fatal(err)
        os.Exit(1)
    }

    for _, comic := range *comics {
        if !comic.Indexed {
            // index title and transcript
            textToIndex := strings.Join([]string{comic.SafeTitle, comic.Title, comic.Transcript}, " ")
            scanner := bufio.NewScanner(strings.NewReader(textToIndex))
            scanner.Split(scanLetters)
            for scanner.Scan() {
                term := strings.ToLower(scanner.Text())
                addComicIndex(index, term, comic.Id)
            }

            comic.Indexed = true 
        }
    }
    flushIndexToDisk(&index)
    return index
}

// Yeah.. ChatGPT'd this one
// Should practice writing custom scanners outside of this exercise
func scanLetters(data []byte, atEOF bool) (advance int, token []byte, err error) {
    start := 0
    for ; start < len(data); start++ {
        if unicode.IsLetter(rune(data[start])) {
            break
        }
    }

    for i := start; i < len(data); i++ {
        if !unicode.IsLetter(rune(data[i])) {
            return i + 1, data[start:i], nil
        }
    }

    
    if !atEOF && unicode.IsLetter(rune(data[len(data)-1])) {
        return start, nil, nil
    }

    return len(data) + 1, data[start:], nil
}

func flushIndexToDisk(index *ComicIndex) {
    data, err := json.MarshalIndent(index, "", " ")
    if err != nil {
        log.Fatal(err)
    }

    err = ioutil.WriteFile(COMIC_INDEX_STORAGE, data, 0644)
    if err != nil {
        fmt.Printf("HERE\n")
        log.Fatal(err)
    }
}


// first version inverted word index?
// Probably fine for first attempt

func loadComics(comics *[]XKCDComic) {
    var wg sync.WaitGroup
    // Comics start at 1
    i := 1
    finished := false
    for {
        // Very funny xkcd.. Very funny..
        if i == 404 {
            i++
        }

        wg.Add(1)
        // id is incrementing and passed by value so it is safe
        // each routine will only access its index in the array
        //fmt.Fprintf(os.Stderr, "Loading comic %v\n", i)
        go tryLoadFromDiskFallbackToInternet(i, &(*comics)[i], &wg, &finished)

        // Run lots of go routines for the first 2000
        // Then slow down once we go past 2500
        // If we hit xkcd too aggressively we may start to timeout
        if i < 2500 {
            if i % 500 == 0 {
                wg.Wait()
            }
        } else {
            if i % 10 == 0 {
                wg.Wait()
            }
        }

        if finished {
            // A routine hit the end of the comic list
            // This is safe because even if multiple routines
            // flipped the flag they all would flipped it to true
            // (Nothing flips it to false, no race conditions here)
            break
        }

        // Resize comics array if running out of space
        if i + 500 > len(*comics) {
            // resize
            newComics := make([]XKCDComic, len(*comics) * 2, cap(*comics) * 2)
            copy(newComics, *comics)
            *comics = newComics
        }

        i++
    }
}

func tryLoadFromDiskFallbackToInternet(id int, xkcdComic *XKCDComic, wg *sync.WaitGroup, finished *bool) {
    data, err := os.ReadFile(fmt.Sprintf(COMIC_STORAGE_FMT, id))
    if err != nil {
        if os.IsNotExist(err) {
            go loadComicFromInternetAndFlushToDisk(id, xkcdComic, wg, finished)
        } else {
            log.Fatal(err)
            os.Exit(1)
        }
    } else {
        err = json.Unmarshal(data, xkcdComic)
        if err != nil {
            log.Fatal(err)
            os.Exit(1)
        }
        // Finished loading from file
        wg.Done()
    }
}

func loadComicFromInternetAndFlushToDisk(id int, xkcdComic *XKCDComic, wg *sync.WaitGroup, finished *bool) {
    url := fmt.Sprintf(XKCD_URL_JSON_FMT, id)
    resp, err := http.Get(url)
    if err != nil {
        log.Fatal(err)
    }

    defer resp.Body.Close()

    if resp.StatusCode == http.StatusOK {
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            log.Fatal(err)
        }
    
        err = json.Unmarshal(body, &xkcdComic)
        if err != nil {
            log.Fatal(err)
        }
        flushComicToDisk(id, xkcdComic, wg)
    } else if resp.StatusCode == http.StatusNotFound {
        // We've likely reached the end of available comics.
        // Commuicate that back to the caller via the boolean
        wg.Done()
        *finished = true 
    } else {
        fmt.Printf("\n----\nohno\n----\n")
    }
}

func flushComicToDisk(id int, xkcdComic *XKCDComic, wg *sync.WaitGroup) {
    data, err := json.MarshalIndent(xkcdComic, "", " ")
    if err != nil {
        log.Fatal(err)
    }

    err = ioutil.WriteFile(fmt.Sprintf(COMIC_STORAGE_FMT, id), data, 0644)
    if err != nil {
        log.Fatal(err)
    }
    // Finished will chain of events
    // WaitGroup is optional 
    if wg != nil {
        wg.Done()
    }
}

//
//func refreshAndLoadIndex(comics []XKCDComic) ComicIndex {
//    // load existing index
//    // refresh
//    index := new ComicIndex
//    return index
//}
//
//func fetchComic(number int) {
//}
//
//
//
//
//func searchComics(comics []XKCDComic, index ComicIndex) int[] {
//    return comics[0]
//}
//
//
//
//
//
