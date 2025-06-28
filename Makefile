build:
	@echo "Building WASM filter..."
	cd wasm-filter && \
	go mod tidy && \
	env GOOS=wasip1 GOARCH=wasm go build -buildmode=c-shared -o ../filter.wasm main.go
	@echo "Build complete!"

# Start services
run: build
	docker-compose up --build

# Run tests
test:
	@echo "=== Health check (no auth required) ==="
	curl -i http://localhost:10000/health
	@echo "\n\n=== Request without auth (should fail) ==="
	curl -i http://localhost:10000/api/data
	@echo "\n\n=== Request with invalid token (should fail) ==="
	curl -i -H "Authorization: Bearer invalid-token" http://localhost:10000/api/data
	@echo "\n\n=== Request with valid user token (should succeed) ==="
	curl -i -H "Authorization: Bearer secret-token-123" http://localhost:10000/api/data
	@echo "\n\n=== Request with admin token (should succeed) ==="
	curl -i -H "Authorization: Bearer admin-token-456" http://localhost:10000/api/data
	@echo "\n\n=== Envoy admin stats ==="
	curl -s http://localhost:9901/stats/prometheus | grep wasm

# Cleanup
clean:
	docker-compose down
	rm -f filter.wasm
