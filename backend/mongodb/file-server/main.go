package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fly602/golang/backend/mongodb/file-server/fileserver"
	"github.com/gin-gonic/gin"
)

var conn *fileserver.Conn

var (
	uri    = "mongodb://admin:admin@127.0.0.1:27017/"
	dbName = "picture"
)

func HandlePicture(c *gin.Context) {
	files := conn.ListFile()
	log.Println("ListFile:", files)
	// 构建文件链接
	var fileLinks []string
	baseURL := "http://127.0.0.1:50001/download/" // 替换为您的文件基础URL

	for _, fileName := range files {
		fileLink := fmt.Sprintf("%s%s", baseURL, fileName)
		fileLinks = append(fileLinks, fileLink)
	}

	// 返回文件链接到前端
	c.JSON(http.StatusOK, gin.H{"fileLinks": fileLinks})
}

func handleHome(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func handleUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.String(500, "上传图片出错")
	}
	// c.JSON(200, gin.H{"message": file.Header.Context})
	log.Println("upload", file.Filename)
	src, err := file.Open()
	if err != nil {
		c.String(http.StatusInternalServerError, "上传文件失败")
		return
	}
	defer src.Close()
	newFileName := conn.UploadFile(file.Filename, src)
	if newFileName == "" {
		c.String(http.StatusInternalServerError, "上传文件失败")
		return
	}
	c.String(http.StatusOK, newFileName)

}

func handleDownload(c *gin.Context) {
	filename := c.Param("filename")
	err := conn.DownloadFile(filename, c.Writer)
	if err != nil {
		c.String(http.StatusInternalServerError, "下载文件失败")
		return
	}
	c.String(http.StatusOK, "download OK")
}

func handleDelete(c *gin.Context) {
	filename := c.Param("filename")
	err := conn.DeleteFile(filename)
	if err != nil {
		c.String(http.StatusInternalServerError, "删除文件失败")
		return
	}
	c.String(http.StatusOK, "delete OK")
}

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime)
	conn = fileserver.Connect(uri, dbName)
	r := gin.Default()
	r.LoadHTMLFiles("html/index.html")
	r.GET("/", handleHome)
	r.POST("/upload", handleUpload)
	r.GET("/download/:filename", handleDownload)
	r.GET("/picture", HandlePicture)
	r.GET("/delete/:filename", handleDelete)
	r.Run(":50001")
}
