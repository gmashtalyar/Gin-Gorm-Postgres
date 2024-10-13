FROM golang:1.22-alpine


WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download the Go modules
RUN go mod download


COPY . .


EXPOSE 8080


CMD ["go", "run", "main.go"]
