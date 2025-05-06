package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	config "csis_project/config"
	storage "csis_project/storage"
)

// Extracts server ID from environment or config
func getServerID(cfg *config.Config) int {
	return cfg.ServerID // Add this field in config if not already there
}

func HybridQueryHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serverID := getServerID(cfg)

		var req map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		qtype := req["type"].(string)
		fileID := int(req["file_index"].(float64))

		fname := fmt.Sprintf("file%d.bin", fileID)

		if qtype == "coded" {
			// For now, skip matrix logic — just return shard
			path := fmt.Sprintf("data/coded/%s.shard.%d", fname, serverID)
			data, err := ioutil.ReadFile(path)
			if err != nil {
				http.Error(w, "Coded file missing", 404)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(data)

		} else if qtype == "uncoded" {
			path := fmt.Sprintf("data/uncoded/%s.replica.%d", fname, serverID)
			data, err := ioutil.ReadFile(path)
			if err != nil {
				http.Error(w, "Uncoded file missing", 404)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(data)
		} else {
			http.Error(w, "Unknown query type", 400)
		}
	}
}
func PIRHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// load metadata
		meta, err := storage.FetchMeta(cfg.DBPath)
		if err != nil {
			http.Error(w, "meta err", 500)
			return
		}
		// respond with JSON fraction
		data, _ := storage.JSONFraction(meta)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}
