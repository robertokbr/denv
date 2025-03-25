package cli

import (
	"fmt"
	"log"

	"github.com/robertokbr/denv/config"
)

func ConfigureApplication() {
	var creds config.AWSCredentials
	
	fmt.Println("🚧 Insert your AWS Access key")
	fmt.Scan(&creds.AccessKey)
	
	fmt.Println("🚧 Insert your AWS Secret key")
	fmt.Scan(&creds.SecretKey)
	
	fmt.Println("🚧 Insert your AWS Bucket name")
	fmt.Scan(&creds.BucketName)
	
	fmt.Println("🚧 Insert your AWS Bucket region")
	fmt.Scan(&creds.BucketRegion)
	
	err := config.SaveCredentials(creds)
	if err != nil {
		log.Fatalf("Failed to save configuration: %v", err)
	}
	
	PrintSuccessConfig()
}