package main

import (
	"log"

	"github.com/garyburd/redigo/redis"
)

func MutiSet(c redis.Conn) {
	_, err := c.Do("Mset", "user:fly:age", 30, "user:cc:age", 28, "user:fsz:age", 2)
	if err != nil {
		log.Fatalln("Mset err", err)
	}
}

func MutiGet(c redis.Conn) {
	rep, err := redis.Ints(c.Do("Mget", "user:fly:age", "user:cc:age", "user:fsz:age"))
	if err != nil {
		log.Fatalln("Mset err", err)
	}
	for _, v := range rep {
		log.Println("Mget", v)

	}
}

func ExpireSet(c redis.Conn) {
	_, err := c.Do("expire", "user:fly:age", 1)
	if err != nil {
		log.Fatalln("redis SetExpire err", err)
	}
}

func ListSet(c redis.Conn) {
	for i := 0; i < 10000; i++ {
		_, err := c.Do("lpush", "number", i)
		if err != nil {
			log.Fatalln("redis lpush err", err)
		}
	}
}

func ListPop(c redis.Conn) {
	for i := 0; i < 10000; i++ {
		rep, err := redis.Int(c.Do("lpop", "number"))
		if err != nil {
			log.Fatalln("redis Dail err", err)
		}
		log.Println("Lpop", rep)
	}
}

func HashSet(c redis.Conn) {
	_, err := c.Do("HMset", "user", "name", "fly", "age", 30, "brith", "1992-06-12")
	if err != nil {
		log.Fatalln("redis Dail err", err)
	}
}

func main() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.Ldate | log.Ltime)
	conn, err := redis.Dial("tcp", "127.0.0.1:56379")
	if err != nil {
		log.Fatalln("redis Dail err", err)
	}
	defer conn.Close()
	_, err = conn.Do("set", "user:fly:age", 30)
	if err != nil {
		log.Fatal("Set err", err)
	}
	r, err := redis.Int(conn.Do("get", "user:fly:age"))
	if err != nil {
		log.Fatal("Set err", err)
	}
	log.Println("user:fly:age=", r)
	MutiSet(conn)
	MutiGet(conn)
	ExpireSet(conn)
	ListSet(conn)
	ListPop(conn)
	HashSet(conn)
}
