#!/bin/bash

# Define color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Get current directory
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
ROOT_DIR="$(dirname "$DIR")"
CONFIG_DIR="$ROOT_DIR/config/development"

# Check if config directory exists
if [ ! -d "$CONFIG_DIR" ]; then
  echo -e "${RED}Error: Config directory not found: $CONFIG_DIR${NC}"
  echo "Please make sure you run this script from the project root directory."
  exit 1
fi

CONFIG_FILE="$CONFIG_DIR/config.yaml"

# Check if config file exists
if [ ! -f "$CONFIG_FILE" ]; then
  echo -e "${RED}Error: Config file not found: $CONFIG_FILE${NC}"
  exit 1
fi

# Check if the application is running
echo -e "${BLUE}Checking if the application is running...${NC}"
if ! curl -s http://localhost:8080/health > /dev/null; then
  echo -e "${YELLOW}Warning: Application does not seem to be running or health endpoint is not available.${NC}"
  echo "This test is more effective when the application is running."
  echo "Would you like to continue anyway? (y/n)"
  read -r response
  if [[ "$response" != "y" ]]; then
    echo "Exiting..."
    exit 0
  fi
fi

# Backup original config
BACKUP_FILE="$CONFIG_FILE.bak"
echo -e "${BLUE}Creating backup of config file to $BACKUP_FILE${NC}"
cp "$CONFIG_FILE" "$BACKUP_FILE"

# Function to restore original config
restore_config() {
  echo -e "\n${BLUE}Restoring original config file...${NC}"
  cp "$BACKUP_FILE" "$CONFIG_FILE"
  rm "$BACKUP_FILE"
  echo -e "${GREEN}Config restored.${NC}"
}

# Set trap to restore config on exit
trap restore_config EXIT

# Modify port in config file
echo -e "${BLUE}Modifying server port in config file...${NC}"
PORT=$(grep "port:" "$CONFIG_FILE" | awk '{print $2}')
NEW_PORT=$((PORT + 1))
sed -i "s/port: $PORT/port: $NEW_PORT/" "$CONFIG_FILE"

echo -e "${GREEN}Changed port from $PORT to $NEW_PORT${NC}"
echo -e "${YELLOW}The application should detect this change and reload its configuration.${NC}"
echo "Check the application logs for messages indicating config reload."

# Wait for a moment
echo "Waiting for 5 seconds..."
sleep 5

# Modify app name in config file
echo -e "\n${BLUE}Modifying app name in config file...${NC}"
APP_NAME=$(grep "name:" "$CONFIG_FILE" | head -1 | sed 's/name: "\(.*\)"/\1/')
NEW_APP_NAME="$APP_NAME (Modified)"
sed -i "s/name: \"$APP_NAME\"/name: \"$NEW_APP_NAME\"/" "$CONFIG_FILE"

echo -e "${GREEN}Changed app name from '$APP_NAME' to '$NEW_APP_NAME'${NC}"
echo -e "${YELLOW}The application should detect this change and reload its configuration.${NC}"
echo "Check the application logs for messages indicating config reload."

# Wait for user to continue
echo -e "\n${BLUE}Test completed. Press Enter to restore the original config...${NC}"
read -r 