package main

import (
	"github.com/gin-gonic/contrib/static"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// r := router.Group("/v1")

	// r.GET("/ping", testGET)

	// r.POST("/post", func(c *gin.Context) {
	// 	id := c.Query("id")
	// 	page := c.DefaultQuery("page", "0")
	// 	name := c.PostForm("name")
	// 	message := c.PostForm("message")

	// 	fmt.Printf("id: %s; page: %s; name: %s; message: %s", id, page, name, message)
	// })

	// r.Use(static.Serve("/", BinaryFileSystem("assets")))

	// router.GET("/", testGET)

	// router.StaticFS("/", http.Dir("templates/build"))

	// Serve frontend static files
	router.Use(static.Serve("/", static.LocalFile("./templates/build", true)))
	router.Use(static.Serve("/404", static.LocalFile("./templates/build", true)))

	// Serve frontend static files
	// router.Use(router.Serve("/", router.LocalFile("templates/build", true)))

	router.Run(":8080") // listen and serve on 0.0.0.0:8080
}

func testGET(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

// func BinaryFileSystem(root string) *binaryFileSystem {
// 	fs := &assetfs.AssetFS{Asset, AssetDir, AssetInfo, root}
// 	return &binaryFileSystem{
// 		fs,
// 	}
// }
