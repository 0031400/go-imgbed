package main

import (
	"fmt"
	"imgbed/config"
	"imgbed/handler"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	var c config.Config
	err := c.Load()
	if err != nil {
		log.Fatalf("loal config fail\n%s\n", err)
	}
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	handler.SetupImgBedRoute(r, c)
	http.ListenAndServe(fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port), r)
}
