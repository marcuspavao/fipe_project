FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

RUN go mod download

RUN go install github.com/air-verse/air@v1.60.0

COPY . .

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]
