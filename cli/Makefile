.DEFAULT_GOAL := build

.PHONY: build
build:
	go build -o ./hyaline -ldflags="-X 'main.Version=development'" ./cmd/hyaline.go

.PHONY: test
test:
	go test ./...

.PHONY: extract-current
extract-current:
	rm -f ./current.db
	make build
	echo "---"
	./hyaline --debug extract current --config ./_example/config.yml --system local --output ./current.db

.PHONY: extract-change
extract-change:
	rm -f ./change.db
	make build
	echo "---"
	./hyaline --debug extract change --config ../../hyaline-example/config.yml --system local --base main --head feat-1 --output ./change.db

.PHONY: check-current
check-current:
	make build
	echo "---"
	./hyaline --debug check --config ./_example/config.yml --current ./current.db --system local

.PHONY: check-change
check-change:
	make build
	echo "---"
	./hyaline --debug check --config ./_example/config.yml --current ./current.db --change ./change.db --system local

.PHONY: recommend-current
recommend-current:
	make build
	echo "---"
	./hyaline --debug check --config ./_example/config.yml --current ./current.db --system local --recommend

.PHONY: recommend-change
recommend-change:
	make build
	echo "---"
	./hyaline --debug check --config ./_example/config.yml --current ./current.db --change ./change.db --system local --recommend