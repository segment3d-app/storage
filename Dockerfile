FROM golang:1.21.6-alpine3.19 as builder

WORKDIR /app

COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o main .

FROM alpine:latest  

RUN apk --no-cache add ca-certificates ffmpeg curl

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/.env .
COPY --from=builder /app/Makefile .
RUN mkdir -p /app/files

CMD ["./main"]