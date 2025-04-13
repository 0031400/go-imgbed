package controller

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"imgbed/config"
	"imgbed/util"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/golang/freetype"
)

type FileInfo struct {
	Name  string `json:"name"`
	IsDir bool   `json:"isDir"`
}

func InternalError(ctx *gin.Context, data string) {
	ctx.JSON(http.StatusInternalServerError, gin.H{"code": 1, "data": data})
}
func UploadImage(ctx *gin.Context, c config.Config) {

	file, err := ctx.FormFile("file")
	if err != nil {
		InternalError(ctx, "fail to get uploaded file")
		log.Println(err)
		return
	}
	ext := filepath.Ext(file.Filename)
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" && ext != ".webp" {
		InternalError(ctx, "not allowed image type")
		log.Println(err)
		return
	}
	timePath := time.Now().Format("2006/01/02")
	originalDirPath := filepath.Join(c.Server.OriginalDir, timePath)
	if err := os.MkdirAll(originalDirPath, 0755); err != nil {
		InternalError(ctx, "Failed to create original directory")
		log.Println(err)
		return
	}
	randomString := util.GenerateRandomString(c.Server.RandomNum)
	newOriginalFileName := randomString + ext
	originalFilePath := filepath.Join(originalDirPath, newOriginalFileName)
	ctx.SaveUploadedFile(file, originalFilePath)
	src, err := os.Open(originalFilePath)
	if err != nil {

		InternalError(ctx, "fail to open uploaded file")
		log.Println(err)
		return
	}
	defer src.Close()
	var img image.Image
	switch ext {
	case ".png":
		img, err = png.Decode(src)
	case ".jpeg":
		img, err = jpeg.Decode(src)
	case ".jpg":
		img, err = jpeg.Decode(src)
	case ".webp":
		img, err = webp.Decode(src)
	default:
		InternalError(ctx, "Failed to decode image")
		log.Println(err)
		return
	}
	if err != nil {
		InternalError(ctx, "Failed to decode image")
		log.Println(err)
		return
	}
	//生成缩略图
	thumbnailImg := imaging.Resize(img, c.Thumbnail.Width, c.Thumbnail.Height, imaging.Lanczos)
	thumbnailDir := filepath.Join(c.Server.ThumbnailDir, timePath)
	if err = os.MkdirAll(thumbnailDir, 0755); err != nil {
		InternalError(ctx, "fail to create thumbnail dir")
		log.Println(err)
	}
	err = imaging.Save(thumbnailImg, filepath.Join(thumbnailDir, newOriginalFileName))
	if err != nil {
		InternalError(ctx, "Failed to save thumbnail file")
		log.Println(err)
		return
	}
	//生成webp图
	webpDirPath := filepath.Join(c.Server.PublicDir, timePath)
	if err := os.MkdirAll(webpDirPath, 0755); err != nil {
		InternalError(ctx, "Failed to create webp directory")
		log.Println(err)
		return
	}
	webpFileName := randomString + ".webp"
	webpFilePath := filepath.Join(webpDirPath, webpFileName)
	webpFile, err := os.Create(webpFilePath)
	if err != nil {

		InternalError(ctx, "Failed to create WebP file")
		log.Println(err)
		return
	}
	defer webpFile.Close()
	webpImg := img
	//水印处理
	if c.WaterMark.Enable {
		NewDrawImg := image.NewRGBA(img.Bounds())
		draw.Draw(NewDrawImg, img.Bounds(), img, image.Point{0, 0}, draw.Src)
		fontBytes, err := os.ReadFile(c.WaterMark.Font)
		if err != nil {
			InternalError(ctx, "Failed to read font file")
			log.Println(err)
			return
		}
		font, err := freetype.ParseFont(fontBytes)
		if err != nil {
			InternalError(ctx, "Failed to parse font file")
			log.Println(err)
			return
		}
		f := freetype.NewContext()
		f.SetDPI(float64(c.WaterMark.Pdi))
		f.SetFont(font)
		f.SetFontSize(float64(c.WaterMark.Size))
		f.SetClip(img.Bounds())
		f.SetDst(NewDrawImg)
		f.SetSrc(image.NewUniform(color.RGBA{R: uint8(c.WaterMark.Color.R), G: uint8(c.WaterMark.Color.G), B: uint8(c.WaterMark.Color.B), A: uint8(c.WaterMark.Color.A)}))
		pt := freetype.Pt(c.WaterMark.Position.X, c.WaterMark.Position.Y+c.WaterMark.Size)
		_, err = f.DrawString(c.WaterMark.Text, pt)
		if err != nil {
			InternalError(ctx, "Failed to draw watermark")
			log.Println(err)
			return
		}
		webpImg = NewDrawImg
	}
	err = webp.Encode(webpFile, webpImg, &webp.Options{Quality: float32(c.Quality)})
	if err != nil {
		InternalError(ctx, "Failed to encode image to WebP")
		log.Println(err)
		return
	}

	imgPath := fmt.Sprintf("%s/%s", timePath, webpFileName)
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": imgPath,
	})
}
func DeleteImage(ctx *gin.Context, c config.Config) {
	var filePathList []string
	if err := ctx.ShouldBindBodyWithJSON(&filePathList); err != nil {
		InternalError(ctx, "fail to decode list")
		log.Println(err)
		return
	}
	filePathListLen := len(filePathList)
	requestOriginalFileName := filePathList[filePathListLen-1]
	id := strings.Split(requestOriginalFileName, ".")[0]

	prefixPath := filepath.Join(filePathList[:filePathListLen-1]...)
	webpFilePath := filepath.Join(c.Server.PublicDir, prefixPath, id+".webp")
	_, err := os.Stat(webpFilePath)
	if err != nil {
		InternalError(ctx, "fail to open webp image")
		log.Println(err)
		return
	}
	err = os.Remove(webpFilePath)
	if err != nil {
		InternalError(ctx, "fail to remove webp image")
		log.Println(err)
		return

	}
	originalFilePath := filepath.Join(c.Server.OriginalDir, prefixPath, requestOriginalFileName)
	_, err = os.Stat(originalFilePath)
	if err != nil {
		InternalError(ctx, "fail to open original image")
		log.Println(err)
		return
	}
	err = os.Remove(originalFilePath)
	if err != nil {

		InternalError(ctx, "fail to remove original file")
		log.Println(err)
		return
	}
	thumbnailFilePath := filepath.Join(c.Server.ThumbnailDir, prefixPath, requestOriginalFileName)
	_, err = os.Stat(thumbnailFilePath)
	if err != nil {
		InternalError(ctx, "fail to open thumbnail image")
		log.Println(err)
		return
	}
	err = os.Remove(thumbnailFilePath)
	if err != nil {

		InternalError(ctx, "fail to remove thumbnail file")
		log.Println(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
	})
}
func List(ctx *gin.Context, c config.Config) {
	var filePathList []string
	if err := ctx.ShouldBindBodyWithJSON(&filePathList); err != nil {
		InternalError(ctx, "fail to decode list")
		log.Println(err)
		return
	}

	originalDir := filepath.Join(append([]string{c.Server.OriginalDir}, filePathList...)...)
	_, err := os.Stat(originalDir)
	if err != nil {
		InternalError(ctx, "fail to open dir")
		log.Println(err)
		return
	}
	files, err := os.ReadDir(originalDir)
	if err != nil {
		InternalError(ctx, "fail to list file")
		log.Println(err)
		return

	}
	var fileList []FileInfo
	for _, v := range files {
		fileList = append(fileList, FileInfo{Name: v.Name(), IsDir: v.Type().IsDir()})
	}
	if len(fileList) == 0 {
		ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": []string{}})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": fileList})
}
func ShowImg(ctx *gin.Context, c config.Config) {
	filePath := ctx.Param("filePath")
	ctx.File(filepath.Join(c.Server.PublicDir, filePath))
}
func ShowThumbnailImg(ctx *gin.Context, c config.Config) {
	filePath := ctx.Param("filePath")
	ctx.File(filepath.Join(c.Server.ThumbnailDir, filePath))
}
