FROM golang:1.22

RUN go install github.com/air-verse/air@latest

WORKDIR /opt/user

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8009

CMD ["air"]
