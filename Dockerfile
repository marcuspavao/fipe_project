FROM golang:1.23

WORKDIR /app

# Copia apenas go.mod e go.sum primeiro
COPY go.mod go.sum ./

# Executa o tidy (que vai fazer download também)
RUN go mod tidy

RUN go mod download

RUN go install github.com/air-verse/air@latest

# Agora copia o restante do código
COPY . .

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]