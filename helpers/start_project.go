package helpers

import (
	"errors"
	"log"
	"os"
	"path"
	"runtime"

	"github.com/joho/godotenv"
)

const (
	READ_WRITE_AND_EXECUTE_PERMISSION = 0755
)

var (
	PROJECT_PATH string
	ENV_PATH     string
)

func setPaths(home string) {
	isWindows := runtime.GOOS == "windows"

	if isWindows {
		PROJECT_PATH = path.Join(home, "denv")
	} else {
		PROJECT_PATH = path.Join(home, ".config", "denv")
	}

	ENV_PATH = path.Join(PROJECT_PATH, ".env")
}

func StartProject() {
	home, err := os.UserHomeDir()

	if err != nil {
		log.Fatalf("Failed to get user config dir: %s", err.Error())
	}

	setPaths(home)

	if _, err := os.Stat(PROJECT_PATH); err != nil {
		err = os.Mkdir(PROJECT_PATH, READ_WRITE_AND_EXECUTE_PERMISSION)

		if err != nil {
			log.Fatalf("Failed to read denv dir: %s", err.Error())
		}

		_, err = os.Create(ENV_PATH)

		if err != nil {
			log.Fatalf("Error to write env file: %s", err.Error())
		}
	}

	err = godotenv.Load(ENV_PATH)

	if err != nil {
		log.Fatalf("Failed to read envs: %s", err.Error())
	}
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
