FROM golang:1.23.0
WORKDIR /app

# Instalar nc (netcat) para o script de espera
RUN apt-get update && apt-get install -y netcat-openbsd

COPY . .
RUN go mod tidy && go build -o main ./main.go
CMD ["./main"]