BIN="./github-data-analyzer"

ifeq (, $(shell which golangci-lint))
$(warning "could not find golangci-lint in $(PATH)")
$(info "run: curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh")
endif

.PHONY: lint test install_deps clean

build: install_deps
	$(info ******************** building github-data-analyzer ********************)
	go build -o github-data-analyzer .

lint:
	$(info ******************** running lint tools ********************)
	golangci-lint run -v

test: install_deps lint
	$(info ******************** running tests ********************)
	go test -v -cover ./...

install_deps:
	$(info ******************** downloading dependencies ********************)
	go mod vendor -v

docker-build:
	$(info ******************** building docker image ********************)
	docker build -t github-data-analyzer:latest .

clean:
	rm -rf $(BIN)