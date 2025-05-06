#!/usr/bin/env bash
set -e

########################################
# 1)  Generate 10 raw 1‑KB files
########################################
mkdir -p data/raw
for i in $(seq 1 10); do
  dd if=/dev/urandom of=data/raw/file${i}.bin bs=1024 count=1
done
echo "✅ 10 raw files written to data/raw/"

########################################
# 2)  Run hybrid placement (your file_manager)
########################################
./bin/file_manager --datadir=./data --k=4 --m=2
echo "✅ Hybrid placement finished"

########################################
# 3)  Build tiny BU‑coded bundles
#     – file1  (demo from earlier)
#     – file3  (needed for full BU test)
########################################
mkdir -p data/coded

# ---- file1 ----
dd if=/dev/urandom bs=32 count=1 of=data/coded/file1.orig 2>/dev/null
go run cmd/mds_encode/*.go encode data/coded/file1.orig 6
# ---- file3 ----
dd if=/dev/urandom bs=32 count=1 of=data/coded/file3.orig 2>/dev/null
go run cmd/mds_encode/*.go encode data/coded/file3.orig 6
echo "✅ BU demo shards written for file1 and file3"
