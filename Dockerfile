FROM golang:latest
LABEL authors="koeno"
WORKDIR /app
COPY . .
RUN go build -x cmd/perrApp/main.go
ENTRYPOINT ./main