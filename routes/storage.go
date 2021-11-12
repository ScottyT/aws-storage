package routes

import (
	"aws-storage/middleware"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
)

var MyBucket string
var MyRegion string
var filepath string

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
	filepath = "https://" + MyBucket + "." + "s3-" + MyRegion + ".amazonaws.com/" + filename
	c.JSON(http.StatusOK, gin.H{
		"filepath": filepath,
	})
}
func ListImage(c *gin.Context) {
	middleware.LoadEnv()
	MyBucket = middleware.GetEnvWithKey("BUCKET_NAME")
	//sess := c.MustGet("sess").(*session.Session)
	//owner := c.MustGet("owner").(string)
	//owner := c.GetHeader("x-amz-expected-bucket-owner")
	//svc := s3.New(sess)
	input := &s3.GetObjectInput{
		Bucket: aws.String(MyBucket),
		Key:    aws.String("1ad56663af6534d929ad85d1f030e907.jpg"),
		//ExpectedBucketOwner: &owner,
	}
	/* result, err := svc.GetObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				fmt.Println(s3.ErrCodeNoSuchKey, aerr.Error())
			case s3.ErrCodeInvalidObjectState:
				fmt.Println(s3.ErrCodeInvalidObjectState, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	} */
	//c.JSON(http.StatusOK, result)
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title": "Main website",
		"image": input.Key,
	})
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

func CreateBucketACL(c *gin.Context) {
	/* sess := c.MustGet("sess").(*session.Session)
	svc := s3.New(sess) */
	user := c.PostForm("user")
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
	/* result, err := svc.GetBucketAcl(&s3.GetBucketAclInput{Bucket: &MyBucket})
	if err != nil {
		log.Fatal(err)
	} */
	/* fmt.Print(result)

	owner := *result.Owner.DisplayName
	ownerId := *result.Owner.ID

	// Existing grants
	grants := result.Grants

	// Create new grantee to add to grants
	var newGrantee = s3.Grantee{EmailAddress: &address, Type: &userType}
	var newGrant = s3.Grant{Grantee: &newGrantee, Permission: &permission} */
}
