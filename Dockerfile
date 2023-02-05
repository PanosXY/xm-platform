FROM golang:latest

# Configure working directory
WORKDIR /app

# Download go packages & verify them
COPY ./src/go.mod ./src/go.sum ./
RUN go mod download && go mod verify

# Live reload
RUN go install github.com/cosmtrek/air@latest

# Initialize air
CMD ["air"]
