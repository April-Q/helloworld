# Build binary
FROM golang:1.14
WORKDIR /workspace

COPY . .
RUN go build -o hello main.go
ENTRYPOINT ["/hello"]
