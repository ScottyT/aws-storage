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
	})
	router.Use(middleware.Authenticate)
	router.POST("/upload", routes.UploadTest)
	router.LoadHTMLGlob("templates/*")
	router.GET("/images", routes.ListImages)
	router.GET("/image", routes.ListImage)
	router.POST("/download", gin.WrapF(routes.DownloadObjects))

	_ = router.Run(":8090")
}
