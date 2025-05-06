package main

import (
	"csis_project/mds"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// EncodeMain is the CLI wrapper to encode a file into RS shards.
func EncodeMain() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: encode <inputFile> <N>")
		os.Exit(1)
	}
	cmd := os.Args[1]
	inFile := os.Args[2]
	N, _ := strconv.Atoi(os.Args[3]) // total shards

	if cmd != "encode" {
		fmt.Println("Only 'encode' supported")
		os.Exit(1)
	}

	data, _ := ioutil.ReadFile(inFile)
	k := N - 1
	m := 1

	shards, err := mds.EncodeBlock(data, k, m)
	if err != nil {
		panic(err)
	}

	base := filepath.Base(inFile)
	prefix := strings.TrimSuffix(base, filepath.Ext(base))
	outDir := filepath.Dir(inFile)
	for i, s := range shards {
		out := filepath.Join(outDir, fmt.Sprintf("%s.shard.%d", prefix, i))
		ioutil.WriteFile(out, s, 0644)
	}
	fmt.Printf("✅ Wrote %d shards to %s\n", len(shards), outDir)
}
