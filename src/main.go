package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"jojika.go/bfs"
	"jojika.go/ids"
)

type Response struct {
	Path    []string `json:"path"`
	Time    int64    `json:"time"`
	Visited int      `json:"visited"`
}

func main() {
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		method := r.URL.Query().Get("method")
		startPage := r.URL.Query().Get("start")
		endPage := r.URL.Query().Get("end")

		switch method {
		case "bfs":
			// Call BFS function from bfs package
			startURL := "https://en.wikipedia.org/wiki/" + startPage
			endURL := "https://en.wikipedia.org/wiki/" + endPage
			startTime := time.Now()
			shortest, visited := bfs.BFS(startURL, endURL)
			elapsedTime := time.Since(startTime)

			shortestPath := make([]string, len(shortest))
			for i, link := range shortest {
				shortestPath[i] = link.URL // Replace this with the correct field name
			}

			// Write shortest path to response
			json.NewEncoder(w).Encode(Response{Path: shortestPath, Time: elapsedTime.Milliseconds(), Visited: visited})
		case "ids":
			// Call IDS function from ids package
			root := &ids.Node{ID: "/wiki/" + startPage, Children: ids.GetLinks(startPage, nil)} // Changed to GetLinks
			startTime := time.Now()
			result, visited := ids.IDS(root, "/wiki/"+endPage)
			elapsedTime := time.Since(startTime)

			path := make([]string, 0)
			for node := result; node != nil; node = node.Parent {
				path = append([]string{node.ID}, path...)
			}

			// Write result path to response
			json.NewEncoder(w).Encode(Response{Path: path, Time: elapsedTime.Milliseconds(), Visited: visited}) // Changed to use Response struct
		default:
			// Handle unknown method
			http.Error(w, "Unknown method: "+method, http.StatusBadRequest)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
