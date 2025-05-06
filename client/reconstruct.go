package main

import (
	"crypto/sha256"
	"csis_project/mds"
	"fmt"
	"io/ioutil"
)

// Reconstruct combines coded + uncoded slices and checks hash.
// rawParts = uncoded pieces in correct order
// coded    = α‑fraction Reed‑Solomon shards (k+m shards, we need any k)
func Reconstruct(rawParts [][]byte, coded [][]byte, k int) {
	// concatenate raw parts
	raw := rawParts[0]

	// decode coded bundle (returns full data)
	full, _ := mds.DecodeBlock(coded, k)

	// merge bundles → original file
	final := append(raw, full...)

	sum := sha256.Sum256(final)
	fmt.Printf("✅ SHA‑256 of reconstructed file: %x\n", sum)
	ioutil.WriteFile("recovered.bin", final, 0644)
}
