FROM golang:1.5.1
EXPOSE 8080

COPY ./glide.yaml /go/src/github.com/ello/streams/glide.yaml
RUN go get github.com/Masterminds/glide
RUN go build github.com/Masterminds/glide
WORKDIR /go/src/github.com/ello/streams/
RUN GO15VENDOREXPERIMENT=1 glide install

COPY . /go/src/github.com/ello/streams/
RUN GO15VENDOREXPERIMENT=1 go build

CMD ["./streams"]
