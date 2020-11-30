FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -a -installsuffix cgo -ldflags="-w -s" -o ngserve .

FROM scratch

COPY --from=builder /build/ngserve /

ENV PORT=8080 \
    APP_ROOT=/
    WEB_ROOT=./app
    NO_CACHE=false

EXPOSE 8080
VOLUME [ "/app" ]

ENTRYPOINT ["/ngserve"]
