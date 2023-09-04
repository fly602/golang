package main

import "github.com/gin-gonic/gin"

func HandleIndex(c *gin.Context) {

}

func main() {
	r := gin.Default()
	r.GET("/", HandleIndex)

	r.Run(":50002")
}
