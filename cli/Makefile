.DEFAULT_GOAL := build

.PHONY: build
build:
	go build -o ./hyaline -ldflags="-X 'main.Version=development'" ./cmd/hyaline.go

.PHONY: db
db:
	sqlc generate

.PHONY: install
install:
	go install -ldflags="-X 'main.Version=development'" ./cmd/hyaline.go

.PHONY: test
test:
	rm -rf .coverdata/unit
	mkdir -p .coverdata/unit
	go test -cover ./cmd/... ./internal/... -args -test.gocoverdir="`pwd`/.coverdata/unit"

.PHONY: e2e
e2e:
	rm -rf ./e2e/_output/*
	rm -rf .coverdata/e2e
	mkdir -p .coverdata/e2e
	go build -o ./hyaline-e2e -ldflags="-X 'main.Version=development'" -cover ./cmd/hyaline.go
	go test ./e2e/...
	go tool covdata percent -i=.coverdata/e2e

.PHONY: coverage
coverage:
	make test
	make e2e
	go tool covdata percent -i=./.coverdata/unit,./.coverdata/e2e
	go tool covdata textfmt -i=./.coverdata/unit,./.coverdata/e2e -o ./.coverdata/profile
	go tool cover -html=./.coverdata/profile

.PHONY: e2e-update
e2e-update:
	rm -rf ./e2e/_output/*
	go build -o ./hyaline-e2e -ldflags="-X 'main.Version=development'" ./cmd/hyaline.go
	go test ./e2e/... -update

.PHONY: release
release:
	./scripts/release.sh

.PHONY: git-hooks
git-hooks:
	cp scripts/hooks/* ../.git/hooks/

.PHONY: release-v1
release-v1:
	./scripts/release-v1.sh
