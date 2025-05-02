#!/bin/bash

# Create hooks directory if it doesn't exist
mkdir -p .git/hooks

# Function to install a hook
install_hook() {
    HOOK_NAME=$1
    HOOK_PATH=".git/hooks/$HOOK_NAME"
    
    echo "Installing $HOOK_NAME hook..."
    
    # Copy the hook file
    cp "backend/scripts/hooks/$HOOK_NAME" "$HOOK_PATH"
    
    # Make it executable
    chmod +x "$HOOK_PATH"
    
    echo "$HOOK_NAME hook installed successfully!"
}

# Install hooks
install_hook "pre-commit"

echo "All Git hooks have been installed!" 