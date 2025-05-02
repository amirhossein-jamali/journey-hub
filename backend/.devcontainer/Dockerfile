FROM golang:1.24-alpine

# Set the working directory
WORKDIR /app

# Install git, bash and other useful tools
RUN apk add --no-cache git bash curl make gcc libc-dev zip unzip python3 py3-pip openssh-client jq

# Install golangci-lint
RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.54.2

# Install swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Install govulncheck for dependency vulnerability checking
RUN go install golang.org/x/vuln/cmd/govulncheck@latest

# Install air for hot reload
RUN go install github.com/cosmtrek/air@latest

# Install mockgen for generating mocks
RUN go install github.com/golang/mock/mockgen@latest

# Install gotestsum for better test output
RUN go install gotest.tools/gotestsum@latest

# Install dlv for debugging
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# Install swagger-ui for API documentation
RUN mkdir -p /usr/local/swagger-ui \
    && cd /usr/local/swagger-ui \
    && wget https://github.com/swagger-api/swagger-ui/archive/refs/tags/v5.11.2.zip \
    && unzip v5.11.2.zip \
    && rm v5.11.2.zip \
    && mv swagger-ui-5.11.2/dist /usr/local/swagger-ui/dist \
    && rm -rf swagger-ui-5.11.2

# Add air configuration
COPY .air.toml .

# Set environment variables
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GO111MODULE=on \
    PATH="${PATH}:/root/.local/bin:/go/bin"

# Expose the port
EXPOSE 8080

# Command to start the development server with hot reload
CMD ["air", "-c", ".air.toml"] 