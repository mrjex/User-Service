FROM golang:alpine
WORKDIR /
COPY . .
RUN go mod download
RUN go build -o user-service main.go
CMD ["./user-service"]

