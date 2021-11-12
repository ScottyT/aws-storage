package main

import (
	"aws-storage/middleware"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
)

var AccessKeyID string
var SecretAccessKey string

func GetEnvWithKey(key string) string {
	return os.Getenv(key)
}

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
		os.Exit(1)
	}
}
func main() {
	if len(os.Args) != 2 {
		exitErrorf("Bucket name required\nUsage: go run", os.Args[0], "BUCKET")
	}
	LoadEnv()
	bucket := os.Args[1]
	filePath := "shinobu.jpg"
	fmt.Print(bucket)
	AccessKeyID = GetEnvWithKey("AWS_ACCESS_KEY_ID")
	SecretAccessKey = GetEnvWithKey("AWS_SECRET_ACCESS_KEY")
	sess := middleware.ConnectAws()

	// Create S3 service client
	svc := s3.New(sess)
	err := uploadFileToS3(svc, bucket, filePath)
	if err != nil {
		log.Fatalf("could not upload file: %v", err)
	}

}
func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
func uploadFileToS3(s3Client *s3.S3, bucketName string, filePath string) error {
	// Get the fileName from Path
	fileName := filepath.Base(filePath)

	// Open the file from the file path
	upFile, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not open local filepath [%v]: %+v", filePath, err)
	}
	defer upFile.Close()

	// Get the file info
	upFileInfo, _ := upFile.Stat()
	var fileSize int64 = upFileInfo.Size()
	fileBuffer := make([]byte, fileSize)
	upFile.Read(fileBuffer)

	// Put the file object to s3 with the file name
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(bucketName),
		Key:                  aws.String(fileName),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(fileBuffer),
		ContentLength:        aws.Int64(fileSize),
		ContentType:          aws.String(http.DetectContentType(fileBuffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	fmt.Print("error here:", err)
	return err
}
