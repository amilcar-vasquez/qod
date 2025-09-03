## run/api: run go fmt, then run the cmd/api application

.PHONY: run/api
run/api:
	@echo '-- Running go fmt --'
	@go fmt ./...
	@echo '-- Running application --'
	@go run ./cmd/api