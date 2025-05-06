package main

import (
	"flag"
	"fmt"
	"log"

	"bytes"
	"csis_project/config"
	"csis_project/storage"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
)

func main() {
	alg := flag.Bool("alg", false, "run algebraic BU demo and exit")
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	fileID := flag.Int("file", 1, "file ID to fetch")
	flag.Parse()

	if *alg {
		cfg, _ := config.Load()
		runAlgebraicDemo(cfg, *fileID, 7 /* <= shard size */, 5) // not 32
		return
	}

	// load meta
	meta, err := storage.FetchMeta(cfg.DBPath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Using α=%.4f, r=%d, N=%d\n", meta.Alpha, meta.R, meta.N)

	// fetch in parallel
	results, err := FetchParallel(cfg.N, cfg.BasePort, *fileID, cfg.HTTPTimeout)

	if err != nil {
		log.Fatal(err)
	}

	// aggregate
	for _, r := range results {
		if r.Err != nil {
			fmt.Printf("❌ Server %d error: %v\n", r.Server, r.Err)
		} else {
			fmt.Printf("📦 Server %d returned: %.4f fraction\n", r.Server, r.Fraction)
		}
	}

	// Split into rawParts and codedParts based on server index
	var rawParts [][]byte
	var codedParts [][]byte

	for i, r := range results {
		if r.Err != nil {
			continue
		}
		// Simulate part content (in real use: fetch real content)
		var part []byte
		if r.Server < meta.R {
			// Uncoded request
			req := map[string]interface{}{
				"type":       "uncoded",
				"file_index": *fileID,
			}
			part = SendQuery(cfg.BasePort+r.Server, req)
		} else {
			// Coded request
			req := map[string]interface{}{
				"type":       "coded",
				"file_index": *fileID,
				"query":      GenerateCodedQuery(), // implement this to return []float64
			}
			part = SendQuery(cfg.BasePort+r.Server, req)
		}

		if i < meta.R {
			rawParts = append(rawParts, part)
		} else {
			codedParts = append(codedParts, part)
		}
	}

	Reconstruct(rawParts, codedParts, meta.R)
	ComputeCosts(results, meta)
}

func SendQuery(port int, req map[string]interface{}) []byte {
	url := fmt.Sprintf("http://localhost:%d/query", port)
	body, _ := json.Marshal(req)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("❌ Error sending request to %s: %v", url, err)
	}
	defer resp.Body.Close()
	result, _ := io.ReadAll(resp.Body)
	return result
}

func GenerateCodedQuery() []float64 {
	q := make([]float64, 32) // example length
	for i := range q {
		q[i] = rand.Float64()
	}
	return q
}
