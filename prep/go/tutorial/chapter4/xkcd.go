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
        "sort"
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
    TranscriptSummary string `json:"transcript_summary"`
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
    // Some comics don't have transcripts
    if len(comic.Transcript) >= 100 {
        view.TranscriptSummary = comic.Transcript[:100]
    } else {
        view.TranscriptSummary = comic.Transcript
    }
    view.TranscriptSummary += "..."
    view.Link = fmt.Sprintf(XKCD_URL_FMT, comic.Id)
    return view
}

var XKCD_URL_JSON_FMT string 
var XKCD_URL_FMT string 
var COMIC_STORAGE_FMT string 
var COMIC_INDEX_STORAGE string


func mkComicStorageDir() string { 
    comicStorageDir := "/.local/xkcd_cmd/"
    homeDir, err := os.UserHomeDir()
    if err != nil {
        log.Fatal(err)
        os.Exit(1)
    }
    storageDir := strings.Join([]string{homeDir, comicStorageDir}, "")
    // TODO Just threw file permissions around to get stuff to work
    // Refresh understanding and make least-priviledge
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
    loadComics(&comics)

    // TODO make id lookup fast (only load single comic from disk)
    if err == nil && comicId > 0 {
        // if requested by number print comic
        if comicId > len(comics) {
            printComics([]int{}, &comics)
            os.Exit(0)
        }
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

type WeightedId struct {
    Id int
    Weight int
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

// Case insensitive search the term-based inverted index by matching exact terms
// The more terms that match the higher the weight
// For example "Three Laws of Robotics" will return the comics that
// Mention that phrase at the top of the list.
// Hardcode results to 3 since it can be comfortable viewed
// TODO make flag for result limit
// TODO (Algorithms) could use levelstien distance to search on typo'd terms
func searchIndex(terms []string, index ComicIndex) []int {
    weightedIds := make([]WeightedId, 0, 0)
    idsMapToWeight := make(map[int]int)
    ids := make([]int, 0, 3)

    for _, term := range terms {
        // ignore short terms
        if len(term) < 3 {
           continue
        }
        // If term is in index
        if _, value := index[strings.ToLower(term)]; value {
            // Add those comic ids
            for comicId, _ := range index[strings.ToLower(term)] {
                idsMapToWeight[comicId]++
            }
        } else {
            // Try to match term using levenshtein distance
            termMatch := ""
            termMatchDistance := 100
            for indexTerm, _ := range index {
                distance := lev(term, indexTerm)
                if (distance < termMatchDistance) {
                    termMatchDistance = distance
                    termMatch = indexTerm
                }
            }
            // Once term is matched add comic ids
            if termMatch != "" {
                for comicId, _ := range index[strings.ToLower(termMatch)] {
                    idsMapToWeight[comicId]++
                }
            }
        }
    }
    for id, weight := range idsMapToWeight {
        weightedIds = append(weightedIds, WeightedId{id, weight})
    }
    // TODO Stack overflowed this but it is pretty easy to understand
    // Worth reading more about sort module though
    sort.Slice(weightedIds, func(i, j int) bool {
        return weightedIds[i].Weight > weightedIds[j].Weight
    })
    idsToReturn := min(3, len(weightedIds))
    for i := 0; i < idsToReturn; i++ {
        ids = append(ids, weightedIds[i].Id)
    }
    return ids[:idsToReturn]
}

// Creates a simple case-insensitive inverted index of alphabetic terms
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
                // Skip words like "a" and "to"
                if len(term) < 3 {
                    continue
                }
                addComicIndex(index, term, comic.Id)
            }

            // TODO flush comics back to disk with Index = True so it
            // doesn't get indexed every run
            // Trying to do it as everything is written now introduces a bug
            // that isn't worth fixing until the performance is desired (profiling exercise? :))
            comic.Indexed = true 
        }
    }
    flushIndexToDisk(&index)
    return index
}

// TODO ChatGPT'd this. Need to understand writing custom scanners outside this exercise
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
        log.Fatal(err)
    }
}


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
        go tryLoadFromDiskFallbackToInternet(i, &(*comics)[i], &wg, &finished)

        // Run lots of go routines for the first 2500
        // Then slow down once we go past 2500
        // If we let this go nuts we may start to timeout
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
            // This is safe because even if multiple routines touch it
            // because all routines only flip false -> true; none go true -> false
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

// TODO handle network issues?
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
        *finished = true 
        wg.Done()
    } else {
        fmt.Printf("\n----\nohno\n----\n")
        os.Exit(1)
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

// ok ok, I implemented the recursive solution myself looking at the formula but it was slow
// I MIGHT have ripped this from someone else's leetcode solution <_<
func lev(word1 string, word2 string) int {
    m := len(word1)
    n := len(word2)
    if m == 0 {
        return n
    }
    if n == 0 {
        return m
    }
    
    dp := make([][]int, m+1)
    for i := 0; i <= m; i++ {
        dp[i] = make([]int, n+1)
    }
    
    for i := 0; i <= m; i ++ {
        dp[i][0] = i
    }
    for i := 0; i <= n; i ++ {
        dp[0][i] = i
    }
    
    
    for i := 1; i <= m ; i++ {
        for j := 1; j <= n; j++ {
            if word1[i - 1] == word2[j - 1] {
                dp[i][j] = dp[i - 1][j - 1]
            }else {
                dp[i][j] = minDis(dp[i -1][j], dp[i][j - 1], dp[i - 1][j - 1]) + 1
            }
        } 
    }
    
    return dp[len(word1)][len(word2)]

}

func minDis(a, b, c int) int {
    return min(a, min(b, c))
}

//TODO weights in index on title and summary
// so weight will be higher if term in title
