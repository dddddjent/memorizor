FROM golang:alpine as compiler
WORKDIR /src

RUN go env -w  GOPROXY=https://goproxy.cn,direct 
RUN go install github.com/cespare/reflex@latest

COPY go.* ./
RUN go mod download

COPY . .
RUN go build -o ./build/run


FROM alpine:latest
WORKDIR /

RUN apk add --no-cache ca-certificates
EXPOSE 8080

COPY --from=compiler /src/build/run ./bin
CMD ./bin/run
