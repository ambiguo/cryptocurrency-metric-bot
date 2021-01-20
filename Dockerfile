FROM golang:1.12.0-alpine3.9 as builder

RUN apk add git

RUN mkdir /app
ADD . /app
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

CMD ["./main"]

