.DEFAULT_GOAL := build

.PHONY: build
build:
	go build -o ./hyaline -ldflags="-X 'main.Version=development'" ./cmd/hyaline.go

.PHONY: extract
extract:
	rm -f ./current.db
	make build
	echo "---"
	./hyaline extract current --config ./_example/config.yml --system local --output ./current.db