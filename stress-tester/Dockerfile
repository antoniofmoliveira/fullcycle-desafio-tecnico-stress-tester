FROM golang:1.24-bookworm AS build
RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive \
    apt-get install --no-install-recommends --assume-yes \
      build-essential \
      libsqlite3-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=1 CGOOS=linux GOARCH=amd64 go build -o stresstester cmd/main.go

FROM debian:bookworm
RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive \
    apt-get install --no-install-recommends --assume-yes \
      libsqlite3-0
COPY --from=build /app/stresstester /usr/bin/stresstester
ENTRYPOINT ["stresstester"]
CMD ["--url=http://localhost:8080", "--requests=100", "--concurrency=10"]