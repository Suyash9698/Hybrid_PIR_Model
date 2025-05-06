package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"csis_project/config"
	"csis_project/pir"
)

func runAlgebraicDemo(cfg *config.Config, fileID, L, theta int) {
	vec := pir.BuildQuery(L, theta)
	cli := &http.Client{Timeout: 4 * time.Second}
	resps := make([]byte, cfg.N)

	for i := 0; i < cfg.N; i++ {
		url := fmt.Sprintf("http://localhost:%d/alg", cfg.BasePort+i)
		payload, _ := json.Marshal(map[string]interface{}{
			"vector": vec, "file_id": fileID,
		})
		resp, err := cli.Post(url, "application/json",
			bytes.NewReader(payload))
		if err != nil {
			panic(err)
		}
		b, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		resps[i] = b[0]
		fmt.Printf("Srv%02d → 0x%02x\n", i+1, b[0])
	}
	sym := pir.DecodeSymbol(resps)
	fmt.Printf("🔒 decoded coded symbol θ=%d = 0x%02x\n", theta, sym)
}
