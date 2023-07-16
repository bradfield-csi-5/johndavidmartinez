// most common strategy is to propagate an error

func main() {
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
}

// chain errors together
doc, err := html.Parse(resp.Body)
resp.Body.Close()
if err != nil {
    return nil, fmt.Errorf("parsing %s as HTML %v", url, error)
}
// so now if it is to be bubbled up adding another error will give
// An easy to understand "trace" of what failed
// Failed to fetch webpage: parsing web as HTML failed

