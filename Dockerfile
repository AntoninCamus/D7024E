FROM golang:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Install system dependancies
RUN apt-get update && apt-get install protobuf-compiler -y
RUN go get -u github.com/golang/protobuf/protoc-gen-go

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN ./build.sh

# Expose port 8080 & 9090 to the outside world
EXPOSE 8080
EXPOSE 9090

# Command to run the executable
CMD ["./main.run"]