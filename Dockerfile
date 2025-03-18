FROM golang:1.24

WORKDIR /go/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags "-s -w"

CMD ["./vinted_go"]