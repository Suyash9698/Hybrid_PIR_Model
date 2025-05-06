package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"sync"
	"time"

	"csis_project/storage"
)

// FetchResult holds one server’s response
type FetchResult struct {
	Server    int
	BytesRead int
	Fraction  float64 // fraction of file
	Err       error
	Duration  time.Duration
}

// FetchParallel fetches /query on N servers and stops once
// downloaded ≥ (α + (1–α)/r)*OriginalSize bytes.
func FetchParallel(n, basePort, fileID int, timeout time.Duration) ([]FetchResult, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := &http.Client{Timeout: timeout}
	results := make([]FetchResult, n)

	// Load metadata
	meta, err := storage.FetchMeta("./data/meta.db")
	if err != nil {
		return nil, fmt.Errorf("could not load meta: %v", err)
	}

	origSize := meta.OriginalSize
	thresholdFrac := meta.Alpha + (1-meta.Alpha)*float64(meta.R)/float64(meta.R-1)

	thresholdBytes := int(math.Ceil(thresholdFrac * float64(origSize)))

	var (
		totalBytes int
		mu         sync.Mutex
		done       bool
		wg         sync.WaitGroup
	)

	fmt.Printf("Stopping once downloaded ≥ %.2f%% of file (%d bytes)\n",
		thresholdFrac*100, thresholdBytes)

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(idx int) {
			defer wg.Done()

			queryType := "uncoded"
			if idx >= meta.R {
				queryType = "coded"
			}
			reqBody, _ := json.Marshal(map[string]interface{}{
				"type":       queryType,
				"file_index": fileID,
			})
			url := fmt.Sprintf("http://localhost:%d/query", basePort+idx)
			req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")

			start := time.Now()
			resp, err := client.Do(req)
			duration := time.Since(start)
			if err != nil {
				results[idx] = FetchResult{Server: idx + 1, Err: err}
				return
			}
			defer resp.Body.Close()

			data, _ := io.ReadAll(resp.Body)
			bytesRead := len(data)
			frac := float64(bytesRead) / float64(origSize)

			mu.Lock()
			if !done {
				totalBytes += bytesRead
				results[idx] = FetchResult{
					Server:    idx + 1,
					BytesRead: bytesRead,
					Fraction:  frac,
					Duration:  duration,
				}

				fmt.Printf("📦 Server %d returned %d bytes (%.2f%%) in %v — total %d bytes (%.2f%%)\n",
					idx+1,
					bytesRead, frac*100,
					duration,
					totalBytes, float64(totalBytes)/float64(origSize)*100,
				)

				if totalBytes >= thresholdBytes {
					done = true
					cancel()
					fmt.Printf("✅ Threshold reached: downloaded %.2f%% of file, stopping.\n",
						float64(totalBytes)/float64(origSize)*100)
				}
			}
			mu.Unlock()
		}(i)
	}

	wg.Wait()
	return results, nil
}
