package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

type Link struct {
	URL string
}

func fetchLinks(pageURL string) ([]Link, error) {
	c := colly.NewCollector(
		colly.AllowedDomains("en.wikipedia.org"),
	)

	var links []Link

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if strings.HasPrefix(href, "/wiki/") {
			link := Link{
				URL: "https://en.wikipedia.org" + href,
			}
			links = append(links, link)
		}
	})

	err := c.Visit(pageURL)
	if err != nil {
		return nil, err
	}

	return links, nil
}

func BFS(startURL, endURL string) []Link {
	queue := []Link{{URL: startURL}}
	visited := make(map[string]bool)
	path := make(map[string][]Link)

	for len(queue) > 0 {
		currentLink := queue[0]
		queue = queue[1:]

		if currentLink.URL == endURL {
			return path[currentLink.URL]
		}

		if visited[currentLink.URL] {
			continue
		}

		visited[currentLink.URL] = true

		links, err := fetchLinks(currentLink.URL)
		if err != nil {
			fmt.Println("Error fetching links:", err)
			continue
		}

		for _, link := range links {
			if !visited[link.URL] {
				path[link.URL] = append(path[currentLink.URL], link)
				queue = append(queue, link)

				if link.URL == endURL {
					return path[link.URL]
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
