package helpers

import (
	"log"
	"os"
	"path"

	"github.com/joho/godotenv"
)

const (
	READ_WRITE_AND_EXECUTE_PERMISSION = 0755
)

var (
	PROJECT_PATH string
	ENV_PATH     string
)

func StartProject() {
	config, err := os.UserHomeDir()

	if err != nil {
		log.Fatalf("Failed to get user config dir: %v", err)
	}

	PROJECT_PATH = path.Join(config, ".config", "denv")

	ENV_PATH = path.Join(PROJECT_PATH, ".env")

	if _, err := os.Stat(PROJECT_PATH); os.IsNotExist(err) {
		err = os.Mkdir(PROJECT_PATH, READ_WRITE_AND_EXECUTE_PERMISSION)

		if err != nil {
			log.Fatalf("Failed to read denv dir: %v", err)
		}

		_, err = os.Create(ENV_PATH)

		if err != nil {
			log.Fatalf("Error to write env file: %v", err)
		}
	}

	err = godotenv.Load(ENV_PATH)

	if err != nil {
		log.Fatalf("Failed to read envs: %v", err)
	}
}

func CheckEnvs() bool {
	err := godotenv.Load(ENV_PATH)

	if err != nil {
		log.Fatalf("Failed to read envs: %v", err)
	}

	accessKey := os.Getenv("AWS_ACCESS_KEY")
	secretKey := os.Getenv("AWS_SECRET_KEY")
	bucketName := os.Getenv("AWS_BUCKET_NAME")

	if accessKey == "" || secretKey == "" || bucketName == "" {
		return false
	}

	return true
}
