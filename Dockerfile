FROM golang:1.22

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY . .
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /gomax-ecs

# Run
CMD ["/gomax-ecs"]
