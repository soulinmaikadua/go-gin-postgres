# Build Stage
FROM golang:latest AS build

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files to the working directory
COPY go.mod go.sum ./

# Download and install Go dependencies
RUN go mod download

# Copy the rest of the application source code to the working directory
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Runtime Stage
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the built executable from the build stage to the runtime image
COPY --from=build /app/main .

# Expose the port on which the Go application will run
EXPOSE 1234

# Command to run the Go application
CMD ["./main"]
