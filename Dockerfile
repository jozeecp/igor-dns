FROM golang:latest

WORKDIR /app

COPY . .

RUN go get -d -v ./...
RUN go build main

EXPOSE 53/udp

CMD ["./main"]
