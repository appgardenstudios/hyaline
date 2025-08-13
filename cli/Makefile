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

.PHONY: benchmark-prompts-refactor
benchmark-prompts-refactor:
	make build
	cd benchmarks/prompts && go test -v -run TestCheckChangeRefactor

.PHONY: benchmark-prompts-remove-feature
benchmark-prompts-remove-feature:
	make build
	cd benchmarks/prompts && go test -v -run TestCheckChangeRemoveFeature

.PHONY: benchmark-prompts-api-contract
benchmark-prompts-api-contract:
	make build
	cd benchmarks/prompts && go test -v -run TestCheckChangeAPIContract

.PHONY: benchmark-prompts-add-feature
benchmark-prompts-add-feature:
	make build
	cd benchmarks/prompts && go test -v -run TestCheckChangeAddFeature

.PHONY: benchmark-prompts-check-current-api
benchmark-prompts-check-current-api:
	make build
	cd benchmarks/prompts && go test -v -run TestCheckCurrentAPIDocs

.PHONY: benchmark-prompts-check-current-customer
benchmark-prompts-check-current-customer:
	make build
	cd benchmarks/prompts && go test -v -run TestCheckCurrentCustomerDocs

.PHONY: benchmark-prompts
benchmark-prompts:
	make build
	cd benchmarks/prompts && go test -v

.PHONY: release
release:
	./scripts/release.sh

.PHONY: git-hooks
git-hooks:
	cp scripts/hooks/* ../.git/hooks/

.PHONY: release-v1
release-v1:
	./scripts/release-v1.sh
