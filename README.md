# 🔐 Hybrid Private Information Retrieval (PIR) System

A research-grade implementation of **Hybrid PIR** combining:
- 📦 **Uncoded replication** (Attia–Kumar–Tandon scheme)
- 🔢 **Coded PIR** using algebraic queries (Banawan–Ulukus model)

It automatically selects the best download strategy depending on the storage fraction **μ** — blending privacy, efficiency, and practicality.

---

## 📁 File Structure
```text
CSIS_PROJECT/
├── bin/                            # Compiled binaries
├── client/                         # Main PIR client logic: Fetch, Reconstruct
│   ├── aggregator.go
│   ├── alg_client.go
│   ├── Dockerfile
│   ├── fetcher.go
│   ├── main.go
│   ├── reconstruct.go
├── cmd/
│   ├── file_manager/              # Placement tool (replica + coded placement)
│   │   ├── main.go
│   │   └── mds_encode/
│   │       ├── encode_main.go
│   │       └── main.go
├── config/                         # Config file loader
│   └── config.go
├── data/
│   ├── raw/                        # Input files
│   ├── coded/                      # RS-coded shards
│   ├── uncoded/                    # Replicated full blocks
│   └── meta.db                     # Metadata file storing α, r, N, file size
├── mds/                            # Matrix-based PIR utils (gf256)
│   ├── bu_codec.go
│   ├── gf256.go
│   └── mds.go
├── pir/                            # BU-coded PIR codec functions
│   └── handler_alg.go
├── scripts/
│   ├── cleanup.sh                  # Wipe old data
│   └── populate_data.sh           # Populate fresh test files
├── server/
│   ├── main.go                     # Starts a PIR server instance
│   ├── handler_alg.go             # Algebraic PIR handler (optional)
│   └── handlers.go                # HybridQueryHandler (coded + uncoded)
├── storage/                        # Placement + metadata storage logic
│   ├── storeFile.go
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── Makefile
└── recovered.bin
```


## 🚀 How to Run

### 1️⃣ Populate data
```bash
bash scripts/populate_data.sh
```
This generates:

- 10 files of 1KB each  
- Stores meta info (`α`, `r`, `N`)  
- Splits into coded + uncoded per hybrid PIR strategy  


2️⃣ Start Servers
```bash
for i in {0..5}; do
  PIR_SERVERID=$i PIR_BASEPORT=$((8000+i)) go run server/*.go &
done
```
- Starts 6 independent PIR servers on ports 8000–8005.

3️⃣ Run Client
```bash
go run cmd/client/main.go -file=3
```



## 🧠 How It Works – Detailed Code Walkthrough

This project implements a **Hybrid Private Information Retrieval (PIR)** system combining **replication-based** and **Reed–Solomon-coded** strategies to balance download cost and storage constraints.

---

### 🔧 Code Architecture Breakdown

| Component               | Path                          | Description |
|------------------------|-------------------------------|-------------|
| `client/main.go`       | `cmd/client/`                 | Entry point for the PIR client. Loads metadata, calls `FetchParallel()`, splits results, reconstructs file, computes costs. |
| `FetchParallel()`      | `cmd/client/fetcher.go`       | Sends concurrent requests to PIR servers, stops after threshold bytes are fetched. |
| `storeFile.go`         | `storage/`                    | Handles hybrid placement logic: splits raw file into coded and uncoded shards and stores them. |
| `HybridQueryHandler()` | `server/handlers.go`          | Responds to PIR queries by serving either uncoded replicas or coded shards, depending on request type. |
| `populate_data.sh`     | `scripts/`                    | Bash script to generate 10 files, store meta info, and perform placement. |
| `file_manager`         | `cmd/file_manager`            | Binaries that execute file placement (splitting and encoding). |
| `meta.db`              | `data/`                       | Stores metadata like file size, α (coded fraction), r (replication factor), and total N servers. |
| `ComputeCosts()`       | `cmd/client/main.go`          | Compares theoretical vs measured download costs. |
| `SendQuery()`          | `cmd/client/main.go`          | Utility to make HTTP requests and fetch PIR responses from server. |

---

### 📈 Workflow: End-to-End Execution

#### ✅ Step 1: File Placement (Hybrid Strategy)
- `bash scripts/populate_data.sh`:
  - Generates 10 × 1KB files in `data/raw/`
  - Runs `file_manager` to:
    - Store `(1–α)` fraction using **r replicated** uncoded replicas.
    - Store `α` fraction as **coded shards** using Reed–Solomon (RS) encoding.
  - Metadata like α, r, N, and file size stored in `data/meta.db`.

#### ✅ Step 2: Launch PIR Servers
- Each server knows its `ServerID` and responds only to matching `fileX.replica.Y` or `fileX.shard.Y` from disk.
- `HybridQueryHandler()` serves PIR queries using this rule:
  ```go
  if serverID < r ⇒ serve uncoded replica
  else ⇒ serve coded shard
  ```

#### ✅ Step 3: Client-Side Parallel Fetch
- Client runs:
  ```bash
  go run cmd/client/main.go -file=3
  ```
  - Loads α, r, N, file size from `meta.db`
  - Calculates download **threshold**:
    ```
    threshold = α + (1–α)·r / (r–1)
    ```
  - Sends parallel requests to all N servers.
  - Collects `[]byte` parts.
  - Stops early if threshold met using `context.Cancel()`.

#### ✅ Step 4: File Reconstruction
- The client:
  - Collects valid uncoded and coded parts.
  - Reconstructs file using appropriate decoding method (trivial copy if uncoded or RS decode if coded).
  - Verifies file integrity using `SHA-256`.

#### ✅ Step 5: Reporting Cost
- `ComputeCosts()` compares:
  - **Theoretical Cost**: based on α and r
  - **Measured Cost**: based on total bytes received
  - Output:
    ```text
    ✅ Theoretical total = 1.7500
    ✅ Measured total   = 1.7500
    ```

---

### 📌 PIR Query Format

Client sends:
```json
{
  "type": "coded" | "uncoded",
  "file_index": 3
}
```

Server responds with corresponding block from disk (`.replica.X` or `.shard.X`).

---

### ✅ Summary

This system demonstrates how **storage–download tradeoffs** in PIR can be tuned using a **hybrid coded + uncoded strategy**, offering flexibility for real-world deployment under different network or storage constraints.



## 📤 Sample Output:

```
Using α=0.2500, r=2, N=6
Stopping once downloaded ≥ 175.00% of file (1792 bytes)
📦 Server 4 returned 256 bytes (25.00%) in 5.235791ms — total 256 bytes (25.00%)
📦 Server 1 returned 1024 bytes (100.00%) in 5.401459ms — total 1280 bytes (125.00%)
📦 Server 3 returned 256 bytes (25.00%) in 6.033458ms — total 1536 bytes (150.00%)
📦 Server 6 returned 256 bytes (25.00%) in 6.06425ms — total 1792 bytes (175.00%)
✅ Threshold reached: downloaded 175.00% of file, stopping.
📦 Server 1 returned: 1.0000 fraction
📦 Server 0 returned: 0.0000 fraction
📦 Server 3 returned: 0.2500 fraction
📦 Server 4 returned: 0.2500 fraction
📦 Server 0 returned: 0.0000 fraction
📦 Server 6 returned: 0.2500 fraction
✅ SHA‑256 of reconstructed file: a138f712960033894450eace81cd056e97003f8d36df270ff8b30bf42503a88e
Download Costs:
 • codedCost = 0.2500
 • uncodedCost = 1.5000
 ✅ Theoretical total = 1.7500
 ✅ Measured total   = 1.7500
Difference = 0.0000 (0.00%)
```



## 📊 Theory vs Measured Cost

For given μ and r:

✅ **Our system:**
- Measures actual bytes fetched  
- Verifies file integrity (SHA‑256)  
- Stops as soon as threshold is met  

---

## ⚙️ Internals

| Component           | Role                                      | Location                  |
|---------------------|-------------------------------------------|---------------------------|
| `HybridQueryHandler` | Serves coded/uncoded shards               | `server/handlers.go`      |
| `FetchParallel`     | Makes parallel requests + stops on threshold | `cmd/client/fetcher.go`   |
| `storeFile.go`      | Places replicas and coded shards          | `storage/`                |
| `meta.db`           | Stores α, r, N, file size                 | `data/`                   |
| `ComputeCosts`      | Reports theoretical vs actual download    | `cmd/client/`             |



### 🧪 Sanity Checks
✅ α, r, N are consistent from meta.db

✅ Fraction matches per byte count

✅ Measured == Theoretical within 4 decimals

✅ File is reconstructed exactly (via SHA‑256)


### 🧾 Credits

- Inspired by:

  - 📘 *"The Capacity of Private Information Retrieval from Coded Databases"*  
    *Banawan & Ulukus*

  - 📘 *"The Capacity of PIR with Heterogeneous Replication"*  
    *Attia–Kumar–Tandon*
