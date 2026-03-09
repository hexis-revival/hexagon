FROM golang:1.25.8 AS build

WORKDIR /app

# Update certificates
RUN apt install -y ca-certificates
RUN update-ca-certificates

# Install ffmpeg
RUN apt-get update && apt-get install -y ffmpeg

# Copy module files
COPY ./go.mod .
COPY ./go.sum .

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o ./hexagon

# Run the compiled binary
ENTRYPOINT ["./hexagon"]
