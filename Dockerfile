# Estágio 1: Build
FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod ./

# Copia o código e sincroniza o go.sum automaticamente
COPY . .
RUN go mod tidy
RUN go mod download

# Compila o binário estático
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server main.go

# Estágio 2: Imagem final mínima
FROM alpine:3.20

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/server .

EXPOSE 3000

CMD ["./server"]