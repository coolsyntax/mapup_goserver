# Start from a base image containing the Go runtime
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the Go server code into the container
COPY . .

# Build the Go application
RUN go build -o app server.go

# Expose port 8000 to the outside world
EXPOSE 8000

# Command to run the Go server executable
CMD ["./app"]
