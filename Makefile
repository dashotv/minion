include .env
export $(shell sed 's/=.*//' .env)

NAME := runic

all: test

clean:
	rm -rf static/assets/*

test: generate
	go test ./...

generate:
	golem generate

build: generate
	go build

install: build
	go install

server:
	go run main.go server

docker:
	docker build --progress=plain -t $(NAME) .

docker-run:
	docker run --rm --name $(NAME)-test -p 1$(PORT):1$(PORT) \
	-e DOTENV_KEY="dotenv://:key_35ced16628fdf84621dfe41854d304f6ea44f6d12958a51920e71ecb5a4d2ce8@dotenv.local/vault/.env.vault?environment=development" \
	$(NAME)

dotenv:
	npx dotenv-vault local build

ui: clean
	cd ui && yarn build

deps:
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/dashotv/golem@latest

.PHONY: server receiver test deps docker docker-run ui
