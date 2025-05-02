#!/bin/bash

# Exit on error
set -e

# Color for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

echo "${YELLOW}Initializing Swagger documentation...${NC}"

# Check if swag is installed
if ! command -v swag &> /dev/null; then
    echo "${RED}swag is not installed. Installing...${NC}"
    go install github.com/swaggo/swag/cmd/swag@latest
fi

# Navigate to backend directory
cd "$(dirname "$0")/.." || exit

# Create docs directory if it doesn't exist
mkdir -p api/swagger/gen

# Generate Swagger docs
echo "${YELLOW}Generating Swagger documentation...${NC}"
swag init --parseDependency --parseInternal --parseDepth 1 \
    --dir ./cmd/api,./internal/interfaces/api/http/handlers \
    --generatedTime \
    --outputTypes json,yaml \
    --output ./api/swagger/gen

# Check if the generation was successful
if [ $? -eq 0 ]; then
    echo "${GREEN}Swagger documentation generated successfully!${NC}"
    echo "${YELLOW}Documentation available at: ${GREEN}http://localhost:8080/swagger/index.html${NC}"
else
    echo "${RED}Failed to generate Swagger documentation!${NC}"
    exit 1
fi 