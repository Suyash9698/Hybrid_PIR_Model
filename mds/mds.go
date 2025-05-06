package mds

import (
	"bytes"

	"github.com/klauspost/reedsolomon"
)

// EncodeBlock splits data into k data shards + m parity shards.
func EncodeBlock(data []byte, k, m int) ([][]byte, error) {
	enc, err := reedsolomon.New(k, m)
	if err != nil {
		return nil, err
	}
	shards, err := enc.Split(data)
	if err != nil {
		return nil, err
	}
	if err := enc.Encode(shards); err != nil {
		return nil, err
	}
	return shards, nil
}

// DecodeBlock reconstructs original data from any k shards.
func DecodeBlock(shards [][]byte, k int) ([]byte, error) {
	enc, err := reedsolomon.New(k, len(shards)-k)
	if err != nil {
		return nil, err
	}
	ok, err := enc.Verify(shards)
	if err != nil || !ok {
		if err := enc.Reconstruct(shards); err != nil {
			return nil, err
		}
	}
	buf := &bytes.Buffer{}
	if err := enc.Join(buf, shards, len(shards[0])*k); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
