FROM golang:1.11
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN GOARCH=wasm GOOS=js go build -o test.wasm main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o app server.go

EXPOSE 8080
CMD ["./app"]


