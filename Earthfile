VERSION 0.6

build:
	FROM golang:1.19
	WORKDIR /build

	COPY go.mod go.sum ./
	RUN go mod download

	COPY . .
	ENV CGO_ENABLED=0
	ENV GOOS=linux
	ENV GOARCH=amd64
	RUN go build -o app ./cmd/server

	SAVE ARTIFACT app AS LOCAL tmp/app

test:
	FROM +build
	RUN go test ./...

lint:
	FROM +build
	RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.50.1
	RUN golangci-lint run

image:
	ARG EARTHLY_TARGET_TAG

	FROM busybox
	WORKDIR /build

	COPY +build/app .
	ENTRYPOINT ["/build/app"]

	SAVE IMAGE --push mplewis/gemocities:latest
	SAVE IMAGE --push mplewis/gemocities:$EARTHLY_TARGET_TAG
