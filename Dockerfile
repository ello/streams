FROM golang:1.17.5
EXPOSE 8080

WORKDIR /go/src/github.com/ello/streams/

COPY . /go/src/github.com/ello/streams/
RUN go build ./..

CMD ["./streams"]
