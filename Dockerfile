FROM golang:1.11-stretch

RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
#RUN rm main.go




RUN CGO_ENABLED=1 go get -u "github.com/dennwc/dom"
RUN CGO_ENABLED=1 go get -u "github.com/boombuler/barcode"
#RUN go get -u "syscall/js" 

RUN CGO_ENABLED=1 GOARCH=wasm GOOS=js go build -o test.wasm main.go

RUN CGO_ENABLED=0 GOOS=linux go build -o app server.go

RUN go build -o app server.go

EXPOSE 8080
CMD ["./app"]


