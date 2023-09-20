package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	webcache "github.com/fly602/golang/backend/redis/example-string/web-cache"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

type usrInfo struct {
	Name string
	Tel  string
	Pwd  string
}

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
	webcache.CleanWebCache(redisClient)
}

func handleIndex(c *gin.Context) {
	data, err := webcache.GetWebPageContent(redisClient, "html/index.html")
	if err != nil {
		c.String(http.StatusOK, "页面不存在")
	}
	c.Writer.Write(data)
}

func handleLogin(c *gin.Context) {
	c.Request.ParseForm()
	form := c.Request.Form

	// 查询缓存是否存在
	key := fmt.Sprintf("userinfo:%v", form["username"][0])
	val, err := redisClient.Get(context.Background(), key).Result()
	if err != nil {
		log.Println("passwd err")
		c.String(http.StatusOK, "用户不存在")
		return
	}
	user := &usrInfo{}
	err = json.Unmarshal([]byte(val), user)
	if err != nil {
		c.String(http.StatusOK, "数据解析失败")
		return
	}
	pwd := fmt.Sprintf("%x", md5.Sum([]byte(form["pwd"][0])))
	if pwd != user.Pwd {
		log.Println("passwd err")
		c.String(http.StatusOK, "密码错误")
		return
	}

	// pwd := form["password"]
	tmp, err := template.ParseFiles("html/home.tmpl")
	if err != nil {
		c.String(http.StatusOK, "数据解析失败")
		return
	}
	tmp.Execute(c.Writer, user.Name)
}

func handleRegister(c *gin.Context) {
	user := &usrInfo{}
	c.Request.ParseForm()
	form := c.Request.Form
	user.Name = form["username"][0]
	user.Tel = form["tel"][0]
	user.Pwd = fmt.Sprintf("%x", md5.Sum([]byte(form["pwd"][0])))

	value, err := json.Marshal(user)
	if err != nil {
		c.String(http.StatusOK, "数据解析失败")
		return
	}

	// 注册用户存储数据库
	key := fmt.Sprintf("userinfo:%v", user.Name)
	err = redisClient.SetNX(context.Background(), key, value, 0).Err()
	if err != nil {
		c.String(http.StatusOK, "服务器处理失败")
		return
	}
	tmp, err := template.ParseFiles("html/home.tmpl")
	if err != nil {
		c.String(http.StatusOK, "数据解析失败")
		return
	}
	tmp.Execute(c.Writer, user.Name)
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("html/*")
	r.GET("/", handleIndex)
	r.POST("/login", handleLogin)
	r.GET("/login", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/register", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "register.html", nil)
	})
	r.POST("/register", handleRegister)
	r.Run(":50001")
	err := redisClient.Close()
	if err != nil {
		log.Fatal("close redis client err:", err)
	}
}
