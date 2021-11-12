package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var AccessKeyID string
var SecretAccessKey string
var MyRegion string
var MyBucket string

func GetEnvWithKey(key string) string {
	return os.Getenv(key)
}

func LoadEnv() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
		os.Exit(1)
	}
}

//Add firebase middleware here
func ConnectAws() *session.Session {
	AccessKeyID = GetEnvWithKey("AWS_ACCESS_KEY_ID")
	SecretAccessKey = GetEnvWithKey("AWS_SECRET_ACCESS_KEY")
	MyRegion = GetEnvWithKey("AWS_REGION")

	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(MyRegion),
			Credentials: credentials.NewStaticCredentials(
				AccessKeyID,
				SecretAccessKey,
				"",
			),
		})

	if err != nil {
		panic(err)
	}

	return sess
}
func Authenticate(c *gin.Context) {
	LoadEnv()
	AccessKeyID = GetEnvWithKey("AWS_ACCESS_KEY_ID")
	SecretAccessKey = GetEnvWithKey("AWS_SECRET_ACCESS_KEY")

}
func UpdateBucket(c *gin.Context) {
	//permission := "READ"
	LoadEnv()
	MyBucket = GetEnvWithKey("BUCKET_NAME")

	/* if auth == "" {
		return
	} */

	sess := c.MustGet("sess").(*session.Session)
	// Create S3 service client
	svc := s3.New(sess)
	readOnly := map[string]interface{}{
		"Version": "2012-10-17",
		"Statement": []map[string]interface{}{
			{
				"Sid":       "AddPerm",
				"Effect":    "Allow",
				"Principal": "*",
				"Action": []string{
					"s3:GetObject",
				},
				"Resource": []string{
					fmt.Sprintf("arn:aws:s3:::%s/*", MyBucket),
				},
			},
		},
	}
	policy, err := json.Marshal(readOnly)
	if err != nil {
		log.Fatal("Failed to markshal policy")
	}
	_, err = svc.PutBucketPolicy(&s3.PutBucketPolicyInput{
		Bucket: aws.String(MyBucket),
		Policy: aws.String(string(policy)),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == s3.ErrCodeNoSuchBucket {
			// Special error handling for the when the bucket doesn't
			// exists so we can give a more direct error message from the CLI.
			log.Fatalf("Bucket %q does not exist", MyBucket)
		}
		log.Fatalf("Unable to set bucket %q policy, %v", MyBucket, err)
	}
	c.Next()
}
