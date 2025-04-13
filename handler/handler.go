package handler

import (
	"imgbed/config"
	"imgbed/controller"
	"imgbed/middleware"

	"github.com/gin-gonic/gin"
)

func SetupImgBedRoute(router *gin.Engine, c config.Config) {

	ImgBedGroup := router.Group("/")
	{
		ImgBedGroup.POST("upload", middleware.Auth(c), func(ctx *gin.Context) {
			controller.UploadImage(ctx, c)
		})
		ImgBedGroup.POST("delete", middleware.Auth(c), func(ctx *gin.Context) {
			controller.DeleteImage(ctx, c)
		})
		ImgBedGroup.POST("list", middleware.Auth(c), func(ctx *gin.Context) { controller.List(ctx, c) })
		ImgBedGroup.GET("i/*filePath", func(ctx *gin.Context) { controller.ShowImg(ctx, c) })
		ImgBedGroup.GET("thumbnail/*filePath", func(ctx *gin.Context) { controller.ShowThumbnailImg(ctx, c) })
	}
}
