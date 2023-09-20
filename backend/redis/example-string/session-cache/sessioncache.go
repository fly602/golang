package sessioncache

import (
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
)

func Set() {
	store, err := redis.NewStore(10, "tcp", "localhost:56379", "", []byte(""))
	if err != nil {
		log.Fatal("NewStore err")
	}
	sessions.Sessions("mysession", store)
}
