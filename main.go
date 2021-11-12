package main

import (
	"aws-storage/middleware"
	"aws-storage/routes"

	"github.com/gin-gonic/gin"
)

var MyBucket string

func main() {
	middleware.LoadEnv()
	MyBucket = middleware.GetEnvWithKey("BUCKET_NAME")
	sess := middleware.ConnectAws()

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("sess", sess)
		c.Next()
	})
	router.POST("/upload", routes.UploadImage)
	router.LoadHTMLGlob("templates/*")
	router.GET("/images", routes.ListImages)
	router.GET("/image", routes.ListImage)

	_ = router.Run(":8082")
}
