FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git
RUN go install github.com/a-h/templ/cmd/templ@latest

WORKDIR /app
COPY . .
RUN go mod tidy
RUN templ generate
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server .

FROM alpine:3.23
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/static ./static
EXPOSE 3000
CMD ["./server"]
