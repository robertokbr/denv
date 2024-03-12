package bucket

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Bucket struct {
	bucket     *s3.S3
	bucketName string
}

func NewS3Bucket() *S3Bucket {
	accessKey := os.Getenv("AWS_ACCESS_KEY")
	secretKey := os.Getenv("AWS_SECRET_KEY")
	bucketName := os.Getenv("AWS_BUCKET_NAME")

	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	}))

	bucket := s3.New(sess)

	s3Bucket := S3Bucket{
		bucket:     bucket,
		bucketName: bucketName,
	}

	return &s3Bucket
}

func (s3b *S3Bucket) UploadFile(filePath, name string) {
	fmt.Println("🚚 Upload in progress...")

	file, err := os.Open(filePath)

	if err != nil {
		log.Fatalf("Failed to read file: %s", err.Error())
	}

	defer file.Close()

	fileStat, _ := file.Stat()
	size := fileStat.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	_, err = s3b.bucket.PutObject(&s3.PutObjectInput{
		Bucket:             aws.String(s3b.bucketName),
		Key:                aws.String(fmt.Sprintf("%s.txt", name)),
		Body:               strings.NewReader(string(buffer)),
		ACL:                aws.String("private"),
		ContentDisposition: aws.String("attachment"),
		ContentType:        aws.String("application/octet-stream"),
	})

	if err != nil {
		log.Fatalf("Failed to upload file to s3: %s", err.Error())
	}

	fmt.Println("🥳 Filed uploaded!!!")
}

func (s3b *S3Bucket) DownloadFile(name, outputName string) {
	fmt.Println("🚚 Download in progress...")

	res, err := s3b.bucket.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s3b.bucketName),
		Key:    aws.String(fmt.Sprintf("%s.txt", name)),
	})

	if err != nil {
		log.Fatalf("Failed to download the file: %s", err.Error())
	}

	defer res.Body.Close()

	var fileName string

	if outputName != "" {
		fileName = outputName
	} else {
		fileName = fmt.Sprintf("%s.txt", name)
	}

	file, err := os.Create(fileName)

	if err != nil {
		log.Fatalf("Failed to create env file: %s", err.Error())
	}

	defer file.Close()

	_, err = io.Copy(file, res.Body)

	if err != nil {
		log.Fatalf("Failed to download file: %s", err.Error())
	}

	fmt.Println("🥳 Download succeed!!!")
}

func (s3b *S3Bucket) ListFiles() {
	fmt.Println("🚚 List in progress...")

	res, err := s3b.bucket.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(s3b.bucketName),
	})

	if err != nil {
		log.Fatalf("Failed to list files: %s", err.Error())
	}

	fmt.Println("🥳 Files in the bucket:")

	for _, item := range res.Contents {
		key := *item.Key
		fmt.Println(key[:len(key)-4])
	}
}

func (s3b *S3Bucket) DeleteFile(name string) {
	fmt.Println("🚚 Delete in progress...")

	_, err := s3b.bucket.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s3b.bucketName),
		Key:    aws.String(fmt.Sprintf("%s.txt", name)),
	})

	if err != nil {
		log.Fatalf("Failed to delete file: %s", err.Error())
	}

	fmt.Println("🥳 File deleted!!!")
}
