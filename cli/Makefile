.DEFAULT_GOAL := build

.PHONY: build
build:
	go build -o ./hyaline -ldflags="-X 'main.Version=development'" ./cmd/hyaline.go

.PHONY: test
test:
	go test ./cmd/... ./internal/...

.PHONY: e2e
e2e:
	rm -rf ./e2e/_output/*
	rm -rf .coverdata
	mkdir -p .coverdata
	go build -o ./hyaline-e2e -ldflags="-X 'main.Version=development'" -cover ./cmd/hyaline.go
	go test ./e2e/...
	go tool covdata percent -i=.coverdata

.PHONY: e2e-update
e2e-update:
	rm -rf ./e2e/_output/*
	go build -o ./hyaline-e2e -ldflags="-X 'main.Version=development'" ./cmd/hyaline.go
	go test ./e2e/... -update

.PHONY: extract-current
extract-current:
	rm -f ./current.db
	make build
	echo "---"
	./hyaline --debug extract current --config ./_example/config.yml --system local --output ./current.db

.PHONY: extract-current-git
extract-current-git:
	rm -f ./current.db
	make build
	echo "---"
	./hyaline --debug extract current --config ./_example/config.yml --system git --output ./current.db

.PHONY: extract-current-http
extract-current-http:
	rm -f ./current.db
	make build
	echo "---"
	./hyaline --debug extract current --config ./_example/config.yml --system http --output ./current.db

.PHONY: extract-change
extract-change:
	rm -f ./change.db
	make build
	echo "---"
	./hyaline --debug extract change --config ./_example/config.yml --system git --base main --head origin/feat-1 --pull-request appgardenstudios/hyaline-example/1 --issue appgardenstudios/hyaline-example/2 --issue appgardenstudios/hyaline-example/3  --output ./change.db

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

.PHONY: merge
merge:
	rm -f ./merged.db
	make build
	echo "---"
	./hyaline --debug merge --input ./current.db --output ./merged.db