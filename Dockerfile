FROM golang:1.22

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /gomax-ecs

ENV ECS_ENABLE_CONTAINER_METADATA=true

# Run
CMD ["/gomax-ecs"]
