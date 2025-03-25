package config

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type AWSCredentials struct {
	AccessKey   string
	SecretKey   string
	BucketName  string
	BucketRegion string
}

func SetupEnvironment() error {
	// Ensure project directory exists
	if err := ensureProjectDir(); err != nil {
		return err
	}
	
	// Load environment variables
	return loadEnv()
}

func ensureProjectDir() error {
	if _, err := os.Stat(ProjectPath); err != nil {
		err = os.MkdirAll(ProjectPath, ReadWriteExecutePermission)
		if err != nil {
			return fmt.Errorf("failed to create denv directory: %s", err.Error())
		}
		
		_, err = os.Create(EnvPath)
		if err != nil {
			return fmt.Errorf("failed to create env file: %s", err.Error())
		}
	}
	
	return nil
}

func loadEnv() error {
	if _, err := os.Stat(EnvPath); err != nil {
		// File doesn't exist, create it
		_, err = os.Create(EnvPath)
		if err != nil {
			return fmt.Errorf("failed to create env file: %s", err.Error())
		}
		return nil
	}
	
	err := godotenv.Load(EnvPath)
	if err != nil {
		return fmt.Errorf("failed to load environment: %s", err.Error())
	}
	
	return nil
}

func ValidateEnvironment() error {
	err := godotenv.Load(EnvPath)
	if err != nil {
		return fmt.Errorf("failed to read environment: %s", err.Error())
	}
	
	accessKey := os.Getenv("AWS_ACCESS_KEY")
	secretKey := os.Getenv("AWS_SECRET_KEY")
	bucketName := os.Getenv("AWS_BUCKET_NAME")
	
	if accessKey == "" || secretKey == "" || bucketName == "" {
		return errors.New("environment variables not properly set")
	}
	
	return nil
}

func SaveCredentials(creds AWSCredentials) error {
	file, err := os.Create(EnvPath)
	if err != nil {
		return fmt.Errorf("failed to create env file: %s", err.Error())
	}
	defer file.Close()
	
	data := "AWS_ACCESS_KEY=%s\nAWS_SECRET_KEY=%s\nAWS_BUCKET_NAME=%s\nAWS_BUCKET_REGION=%s"
	data = fmt.Sprintf(data, creds.AccessKey, creds.SecretKey, creds.BucketName, creds.BucketRegion)
	
	_, err = io.Copy(file, strings.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to write credentials: %s", err.Error())
	}
	
	return nil
}

func GetAWSCredentials() AWSCredentials {
	return AWSCredentials{
		AccessKey:    os.Getenv("AWS_ACCESS_KEY"),
		SecretKey:    os.Getenv("AWS_SECRET_KEY"),
		BucketName:   os.Getenv("AWS_BUCKET_NAME"),
		BucketRegion: os.Getenv("AWS_BUCKET_REGION"),
	}
}