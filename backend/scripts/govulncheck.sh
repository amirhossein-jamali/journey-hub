#!/bin/bash

# govulncheck.sh - Script to check Go dependencies for vulnerabilities
# This script is part of the Journey Hub security workflow

# Print colored output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Journey Hub - Vulnerability Scanner${NC}"
echo "=========================================="

# Check if govulncheck is installed
if ! command -v govulncheck &> /dev/null; then
    echo -e "${YELLOW}Installing govulncheck...${NC}"
    go install golang.org/x/vuln/cmd/govulncheck@latest
    
    # Verify installation
    if ! command -v govulncheck &> /dev/null; then
        echo -e "${RED}Failed to install govulncheck. Please install it manually:${NC}"
        echo "go install golang.org/x/vuln/cmd/govulncheck@latest"
        exit 1
    fi
    echo -e "${GREEN}govulncheck installed successfully.${NC}"
fi

echo -e "${YELLOW}Scanning for vulnerabilities...${NC}"
echo "=========================================="

# Run govulncheck with output formatting
govulncheck ./...

# Check the exit status
STATUS=$?
if [ $STATUS -eq 0 ]; then
    echo -e "${GREEN}No vulnerabilities found!${NC}"
else
    echo -e "${RED}Vulnerabilities detected! Please review the findings above.${NC}"
    echo "For more information about the vulnerabilities, visit: https://pkg.go.dev/golang.org/x/vuln"
fi

exit $STATUS 