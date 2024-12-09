.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

.PHONY: build
build:
	go build -o build/web-scraper main.go

.PHONY: clean
clean:
	rm -rf build

