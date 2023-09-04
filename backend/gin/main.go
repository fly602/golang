package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type User struct {
	Name string
	Age  int
}

func handleHello(c *gin.Context) {
	users := []User{
		{"aaa", 1}, {"bbb", 2},
	}
	cookie, err := c.Cookie("key_cookie")
	if err != nil {
		cookie = "NotSet"
		c.SetCookie("key_cookie", "value_cookie", -1, "/", "localhost", false, true)
	}
	log.Printf("cookie的值是： %s\n", cookie)
	//c.String(http.StatusOK, "hello world")
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "你好，世界！",
		"body":  "你好，世界！！",
		"User":  users,
	})
}

func handleUser(c *gin.Context) {

	c.String(http.StatusOK, "hello User")
}

func handleCookie(c *gin.Context) {
	maxage := 60 * 60 * 24
	fmt.Println("fullpath=", c.FullPath())
	c.SetCookie("user_name", "abc", maxage, c.FullPath(), "localhost", false, false)
	c.SetCookie("mobile", "1234567890", maxage, c.FullPath(), "localhost", false, false)

	c.JSON(http.StatusOK, "登录成功")
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if cookie, err := c.Cookie("user_name"); err == nil {
			if cookie == "abc" {
				c.Next()
				return
			} else if cookie == "" {
				fmt.Println("need auth...")
				c.Next()
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err,
			})
			c.Abort()
		}
	}

}

func handleSession(c *gin.Context) {
	action := c.Query("action")
	session := sessions.Default(c)
	maxage := 60 * 60 * 24
	session.Options(sessions.Options{
		Path:     c.FullPath(),
		Domain:   "localhost",
		MaxAge:   maxage,
		Secure:   false,
		HttpOnly: false,
	})
	fmt.Println("action=", action)
	if action == "add" || action == "" {
		if session.Get("hello") != "world" {
			fmt.Println("set session...")
			session.Set("hello", "world")
			session.Set("name", "fuleyi")
			session.Save()
		}
	} else if action == "del" {
		fmt.Println("delete session")
		session.Delete("hello")
		session.Save()
	}

	c.JSON(http.StatusOK, gin.H{
		"hello": session.Get("hello"),
	})
}

func handleV2(c *gin.Context) {
	c.HTML(http.StatusOK, "v2.html", nil)
}

func main() {
	r := gin.Default()
	r.Use(gin.BasicAuth(gin.Accounts{
		"admin": "123456",
	}))
	r.LoadHTMLGlob("./temp/*")
	r.GET("/", handleHello)
	g1 := r.Group("/v1")
	g1.GET("/user", handleUser)
	g1.GET("/cookie", AuthMiddleware(), handleCookie)

	g2 := r.Group("v2")
	// g2.Use(middlewares.Cors())
	store := cookie.NewStore([]byte("scret11111"))
	g2.Use(sessions.Sessions("mysessions", store))
	g2.GET("/", handleV2)
	g2.GET("/session", handleSession)

	r.Run(":50002")

}
