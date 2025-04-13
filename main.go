package main

import (
	"fmt"
	"imgbed/config"
	"imgbed/handler"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	//初始化配置
	var c config.Config
	err := c.Load()
	if err != nil {
		log.Fatalf("loal config fail\n%s", err)
	}
	log.Printf("%+v", c)
	//创建三个图片文件夹
	initDir(c)
	if err != nil {
		log.Fatalf("create img dir fail\n%s", err)
	}
	//启动服务
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     c.Server.Cors.Origins,
		AllowMethods:     c.Server.Cors.Methods,
		AllowHeaders:     c.Server.Cors.Headers,
		AllowCredentials: c.Server.Cors.Credentials,
	}))

	handler.SetupImgBedRoute(r, c)
	http.ListenAndServe(fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port), r)
}
func initDir(c config.Config) error {
	threeDir := []string{c.Server.OriginalDir, c.Server.ThumbnailDir, c.Server.PublicDir}
	for _, v := range threeDir {
		err := os.Mkdir(v, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
