# Use the official Go image as a parent image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the local configuration files into the container
COPY go.mod .
COPY go.sum .

# Download and install dependencies
RUN go mod download

# Copy the rest of the application's code
COPY . .

# Build the application
RUN go build -o main ./main.go

# Command to run the executable
CMD ["./main"]