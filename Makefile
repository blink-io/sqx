
.PHONY: build
build:
	go build -v ./...

.PHONY: upgrade
upgrade:
	go get -u -v ./...