package main

import (
	"log"
	"net/http"

	webcache "github.com/fly602/golang/backend/redis/example-string/web-cache"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

var redisClient *redis.Client

func initRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:56379", // Redis 服务器地址
		Password: "",                // Redis 密码，如果有的话
		DB:       0,                 // Redis 数据库编号
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	return client
}

func init() {
	redisClient = initRedisClient()
	log.SetFlags(log.Ltime | log.Lshortfile)
}

func handleIndex(c *gin.Context) {
	data, err := webcache.GetWebPageContent(redisClient, "html/index.html")
	if err != nil {
		c.JSON(http.StatusFound, "页面不存在")
	}
	c.Writer.Write(data)
}

func main() {
	r := gin.Default()
	r.GET("/", handleIndex)

	r.Run(":50001")

	err := redisClient.Close()
	if err != nil {
		log.Fatal("close redis client err:", err)
	}
}
