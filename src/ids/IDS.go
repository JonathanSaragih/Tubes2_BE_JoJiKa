package ids

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

// Struktur Node untuk menyimpan informasi halaman Wikipedia
type Node struct {
	ID       string  // ID unik untuk setiap node yang biasanya adalah URL halaman
	Parent   *Node   // Referensi ke parent node
	Children []*Node // Slice untuk menyimpan referensi ke child node
}

// Cache untuk menyimpan link yang sudah diparsing, mempercepat pencarian yang berulang
var linkCache = make(map[string][]*Node)

// Fungsi GetLinks untuk mengambil link dari halaman Wikipedia menggunakan Colly
func GetLinks(input string, parent *Node) []*Node {
	if nodes, ok := linkCache[input]; ok { // Cek jika link sudah ada di cache
		return nodes
	}

	var nodes []*Node
	collector := colly.NewCollector() // Membuat instance collector baru dari Colly

	// Mengatur filter untuk link yang ditemukan di halaman
	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		// Validasi link yang sesuai dengan kriteria tertentu
		if strings.HasPrefix(href, "/wiki/") && !strings.Contains(href, ":") &&
			!strings.Contains(href, "Main_Page") && !strings.Contains(href, "UTF-8") {
			nodes = append(nodes, &Node{ID: href, Parent: parent})
		}
	})

	// Memulai scraping
	err := collector.Visit(fmt.Sprintf("https://en.wikipedia.org/wiki/%s", input))
	if err != nil {
		log.Fatal("Gagal mengambil halaman: ", err)
	}

	linkCache[input] = nodes // Menyimpan hasil ke cache
	return nodes
}

// (Iterative Deepening Search - IDS)
func IDS(root *Node, goal string) (*Node, int) {
	depth := 0
	visited := 0
	for {
		found := DLS(root, goal, depth, &visited) // Mencoba mencari dengan kedalaman tertentu
		if found != nil {
			return found, visited // Mengembalikan nilai ditemukan dan jumlah halaman yang dikunjungi
		}
		depth++ // Meningkatkan kedalaman dan mencoba lagi
	}
}

// (Depth-Limited Search - DLS)
func DLS(node *Node, goal string, depth int, visited *int) *Node {
	if node.ID == goal {
		return node // Jika ID node sama dengan tujuan, kembalikan node ini
	}
	if depth > 0 {
		node.Children = GetLinks(strings.TrimPrefix(node.ID, "/wiki/"), node) // Mengambil link anak jika belum mencapai batas kedalaman
		*visited += len(node.Children)                                        // Menambahkan jumlah halaman yang dikunjungi
		for _, child := range node.Children {
			// print the current path (parent -> child)
			fmt.Println(node.ID, " -> ", child.ID)
			found := DLS(child, goal, depth-1, visited) // Rekursif mencari pada anak dengan kedalaman dikurangi satu
			if found != nil {
				return found
			}
		}
	}
	return nil // Kembalikan nil jika tidak ditemukan
}

// func DriverIDS() {
// 	http.HandleFunc("/api/ids", func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Access-Control-Allow-Origin", "*")
// 		fmt.Println("Received a request")

// 		startPage := r.URL.Query().Get("start")
// 		endPage := r.URL.Query().Get("end")

// 		fmt.Println("Connecting to Wikipedia...")
// 		startTime := time.Now()
// 		root := &Node{ID: "/wiki/" + startPage, Children: get_links(startPage, nil)}
// 		result := IDS(root, "/wiki/"+endPage)
// 		elapsedTime := time.Since(startTime)

// 		if result != nil {
// 			path := make([]string, 0)
// 			for node := result; node != nil; node = node.Parent {
// 				path = append([]string{node.ID}, path...)
// 			}

// 			response := map[string]interface{}{
// 				"found":         true,
// 				"node":          result.ID,
// 				"path":          path,
// 				"executionTime": elapsedTime.String(),
// 			}

// 			json.NewEncoder(w).Encode(response)
// 		} else {
// 			response := map[string]interface{}{
// 				"found":         false,
// 				"executionTime": elapsedTime.String(),
// 			}

// 			json.NewEncoder(w).Encode(response)
// 		}
// 	})

// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }
