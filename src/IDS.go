package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

// Node structure to store Wikipedia page information
type Node struct {
	ID       string
	Parent   *Node
	Children []*Node
}

// Cache to store parsed links, speeding up repeated searches
var linkCache = make(map[string][]*Node)

// Function to fetch links from a given Wikipedia page using Colly
func get_links(input string, parent *Node) []*Node {
	if nodes, ok := linkCache[input]; ok {
		return nodes
	}

	var nodes []*Node
	collector := colly.NewCollector()

	// Filtering links based on criteria
	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if strings.HasPrefix(href, "/wiki/") && !strings.Contains(href, ":") &&
			!strings.Contains(href, "Main_Page") && !strings.Contains(href, "UTF-8") {
			nodes = append(nodes, &Node{ID: href, Parent: parent})
		}
	})

	// Start scraping
	err := collector.Visit(fmt.Sprintf("https://en.wikipedia.org/wiki/%s", input))
	if err != nil {
		log.Fatal("Failed to scrape the page: ", err)
	}

	linkCache[input] = nodes
	return nodes
}

// Iterative Deepening Search (IDS)
func IDS(root *Node, goal string) *Node {
	depth := 0
	for {
		found := DLS(root, goal, depth)
		if found != nil {
			return found
		}
		depth++
	}
}

// Depth-Limited Search (DLS)
func DLS(node *Node, goal string, depth int) *Node {
	if node.ID == goal {
		return node
	}
	if depth > 0 {
		node.Children = get_links(strings.TrimPrefix(node.ID, "/wiki/"), node)
		for _, child := range node.Children {
			found := DLS(child, goal, depth-1)
			if found != nil {
				return found
			}
		}
	}
	return nil
}

// Main function to set up HTTP server and handle requests
func main() {
	http.HandleFunc("/api/ids", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Println("Received a request")

		startPage := r.URL.Query().Get("start")
		endPage := r.URL.Query().Get("end")

		fmt.Println("Connecting to Wikipedia...")
		startTime := time.Now()
		root := &Node{ID: "/wiki/" + startPage, Children: get_links(startPage, nil)}
		result := IDS(root, "/wiki/"+endPage)
		elapsedTime := time.Since(startTime)

		if result != nil {
			path := make([]string, 0)
			for node := result; node != nil; node = node.Parent {
				path = append([]string{node.ID}, path...)
			}

			response := map[string]interface{}{
				"found":         true,
				"node":          result.ID,
				"path":          path,
				"executionTime": elapsedTime.String(),
			}

			json.NewEncoder(w).Encode(response)
		} else {
			response := map[string]interface{}{
				"found":         false,
				"executionTime": elapsedTime.String(),
			}

			json.NewEncoder(w).Encode(response)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
