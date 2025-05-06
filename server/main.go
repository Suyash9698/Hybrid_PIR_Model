package main

import (
	"fmt"
	"log"
	"net/http"

	"csis_project/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/pir", PIRHandler(cfg))
	http.HandleFunc("/query", HybridQueryHandler(cfg))
	addr := fmt.Sprintf(":%d", cfg.BasePort)
	log.Printf("Server listening on HTTP %s", addr)
	http.HandleFunc("/alg", AlgebraicHandler(cfg))
	log.Fatal(http.ListenAndServe(addr, nil))
}
