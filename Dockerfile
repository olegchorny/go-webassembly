FROM golang:1.11-stretch

RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
#RUN rm main.go


RUN go get -u "github.com/boombuler/barcode"

RUN CGO_ENABLED=1 GOOS=js GOARCH=wasm go get -u "github.com/dennwc/dom"
RUN CGO_ENABLED=1 GOOS=js GOARCH=wasm go get -u "syscall/js" 

RUN CGO_ENABLED=1 GOARCH=wasm GOOS=js go build -o test.wasm main.go

RUN go build -o app server.go

EXPOSE 3000
CMD ["./app"]


