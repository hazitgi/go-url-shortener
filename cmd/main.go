package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
	database "github.com/hazi-tgi/go-url-shortner/config"
	"github.com/hazi-tgi/go-url-shortner/controllers"
	"github.com/hazi-tgi/go-url-shortner/handlers"
)

var httpAddress string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading.env file")
	}
	flag.StringVar(&httpAddress, "http-address", ":8080", "HTTP listen address")
	database.Connect()
}

func main() {
	r := gin.Default()
	flag.Parse()

	gin.SetMode(gin.DebugMode)
	r.LoadHTMLGlob("static/*")

	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Home Page",
		})
	})

	urlController := controllers.NewUrlController(database.DB)

	UrlHandler := handlers.NewUrlHandler(urlController)

	UrlHandler.RegisterRoutes(r)

	fmt.Printf("Listening on http://%s\n", httpAddress)
	log.Fatal(r.Run(httpAddress))
}
