package routes

import (
	"aws-storage/middleware"
	"fmt"
	"log"
	"mime"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	"github.com/rhnvrm/simples3"
)

var AccessKeyID string
var SecretAccessKey string
var MyRegion string
var MyBucket string
var filePath string

func UploadImage(c *gin.Context) {
	middleware.LoadEnv()
	sess := c.MustGet("sess").(*session.Session)
	uploader := s3manager.NewUploader(sess)

	MyRegion = middleware.GetEnvWithKey("AWS_REGION")
	file, header, err := c.Request.FormFile("photo")
	if err != nil {
		log.Fatal(err)
	}
	filename := header.Filename
	up, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:               aws.String(MyBucket),
		ACL:                  aws.String("private"),
		Key:                  aws.String(filename),
		Body:                 file,
		ServerSideEncryption: aws.String("AES256"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":    err,
			"uploader": up,
		})
		return
	}
	filePath = "https://" + MyBucket + "." + "s3-" + MyRegion + ".amazonaws.com/" + filename
	c.JSON(http.StatusOK, gin.H{
		"filepath": filePath,
	})
}
func UploadTest(c *gin.Context) {
	middleware.LoadEnv()
	MyBucket = middleware.GetEnvWithKey("BUCKET_NAME")
	AccessKeyID = middleware.GetEnvWithKey("AWS_ACCESS_KEY_ID")
	SecretAccessKey = middleware.GetEnvWithKey("AWS_SECRET_ACCESS_KEY")
	MyRegion = middleware.GetEnvWithKey("AWS_REGION")

	file, header, err := c.Request.FormFile("photo")
	if err != nil {
		log.Fatal(err)
	}
	filename := header.Filename
	mediaType, _, err := mime.ParseMediaType(header.Header.Get("Content-Type"))
	if err != nil {
		log.Fatal(err)
	}

	s := simples3.New(MyRegion, AccessKeyID, SecretAccessKey)
	_, err = s.FileUpload(simples3.UploadInput{
		Bucket:      MyBucket,
		ObjectKey:   filename,
		ContentType: mediaType,
		FileName:    filename,
		Body:        file,
	})
	if err != nil {
		log.Fatalf("Error: %+v", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Uploaded!",
	})
	//sess := c.MustGet("sess").(*session.Session)
}
func ListImage(c *gin.Context) {
	middleware.LoadEnv()
	currentTime := time.Now().String()
	var time, _ = time.Parse(time.RFC1123, currentTime)
	MyBucket = middleware.GetEnvWithKey("BUCKET_NAME")
	AccessKeyID = middleware.GetEnvWithKey("AWS_ACCESS_KEY_ID")
	SecretAccessKey = middleware.GetEnvWithKey("AWS_SECRET_ACCESS_KEY")
	MyRegion = middleware.GetEnvWithKey("AWS_REGION")
	s := simples3.New(MyRegion, AccessKeyID, SecretAccessKey)
	url := s.GeneratePresignedURL(simples3.PresignedInput{
		Bucket:        MyBucket,
		ObjectKey:     "1ad56663af6534d929ad85d1f030e907.jpg",
		Method:        "GET",
		Timestamp:     time,
		ExpirySeconds: 86400,
	})
	fmt.Println("signed url:", url)
	c.HTML(http.StatusOK, "templates/index.tmpl", gin.H{
		"url": "1ad56663af6534d929ad85d1f030e907.jpg",
	})
	//sess := c.MustGet("sess").(*session.Session)
	//c.JSON(http.StatusOK, result)
	/* c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title": "Main website",
		"image": input.Key,
	}) */
}
func ListImages(c *gin.Context) {
	middleware.LoadEnv()
	MyBucket = middleware.GetEnvWithKey("BUCKET_NAME")
	var images []string
	sess := c.MustGet("sess").(*session.Session)
	svc := s3.New(sess)
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(MyBucket)})
	if err != nil {
		log.Fatalf("error listing objects: %+v", err)
	}
	for _, item := range resp.Contents {
		fmt.Println("Name: ", *item.Key)
		fmt.Println("Item: ", *item)
		images = append(images, *item.Key)
	}
	c.JSON(http.StatusOK, images)
}
