.DEFAULT_GOAL := build

.PHONY: build
build:
	go build -o ./hyaline -ldflags="-X 'main.Version=development'" ./cmd/hyaline.go

.PHONY: install
install:
	go install -ldflags="-X 'main.Version=development'" ./cmd/hyaline.go

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

.PHONY: pre-check-current
pre-check-current:
	rm -f ./current.db
	make build
	./hyaline --debug extract current --config ./_example/config.yml --system check-current --output ./current.db

.PHONY: check-current
check-current:
	rm -f ./results.json
	make build
	echo "---"
	./hyaline --debug check current --config ./_example/config.yml --current ./current.db --system check-current --output ./results.json

.PHONY: check-current-purpose
check-current-purpose:
	rm -f ./results.json
	make build
	echo "---"
	./hyaline --debug check current --config ./_example/config.yml --current ./current.db --system check-current --output ./results.json --check-purpose

.PHONY: check-current-completeness
check-current-completeness:
	rm -f ./results.json
	make build
	echo "---"
	./hyaline --debug check current --config ./_example/config.yml --current ./current.db --system check-current --output ./results.json --check-completeness

.PHONY: extract-change
extract-change:
	rm -f ./change.db
	make build
	echo "---"
	./hyaline --debug extract change --config ./_example/config.yml --system git --base main --head origin/feat-1 --pull-request appgardenstudios/hyaline-example/1 --issue appgardenstudios/hyaline-example/2 --issue appgardenstudios/hyaline-example/3  --output ./change.db

.PHONY: pre-check-change
pre-check-change:
	rm -f ./current.db ./change.db
	make build
	./hyaline --debug extract current --config ./_example/config.yml --system check-change --output ./current.db
	./hyaline --debug extract change --config ./_example/config.yml --system check-change --base main --head origin/feat-1 --pull-request appgardenstudios/hyaline-example/1 --issue appgardenstudios/hyaline-example/2 --issue appgardenstudios/hyaline-example/3  --output ./change.db

.PHONY: check-change
check-change:
	rm -f ./results.json
	make build
	echo "---"
	./hyaline --debug check change --config ./_example/config.yml --current ./current.db --change ./change.db --system check-change --output ./results.json

.PHONY: check-change-suggest
check-change-suggest:
	rm -f ./results.json
	make build
	echo "---"
	./hyaline --debug check change --config ./_example/config.yml --current ./current.db --change ./change.db --system check-change --output ./results.json --suggest

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

.PHONY: generate-config
generate-config:
	rm -f ./current.db ./generated-config.yml
	make build
	./hyaline --debug extract current --config ./_example/config.yml --system generate-config --output ./current.db
	echo "---"
	./hyaline --debug generate config --config ./_example/config.yml --current ./current.db --system generate-config --output ./generated-config.yml

.PHONY: generate-config-with-purpose
generate-config-with-purpose:
	rm -f ./current.db ./generated-config.yml
	make build
	./hyaline --debug extract current --config ./_example/config.yml --system generate-config --output ./current.db
	echo "---"
	./hyaline --debug generate config --config ./_example/config.yml --current ./current.db --system generate-config --output ./generated-config.yml --include-purpose

.PHONY: merge
merge:
	rm -f ./merged.db
	make build
	echo "---"
	./hyaline --debug merge --input ./current.db --input ./current-copy.db --input ./change.db --output ./merged.db

.PHONY: pre-update-pr
pre-update-pr:
	make pre-check-change
	make check-change

.PHONY: update-pr
update-pr:
	rm -f ./comment.json
	make build
	echo "---"
	./hyaline --debug update pr --config ./_example/config.yml --pull-request appgardenstudios/hyaline-example/1 --sha b4c5c73 --recommendations ./results.json --output ./comment.json

.PHONY: update-pr-comment
update-pr-comment:
	rm -f ./comment.json
	make build
	echo "---"
	./hyaline --debug update pr --config ./_example/config.yml --pull-request appgardenstudios/hyaline-example/1 --comment appgardenstudios/hyaline-example/2917391402 --sha b4c5c73 --recommendations ./results.json --output ./comment.json

.PHONY: release
release:
	./scripts/release.sh