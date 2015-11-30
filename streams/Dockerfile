FROM golang
EXPOSE 8080
ADD . /go/src/github.com/ello/ello-go/streams/
RUN go get github.com/Masterminds/glide
RUN go build github.com/Masterminds/glide
WORKDIR /go/src/github.com/ello/ello-go/streams/
RUN GO15VENDOREXPERIMENT=1 glide install
RUN GO15VENDOREXPERIMENT=1 go build
ENTRYPOINT ["./streams"]
