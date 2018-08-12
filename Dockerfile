FROM golang:1.8.3 as builder
WORKDIR /go/src/github.com/mmussett/simple_cache_service
RUN go get -d -v github.com/gorilla/mux
RUN go get -d -v github.com/patrickmn/go-cache
COPY main.go  .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o simple-cache-service .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/mmussett/simple_cache_service .
CMD ["./simple-cache-service"]