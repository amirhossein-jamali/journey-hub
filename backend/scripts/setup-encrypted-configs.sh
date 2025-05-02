#!/bin/bash

# Define color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${YELLOW}WARNING: This script is intended for development environments only.${NC}"
echo -e "${YELLOW}DO NOT run this in CI/CD pipelines or automated environments.${NC}"
echo -e "${YELLOW}For CI/CD, set JH_SKIP_CONFIG_ENCRYPTION=true environment variable.${NC}"
echo ""

# Check if running in CI/CD environment
if [ -n "$CI" ] || [ -n "$GITHUB_ACTIONS" ]; then
    echo -e "${RED}ERROR: This script should not be run in CI/CD environments.${NC}"
    echo "Set JH_SKIP_CONFIG_ENCRYPTION=true instead."
    exit 1
fi

# Get current directory
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
ROOT_DIR="$(dirname "$DIR")"
CONFIG_DIR="$ROOT_DIR/config"

# Check if config directory exists
if [ ! -d "$CONFIG_DIR" ]; then
  echo -e "${RED}Error: Config directory not found: $CONFIG_DIR${NC}"
  echo "Please make sure you run this script from the project root directory."
  exit 1
fi

# Generate a random encryption key and set it as environment variable
echo -e "${BLUE}Generating encryption key...${NC}"
ENCRYPTION_KEY=$(openssl rand -base64 32)
export JH_CONFIG_ENCRYPTION_KEY="$ENCRYPTION_KEY"
echo -e "${GREEN}Encryption key generated and set as JH_CONFIG_ENCRYPTION_KEY environment variable.${NC}"
echo "For future use, save this key securely:"
echo "$ENCRYPTION_KEY"
echo ""

# Function to encrypt a value and update config file
encrypt_and_update() {
    local env=$1
    local section=$2
    local key=$3
    local value=$4
    local config_file="$CONFIG_DIR/$env/config.yaml"
    
    if [ ! -f "$config_file" ]; then
        echo -e "${RED}Error: Config file not found: $config_file${NC}"
        return 1
    fi
    
    echo -e "${BLUE}Encrypting $section.$key for $env environment...${NC}"
    
    # Encrypt the value
    encrypted=$(go run "$ROOT_DIR/cmd/encrypt/main.go" -value "$value")
    
    # Extract just the encrypted value (second line of output)
    encrypted_value=$(echo "$encrypted" | sed -n '2p')
    
    if [[ -z "$encrypted_value" ]]; then
        echo -e "${RED}Error: Failed to encrypt value${NC}"
        return 1
    fi
    
    # Update the config file
    # Find the line with the key and replace everything after the colon
    if grep -q "^  $key:" "$config_file"; then
        sed -i "s|^  $key:.*|  $key: \"$encrypted_value\"|" "$config_file"
        echo -e "${GREEN}Updated $section.$key in $env environment${NC}"
    else
        echo -e "${YELLOW}Warning: Could not find $section.$key in $config_file${NC}"
    fi
}

# Process each environment
environments=("development" "test" "production")

for env in "${environments[@]}"; do
    echo -e "\n${BLUE}Processing $env environment...${NC}"
    
    # Database password
    encrypt_and_update "$env" "database" "password" "secure_db_password"
    
    # Redis password
    encrypt_and_update "$env" "redis" "password" "secure_redis_password"
    
    # MinIO secret key
    encrypt_and_update "$env" "storage" "secret_access_key" "secure_minio_secret_key"
    
    # JWT secret
    jwt_secret=$(openssl rand -base64 32)
    encrypt_and_update "$env" "jwt" "secret" "$jwt_secret"
    
    echo -e "${GREEN}Finished setting up encrypted values for $env environment${NC}"
done

echo -e "\n${GREEN}All environments have been set up with encrypted values!${NC}"
echo -e "${YELLOW}Important: Store your encryption key securely for future use:${NC}"
echo "$ENCRYPTION_KEY"
echo -e "${YELLOW}Set this as the JH_CONFIG_ENCRYPTION_KEY environment variable in your production environment.${NC}" 