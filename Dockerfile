FROM golang:latest

WORKDIR /app

COPY . .

RUN go get -d -v ./...
RUN go build -o dns-server .

EXPOSE 53/udp

CMD ["./dns-server"]
