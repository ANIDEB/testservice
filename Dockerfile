# Build stage
FROM golang:1.20-alpine AS builder
WORKDIR /src
COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /server ./

# Final stage
FROM alpine:3.18
RUN addgroup -S app && adduser -S -G app app
WORKDIR /app
COPY --from=builder /server /app/server
RUN chown app:app /app/server
USER app
EXPOSE 8080
ENV PORT=8080
CMD ["/app/server"]
