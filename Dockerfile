FROM golang:1.22.2-bookworm

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o /app.bin ./cmd/server/main.go

EXPOSE 5000

CMD ["/app.bin"]