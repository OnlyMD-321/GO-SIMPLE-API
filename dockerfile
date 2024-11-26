# Step 1: Use Golang base image to build the app
FROM golang:1.23.2-alpine AS builder

# Step 2: Set the working directory inside the container
WORKDIR /app

# Step 3: Copy go.mod and go.sum files to handle dependencies
COPY go.mod go.sum ./

# Step 4: Download the Go dependencies
RUN go mod download

# Step 5: Copy the application code into the container
COPY ./cmd /app
COPY ./internal /app/internal
COPY ./pkg /app/pkg

# Step 6: Build the Go application
RUN go build -o go-simple-api .

# Step 7: Use a smaller image (Alpine) for the final image to reduce size
FROM alpine:latest

# Step 8: Copy the built binary from the builder stage
COPY --from=builder /app/go-simple-api /go-simple-api

# Step 9: Set the working directory for the runtime container
WORKDIR /root/

# Step 10: Expose the port the application will run on
EXPOSE 8080

# Step 11: Command to run the app when the container starts
CMD ["/go-simple-api"]
