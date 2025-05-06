package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"csis_project/config"
	"csis_project/pir"
)

type algQuery struct {
	Vector []byte `json:"vector"`
	FileID int    `json:"file_id"`
}

func AlgebraicHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var q algQuery
		if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
			http.Error(w, "bad json", 400)
			return
		}

		sid := cfg.BasePort - 8000 // 0..N-1
		shard := fmt.Sprintf("%s/coded/file%d.shard.%d",
			cfg.DataDir, q.FileID, sid)
		fmt.Printf("🧮 Server on port %d uses sid=%d (file %d)\n", cfg.BasePort, sid, q.FileID)
		row, err := ioutil.ReadFile(shard)
		if err != nil {
			http.Error(w, "shard missing", 404)
			return
		}
		if len(row) < len(q.Vector) {
			http.Error(w, "vector too long", 400)
			return
		}
		resp := pir.Dot(q.Vector, row[:len(q.Vector)])
		w.Write([]byte{resp})
		fmt.Printf("🧩 row fragment for file %d: %v\n", q.FileID, row[:len(q.Vector)])
	}
}
