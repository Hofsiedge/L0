FROM golang:1.20-alpine AS builder

WORKDIR /app
COPY src/go.mod src/go.sum ./
# RUN apk update && apk upgrade &&
RUN go mod download && go mod verify

COPY src .
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o main cmd/stream-filler/main.go

FROM scratch
WORKDIR /
COPY --from=builder /app/main /app
CMD ["/app"]
