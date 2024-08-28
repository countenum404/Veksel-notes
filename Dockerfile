FROM golang:1.22.6-alpine3.20

RUN apk add --no-cache make
WORKDIR /go/veksel
COPY ./ ./
RUN make build

EXPOSE 4567

CMD ["./veksel"]