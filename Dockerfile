FROM golang:1.22.7 AS build

WORKDIR /app

# Update certificates
RUN apt install -y ca-certificates
RUN update-ca-certificates

# Install ffmpeg
RUN apt-get update && apt-get install -y ffmpeg

# Copy module files
COPY ./go.mod .
COPY ./go.sum .
COPY ./hnet/go.mod ./hnet/go.mod
COPY ./hnet/go.sum ./hnet/go.sum
COPY ./common/go.mod ./common/go.mod
COPY ./common/go.sum ./common/go.sum
COPY ./hscore/go.mod ./hscore/go.mod
COPY ./hscore/go.sum ./hscore/go.sum

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o ./hexagon

# Run the compiled binary
ENTRYPOINT ["./hexagon"]