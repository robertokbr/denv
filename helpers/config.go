package helpers

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func SetupConfig() {
	var AWS_ACCESS_KEY string
	var AWS_SECRET_KEY string
	var AWS_BUCKET_NAME string
	var AWS_BUCKET_REGION string

	fmt.Println("🚧 Insert your AWS Access key")
	fmt.Scan(&AWS_ACCESS_KEY)
	fmt.Println("🚧 Insert your AWS Secret key")
	fmt.Scan(&AWS_SECRET_KEY)
	fmt.Println("🚧 Insert your AWS Bucket name")
	fmt.Scan(&AWS_BUCKET_NAME)
	fmt.Println("🚧 Insert your AWS Bucket region")
	fmt.Scan(&AWS_BUCKET_REGION)

	file, err := os.Create(ENV_PATH)

	if err != nil {
		log.Fatalf("Failed to write env file: %v", err)
	}

	data := "AWS_ACCESS_KEY=%s\nAWS_SECRET_KEY=%s\nAWS_BUCKET_NAME=%s\nAWS_BUCKET_REGION=%s"

	data = fmt.Sprintf(data, AWS_ACCESS_KEY, AWS_SECRET_KEY, AWS_BUCKET_NAME, AWS_BUCKET_REGION)

	sdata := strings.NewReader(data)

	_, err = io.Copy(file, sdata)

	if err != nil {
		log.Fatalf("Failed to copy data to env file: %v", err)
	}

	fmt.Println("🔥 Thank you! Everything is right!")
	fmt.Println("🤓 Type denv --help if you want to see how to use the CLI.")
	fmt.Println("🫢 Type denv --config again if the CLI is not working properly.")
}

func PrintHelp() {
	fmt.Println("denv --config to start the CLI configuration")
	fmt.Println("denv --up [file path] --name [file nickname] to upload some env file")
	fmt.Println("denv --name [file nickname] to download some env file you have uploaded")
	fmt.Println("denv --name [file nickname] --out [file name] to download some env file you have uploaded with some specific name")
	fmt.Println("denv --list to list all files in the bucket")
	fmt.Println("denv --del [file nickname] to delete some file in the bucket")
	fmt.Println("denv --rename [file nickname] --name [new nickname] to rename a file in the bucket")
}

func HandleEnvError() {
	fmt.Println("🤔 Hello! Type 'denv --setup' to start setting up the application")
}

func CheckEnvs() error {
	err := godotenv.Load(ENV_PATH)

	if err != nil {
		log.Fatalf("Failed to read envs: %s", err.Error())
	}

	accessKey := os.Getenv("AWS_ACCESS_KEY")
	secretKey := os.Getenv("AWS_SECRET_KEY")
	bucketName := os.Getenv("AWS_BUCKET_NAME")

	if accessKey == "" || secretKey == "" || bucketName == "" {
		return errors.New("envs are not set")
	}

	return nil
}