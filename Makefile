.PHONY: build clean deploy run

mod:
	go mod vendor -v

build:
	export GO111MODULE=on
	go build -mod=readonly -o bin/gorest ./cmd/gorest

test:
	go test -mod=readonly -race -cover ./... -v

run:
	ENV_PATH=config/envs/local/config.yml \
	./bin/gorest

init-redis:
	redis-server &

clean:
	rm -rf ./bin
	redis-cli shutdown

auth:
	go build -mod=readonly -o bin/auth ./tools/authgenerator
	./bin/auth

all: build test init-redis run
