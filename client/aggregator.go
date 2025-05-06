package main

import (
	storage "csis_project/storage"
	"fmt"
	"math"
)

// ComputeCosts prints theoretical vs. measured cost
func ComputeCosts(results []FetchResult, meta *storage.Meta) {
	// theoretical
	codedCost := meta.Alpha
	uncodedCost := (1 - meta.Alpha) * float64(meta.R) / float64(meta.R-1)
	theory := codedCost + uncodedCost

	// measured sum
	var sum float64
	for _, r := range results {
		if r.Err == nil {
			sum += r.Fraction
		}
	}
	fmt.Println("Download Costs:")
	fmt.Printf(" • codedCost = %.4f\n", codedCost)
	fmt.Printf(" • uncodedCost = %.4f\n", uncodedCost)
	fmt.Printf(" ✅ Theoretical total = %.4f\n", theory)
	fmt.Printf(" ✅ Measured total   = %.4f\n", sum)
	fmt.Printf("Difference = %.4f (%.2f%%)\n",
		math.Abs(theory-sum), math.Abs(theory-sum)/theory*100)
}
