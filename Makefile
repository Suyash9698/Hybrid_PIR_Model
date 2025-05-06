# Makefile for running PIR demo without Prometheus

N := 6
PORT0 := 8000

# Step 1: Setup data and metadata
setup:
	@echo "📦 Populating files and meta.db..."
	@bash scripts/populate_data.sh

# Step 2: Start all 6 servers in background
servers:
	@echo "🚀 Starting $(N) servers on ports $(PORT0) to $(shell echo $(PORT0)+$(N)-1 | bc)..."
	for i in $(shell seq 0 $(shell echo $(N)-1 | bc)); do \
		PIR_BASEPORT=$$(($(PORT0)+$$i)) go run ./server & \
	done
	@sleep 2
	@echo "✅ All servers launched."

# Step 3: Run the client for file 3
client:
	@echo "🧠 Running PIR client..."
	go run ./client -file=3

# Step 4: Demo: setup + servers + client
demo: setup servers client

# Step 5: Kill all Go server processes and cleanup
clean:
	@echo "🧹 Cleaning up servers and artifacts..."
	@pkill -f "./server" || true
	@rm -rf recovered.bin
