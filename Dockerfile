FROM golang:latest as builder
WORKDIR /go/src/github.com/zate/corgibot/
RUN go get -u github.com/golang/dep/cmd/dep
ADD main.go .
RUN dep init && dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o corgibot main.go

FROM scratch
WORKDIR /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/src/github.com/zate/corgibot/corgibot .
COPY .secrets.yaml .

CMD ["/corgibot"]