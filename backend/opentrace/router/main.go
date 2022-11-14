package main

import (
	"go-community/backend/opentrace/proto/product"

	"github.com/gin-gonic/gin"
)

var products product.AllResponse

func main() {
	r := gin.Default()
	r.GET("")
}
