FROM golang
EXPOSE 8080
WORKDIR /go/src/app
RUN go get -v github.com/mesg-foundation/core/api/service && \
    go get -v github.com/mesg-foundation/core/service && \
    go get -v google.golang.org/grpc && \
    go get -v github.com/satori/go.uuid

COPY . .
RUN go install -v ./...
RUN go build
CMD ["./app" ]
