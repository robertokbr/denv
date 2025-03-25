package bucket

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path"
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

func NewS3Bucket(accessKey, secretKey, bucketName, bucketRegion string) *S3Bucket {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(bucketRegion),
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
	if fileStat.IsDir() {
		fmt.Println("🚧 You can't upload a directory")
		return
	}

	size := fileStat.Size()
	ext := path.Ext(filePath)
	buffer := make([]byte, size)
	file.Read(buffer)

	_, err = s3b.bucket.PutObject(&s3.PutObjectInput{
		Bucket:             aws.String(s3b.bucketName),
		Key:                aws.String(fmt.Sprintf("%s.%s", name, ext[1:])),
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
		Key:    aws.String(name),
	})

	if err != nil {
		log.Fatalf("Failed to download the file: %s", err.Error())
	}
	defer res.Body.Close()

	fileName := name
	if outputName != "" {
		fileName = outputName
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

	files, err := s3b.getFilesList()
	if err != nil {
		log.Fatalf("Failed to list files: %s", err.Error())
	}

	fmt.Println("🥳 Files in the bucket:")

	if len(files.Contents) == 0 {
		fmt.Println("No files found in the bucket.")
		return
	}

	fmt.Printf("%-40s | %-20s\n", "File Name", "Last Modified")

	for _, item := range files.Contents {
		key := *item.Key
		lastModified := item.LastModified.Format("2006-01-02 15:04:05")
		fmt.Printf("%-40s | %-20s\n", key, lastModified)
	}
}

// getFilesList returns the raw S3 ListObjectsOutput
func (s3b *S3Bucket) getFilesList() (*s3.ListObjectsOutput, error) {
	return s3b.bucket.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(s3b.bucketName),
	})
}

// ListFileNames returns just the names of files in the bucket for use with autocomplete
func (s3b *S3Bucket) ListFileNames() ([]string, error) {
	res, err := s3b.getFilesList()
	if err != nil {
		return nil, err
	}
	
	fileNames := make([]string, 0, len(res.Contents))
	for _, item := range res.Contents {
		if item.Key != nil {
			fileNames = append(fileNames, *item.Key)
		}
	}
	
	return fileNames, nil
}

func (s3b *S3Bucket) DeleteFile(name string) {
	fmt.Println("🚚 Delete in progress...")

	_, err := s3b.bucket.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s3b.bucketName),
		Key:    aws.String(name),
	})

	if err != nil {
		log.Fatalf("Failed to delete file: %s", err.Error())
	}

	fmt.Println("🥳 File deleted!!!")
}

func (s3b *S3Bucket) RenameFile(oldName, newName string) {
	fmt.Println("🚚 Rename in progress...")
	
	// First, copy the file with the new name
	res, err := s3b.bucket.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s3b.bucketName),
		Key:    aws.String(oldName),
	})
	
	if err != nil {
		log.Fatalf("Failed to find file %s: %s", oldName, err.Error())
	}
	defer res.Body.Close()
	
	// Extract file extension from the old name
	ext := path.Ext(oldName)
	var newKey string
	
	// If old name has an extension, use it for the new name
	if ext != "" {
		baseName := strings.TrimSuffix(newName, ext)
		newKey = baseName + ext
	} else {
		newKey = newName
	}
	
	// Read the entire body
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Failed to read file contents: %s", err.Error())
	}
	
	// Copy object to the new key
	_, err = s3b.bucket.PutObject(&s3.PutObjectInput{
		Bucket:             aws.String(s3b.bucketName),
		Key:                aws.String(newKey),
		Body:               bytes.NewReader(bodyBytes),
		ACL:                aws.String("private"),
		ContentDisposition: aws.String("attachment"),
		ContentType:        aws.String("application/octet-stream"),
	})
	
	if err != nil {
		log.Fatalf("Failed to rename file to %s: %s", newKey, err.Error())
	}
	
	// Delete the old object
	_, err = s3b.bucket.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s3b.bucketName),
		Key:    aws.String(oldName),
	})
	
	if err != nil {
		log.Fatalf("Failed to delete original file after rename: %s", err.Error())
	}
	
	fmt.Printf("🥳 File renamed from %s to %s!!!\n", oldName, newKey)
}
