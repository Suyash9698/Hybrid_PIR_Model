package main

import (
	"csis_project/storage"
	"fmt"
	"log"
)

func main() {
	dataDir := "./data"
	N := 6
	mu := 0.5
	k := 4
	m := 2

	if err := storage.InitStorage(dataDir, dataDir+"/meta.db", N, mu); err != nil {
		log.Fatal("Init error:", err)
	}
	if err := storage.DoPlacement(dataDir, N, mu, k, m); err != nil {
		log.Fatal("Placement error:", err)
	}
	fmt.Println("✅ Placement done")
}
