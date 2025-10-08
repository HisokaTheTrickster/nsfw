# Stage 1: Build the Go app
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod .
RUN go version
RUN echo "---- go.mod contents ----" && cat go.mod && echo "------------------------"
RUN go mod download
COPY . .
RUN go build -o nsfw main.go


# Stage 2: Create a lightweight final image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/nsfw .
COPY --from=builder /app/records.json .
EXPOSE 53/udp
CMD ["./nsfw"]
