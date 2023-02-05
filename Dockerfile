FROM golang:latest

# Configure working directory
WORKDIR /app

# Download go packages & verify them
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Live reload
RUN go install github.com/cosmtrek/air@latest

# Initialize air
CMD ["air"]
