FROM golang:1.22.6-alpine3.20

WORKDIR /go/veksel
COPY ./ ./
RUN go build -o veksel cmd/main/main.go

EXPOSE 4567

CMD ["./veksel"]