package main

import (
	"fmt"

	"github.com/amirhossein-jamali/journey-hub/config"
)

func main() {
	// لود کردن کانفیگ
	conf, err := config.LoadConfig("./config/development")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	// نمایش مقادیر رمزگشایی شده
	fmt.Println("=== Decrypted Values ===")
	fmt.Println("Database Password:", conf.Database.Password)
	fmt.Println("Redis Password:", conf.Redis.Password)
	fmt.Println("Storage Secret Key:", conf.Storage.SecretAccessKey)
	fmt.Println("JWT Secret:", conf.JWT.Secret)
}
