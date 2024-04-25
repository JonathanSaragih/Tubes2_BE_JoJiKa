package main

import (
    "fmt"
    "log"
    "strings"
    "time"

    "net/http"
    "golang.org/x/net/html"
    "container/list"
)


type Link struct {
    URL string
}

func fetchLinks(pageURL string) []Link {
    response, err := http.Get(pageURL)
    if err != nil {
        fmt.Println("Error:", err)
        return nil
    }
    defer response.Body.Close()

    document, err := html.Parse(response.Body)
    if err != nil {
        fmt.Println("Error:", err)
        return nil
    }

    var links []Link
    var explore func(*html.Node)
    explore = func(n *html.Node) {
        if n.Type == html.ElementNode && n.Data == "a" {
            for _, attr := range n.Attr {
                if attr.Key == "href" && strings.HasPrefix(attr.Val, "/wiki/") {
                    validLink := true
                    for _, class := range strings.Fields(attr.Val) {
                        if class == "new" || strings.Contains(strings.ToLower(class), "portal") {
                            validLink = false
                            break
                        }
                    }
                    if validLink && strings.HasPrefix(attr.Val, "/wiki/") && !strings.Contains(attr.Val, ":") {
                        link := Link{
                            URL:   "https://en.wikipedia.org" + attr.Val,
                        }
                        links = append(links, link)
                    }
                }
            }
        }
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            explore(c)
        }
    }
    explore(document)

    return links
}


func BFS(startURL, endURL string) []Link {
    queue := list.New()
    visited := make(map[string]bool)
    path := make(map[string][]Link)

    queue.PushBack([]Link{{URL: startURL}})

    for queue.Len() > 0 {
        currentPath := queue.Remove(queue.Front()).([]Link)
        currentLink := currentPath[len(currentPath)-1]

        if currentLink.URL == endURL {
            return currentPath
        }

        links := fetchLinks(currentLink.URL)
        for _, link := range links {
            if !visited[link.URL] {
                visited[link.URL] = true
                newPath := append(currentPath, link)
                queue.PushBack(newPath)
                path[link.URL] = newPath

                if link.URL == endURL {
                    return newPath
                }
            }
        }
    }

    return nil
}

func main() {
    var startPage, endPage string

    fmt.Print("Masukkan judul halaman awal: ")
    fmt.Scanln(&startPage)

    fmt.Print("Masukkan judul halaman akhir: ")
    fmt.Scanln(&endPage)

    startURL := "https://en.wikipedia.org/wiki/" + startPage
    endURL := "https://en.wikipedia.org/wiki/" + endPage

    startTime := time.Now()

    shortest := BFS(startURL, endURL)
    if shortest == nil {
        log.Fatal("Tidak ditemukan jalur")
    }

    fmt.Print("Jalur terpendek: ")
    for _, link := range shortest {
        fmt.Print(link.URL)
        if link.URL != endURL {
            fmt.Print(" > ")
        } else {
            fmt.Println()
        }
    }
    endTime := time.Now()
    elapsed := endTime.Sub(startTime)
    fmt.Println("Waktu eksekusi:", elapsed)
}