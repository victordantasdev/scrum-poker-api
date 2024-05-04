FROM golang:1.22.2-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o main .

FROM alpine:latest  
WORKDIR /root/
COPY --from=build /app/main .
CMD ["./main"]
