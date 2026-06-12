# Build the application
FROM golang:alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project files
COPY . .

# Build the application binary
RUN go build -o main ./app/cmd

# Stage 2
FROM alpine:latest

RUN apk add --no-cache tzdata

# Set the working directory inside the container
WORKDIR /app

# Copy the entire project files
COPY --from=builder /app/main .

# Copy config and secret
COPY /config ./config

# ENV ENVIRONMENT=production

COPY .env .

# Command to run the application
ENTRYPOINT ["./main"]