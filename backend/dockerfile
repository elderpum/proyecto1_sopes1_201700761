# Use the official Golang image to create a binary.
# This is based on Debian and sets the GOPATH to /go.
# WORKDIR has to be /go as it's the GOPATH.
FROM golang:1.18.1 AS builder
WORKDIR /go/src/app
COPY . .
# Fetch dependencies.
# Using go get.
RUN go get -d -v ./...
# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Use a minimal base image to package the final artifact.
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
# Copy the pre-built binary from the previous stage.
COPY --from=builder /go/src/app/main .
# Expose port 8080 to the outside world
EXPOSE 8080
# Command to run the executable
CMD ["./main"]