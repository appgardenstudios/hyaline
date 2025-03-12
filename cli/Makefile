.DEFAULT_GOAL := build

.PHONY: build
build:
	go build -o ./hyaline -ldflags="-X 'main.Version=development'" ./cmd/hyaline.go

.PHONY: current
current:
	rm -f ./current.db
	make build
	echo "---"
	./hyaline --debug extract current --config ./_example/config.yml --system local --output ./current.db

.PHONY: change
change:
	rm -f ./change.db
	make build
	echo "---"
	./hyaline --debug extract change --config ./_example/config.yml --system local --base main --head feat-1 --output ./change.db

.PHONY: check
check:
	make build
	echo "---"
	./hyaline --debug check --config ./_example/config.yml --current ./current.db --system local

.PHONY: recommend
recommend:
	make build
	echo "---"
	./hyaline --debug check --config ./_example/config.yml --current ./current.db --system local --recommend