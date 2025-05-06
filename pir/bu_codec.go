package pir

import (
	"crypto/rand"
	"fmt"
)

// BuildQuery returns random vector v (len=L) that hides theta.
func BuildQuery(L, theta int) []byte {
	v := make([]byte, L)
	var sum byte
	for i := 0; i < L; i++ {
		if i == theta {
			continue
		}
		_, _ = rand.Read(v[i : i+1])
		if v[i] == 0 {
			v[i] = 1
		}
		sum = Add(sum, v[i])
	}
	v[theta] = Add(1, sum) // ensures ⟨v,row⟩ = row[theta]
	fmt.Printf("🔍 Built query vector (θ=%d): %v\n", theta, v)
	return v
}

// DecodeSymbol XORs (GF) all server responses → coded symbol θ.
func DecodeSymbol(resps []byte) byte {
	var s byte
	for _, b := range resps {
		s ^= b
	}
	return s
}
