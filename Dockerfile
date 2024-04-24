# FROM alpine:latest
# RUN mkdir -p /home/addressbook
# WORKDIR /home/addressbook
# COPY .  .
# RUN  go build /test
# CMD [ "go","run","main.go" ]


# Build stage
# FROM golang:alpine AS builder
# WORKDIR /app
# COPY . .
# RUN go mod download
# RUN go build -o addressbook

# # Final stage
# FROM alpine:latest
# WORKDIR /home/addressbook
# COPY --from=builder /app/addressbook .
# CMD ["./addressbook"]

# Use the official Golang image as a base
FROM golang:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod .
COPY go.sum .

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Copy the .env file into the container
COPY .env .

# Build the Go app
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
