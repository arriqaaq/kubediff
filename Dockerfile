FROM golang:1.16 as builder

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY main.go main.go
COPY pkg/ pkg/
COPY config/ config/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o kubediff main.go



FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /app/kubediff .
USER nonroot:nonroot

ENTRYPOINT ["/kubediff"]
