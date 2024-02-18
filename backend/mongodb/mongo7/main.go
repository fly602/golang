package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type User struct {
	Name    string `bson:"name"`
	Age     int    `bson:"age"`
	Address string `bson:address`
}

var users []User

func main() {
	ctx := context.Background()
	// 使用URI建立连接
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Panicln(err)
	}
	// 插入一条
	cursor, err := client.Database("user").Collection("users").Find(ctx, bson.D{{"name", "fly"}, {"age", 18}})
	if err != nil {
		log.Panicln(err)
	}
	if err := cursor.All(ctx, &users); err != nil {
		log.Panicln(err)
	}
	if len(users) == 0 {
		client.Database("user").Collection("users").InsertOne(ctx, bson.M{"name": "fly", "age": 18})
	} else {
		log.Println("users:", users)
	}
	// 关闭连接
	defer client.Disconnect(ctx)
	// ping测试连接是否可用
	fmt.Println(client.Ping(ctx, readpref.Primary()))
}
