package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"

	mds "csis_project/mds"

	_ "github.com/mattn/go-sqlite3"
)

// ceil
func ceil(x float64) int { return int(math.Ceil(x)) }

// Placement metadata
type Meta struct {
	Alpha        float64
	R            int
	N            int
	OriginalSize int
}

// Initialize storage directory and DB
func InitStorage(dataDir, dbPath string, N int, mu float64) error {
	// create data dir
	if err := os.MkdirAll(dataDir+"/raw", 0755); err != nil {
		return err
	}
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS meta (
            id INTEGER PRIMARY KEY,
            alpha REAL,
            r INTEGER,
            n INTEGER,
			original_size INTEGER
        );
        DELETE FROM meta;
    `)
	return err
}

// Split raw files into coded+uncoded shards & save metadata
func DoPlacement(dataDir string, N int, mu float64, kShard, mShard int) error {
	rawDir := dataDir + "/raw"
	files, err := ioutil.ReadDir(rawDir)
	if err != nil {
		return err
	}
	r := ceil(float64(N)*mu) - 1
	alpha := (mu - float64(r)/float64(N)) / (1 - float64(r)/float64(N))

	// persist meta
	db, err := sql.Open("sqlite3", dataDir+"/meta.db")
	if err != nil {
		return err
	}
	defer db.Close()
	// Read size of first file (e.g., file0.txt)
	var originalSize int
	if len(files) > 0 {
		fpath := filepath.Join(rawDir, files[0].Name())
		data, _ := ioutil.ReadFile(fpath)
		originalSize = len(data)
	}

	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS meta (
        id INTEGER PRIMARY KEY,
        alpha REAL,
        r INTEGER,
        n INTEGER,
        original_size INTEGER
    );
    DELETE FROM meta;
    `)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO meta(alpha,r,n,original_size) VALUES (?, ?, ?, ?)", alpha, r, N, originalSize)

	if err != nil {
		return err
	}

	// for each file
	for _, f := range files {
		filePath := filepath.Join(rawDir, f.Name())
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return err
		}

		// 1) uncoded: replicate full data to r servers
		for i := 0; i < r; i++ {
			out := dataDir + fmt.Sprintf("/uncoded/%s.replica.%d", f.Name(), i)
			os.MkdirAll(filepath.Dir(out), 0755)
			if err := ioutil.WriteFile(out, data, 0644); err != nil {
				return err
			}
		}
		// 2) coded: MDS encode full data
		shards, err := mds.EncodeBlock(data, kShard, mShard)
		if err != nil {
			return err
		}
		for i, shard := range shards {
			out := dataDir + fmt.Sprintf("/coded/%s.shard.%d", f.Name(), i)
			os.MkdirAll(filepath.Dir(out), 0755)
			if err := ioutil.WriteFile(out, shard, 0644); err != nil {
				return err
			}
		}
	}
	return nil
}

// Fetch metadata
func FetchMeta(dbPath string) (*Meta, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	row := db.QueryRow("SELECT alpha,r,n,original_size FROM meta LIMIT 1")

	var m Meta
	if err := row.Scan(&m.Alpha, &m.R, &m.N, &m.OriginalSize); err != nil {
		return nil, err
	}
	return &m, nil
}

// File fraction: α + (1-α)/r
func Fraction(meta *Meta) float64 {
	return meta.Alpha + (1-meta.Alpha)/float64(meta.R)
}

// For JSON output
func JSONFraction(meta *Meta) ([]byte, error) {
	out := map[string]interface{}{
		"fraction": Fraction(meta),
		"alpha":    meta.Alpha,
		"r":        meta.R,
	}
	return json.Marshal(out)
}
