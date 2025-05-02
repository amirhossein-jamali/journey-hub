package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/amirhossein-jamali/journey-hub/config"
)

func main() {
	// Parse command line arguments
	var value string
	var keyEnv string
	var interactive bool

	flag.StringVar(&value, "value", "", "Value to encrypt")
	flag.StringVar(&keyEnv, "key-env", config.EnvKeyName, "Environment variable name for encryption key")
	flag.BoolVar(&interactive, "i", false, "Interactive mode (prompts for value)")
	flag.Parse()

	// Set environment variable for encryption key if provided
	if keyEnv != config.EnvKeyName {
		os.Setenv(config.EnvKeyName, os.Getenv(keyEnv))
	}

	// If in interactive mode, prompt for value
	if interactive {
		fmt.Print("Enter value to encrypt (input will be hidden): ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %s\n", err)
			os.Exit(1)
		}
		value = strings.TrimSpace(input)
	}

	// Check if value is provided
	if value == "" {
		fmt.Println("Error: No value provided. Use -value flag or -i for interactive mode.")
		os.Exit(1)
	}

	// Encrypt the value
	encrypted, err := config.EncryptConfigValue(value)
	if err != nil {
		fmt.Printf("Error encrypting value: %s\n", err)
		os.Exit(1)
	}

	// Print the encrypted value
	fmt.Println("Encrypted Value:")
	fmt.Println(encrypted)
	fmt.Println("\nYou can use this value in your config files.")
}
