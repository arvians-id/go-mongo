package main

import (
	"github.com/arvians-id/go-mongo/adapter/cmd/config"
	"github.com/arvians-id/go-mongo/adapter/middleware"
	controller "github.com/arvians-id/go-mongo/adapter/pkg/user"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"path/filepath"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln("Failed at getting working directory", err)
	}

	dir := filepath.Join(wd, ".env")
	configuration := config.New(dir)

	router := gin.Default()

	router.Use(middleware.SetupCorsMiddleware())
	router.Use(middleware.GinContextToContextMiddleware())

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to My Project Go MongoDB. Created By https://github.com/arvians-id",
		})
	})

	controller.NewUserController(router, configuration)

	port := ":" + configuration.Get("APP_PORT")
	err = router.Run(port)
	if err != nil {
		log.Fatalln(err)
	}
}
