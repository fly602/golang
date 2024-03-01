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
	Name    string      `bson:"name"`
	Age     int         `bson:"age"`
	Address string      `bson:"address"`
	Type    string      `bson:"type"`
	Data    interface{} `bson:"data"`
	Version string      `bson:"version"`
}

var users []User

var client *mongo.Client

var tUser *mongo.Collection

func getUserCollection(ctx context.Context) *mongo.Collection {
	var err error
	if client == nil {
		client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
		if err != nil {
			log.Panicln(err)
		}
	}
	if tUser == nil {
		tUser = client.Database("user").Collection("users")
	}
	return tUser
}

func findIndex(cursor *mongo.Cursor, name string) bool {
	for {
		if cursor.TryNext(context.Background()) {
			var index bson.M
			if err := cursor.Decode(&index); err != nil {
				log.Panicln(err)
			}
			if index["name"] == name {
				log.Println("index found:", index)
				return true
			}
			log.Println("get index:", index)
		} else {
			break
		}
		if err := cursor.Err(); err != nil {
			log.Panicln(err)
		}
	}
	return false
}

func main() {
	ctx := context.Background()
	getUserCollection(ctx)
	tUser := client.Database("user").Collection("users")
	// 修改版本号字段
	_, _ = tUser.UpdateMany(ctx, bson.D{
		{
			"data", bson.D{
				{"$exists", true},
			},
		},
	}, bson.D{{
		"$set", bson.D{{"version", "1.0.3"}},
	}})
	// 修改版本号字段
	_, _ = tUser.UpdateMany(ctx, bson.D{
		{
			"data", bson.D{
				{"$exists", false},
			},
		},
		{
			"address", bson.D{
				{"$exists", true},
			},
		},
	}, bson.D{{
		"$set", bson.D{{"version", "1.0.2"}},
	}})

	// 修改版本号字段
	_, _ = tUser.UpdateMany(ctx, bson.D{
		{
			"address", bson.D{
				{"$exists", false},
			},
		},
		{
			"data", bson.D{
				{"$exists", false},
			},
		},
	}, bson.D{{
		"$set", bson.D{{"version", "1.0.1"}},
	}})

	cursor, err := tUser.Find(ctx, bson.D{{"name", "fly"}, {"age", 18}})
	if err != nil {
		log.Panicln(err)
	}
	if err := cursor.All(ctx, &users); err != nil {
		log.Panicln(err)
	}
	if len(users) == 0 {
		one, err := tUser.InsertOne(ctx, bson.M{
			"name":    "fly",
			"age":     18,
			"address": "hubei",
			"type":    "mixed",
			"data":    bson.A{"bar", "world", 3.14159, bson.D{{"qux", 12345}}},
		})
		if err != nil {
			log.Panicln(err)
		}
		log.Println("InsertOne success", one.InsertedID)
	} else {
		// 修改name为fly 且没有设置address字段的文档
		one, err := tUser.UpdateOne(ctx,
			bson.D{
				{"name", "fly"},
				{
					"address", bson.D{
						{
							"$exists", false,
						},
					},
				},
			},
			bson.D{
				{"$set", bson.M{
					"address": "hubei",
				}},
			}, nil)
		if err != nil {
			log.Panicln(err)
		}
		log.Println("Update success", one.UpsertedID)
	}

	// 尝试插入一条mixed数据
	cursor, err = tUser.Find(ctx, bson.M{
		"name":    "aaa",
		"age":     16,
		"address": "hubei",
		"type":    "mixed",
	})

	if err != nil {
		log.Panicln(err)
	}

	if err := cursor.All(ctx, &users); err != nil {
		log.Panicln(err)
	}

	log.Println(users)
	if len(users) == 0 {
		one, err := tUser.InsertOne(ctx, bson.M{
			"name":    "aaa",
			"age":     16,
			"address": "hubei",
			"type":    "mixed",
			"data":    bson.A{"bar", "world", 3.14159, bson.D{{"qux", 12345}}},
		})
		if err != nil {
			log.Panicln(err)
		}
		log.Println("InsertOne success", one.InsertedID)
	}

	_, _ = tUser.UpdateMany(ctx, bson.D{
		{
			"$or", bson.A{
				bson.D{
					{
						"redundantField", bson.D{
							{
								"$exists", false,
							},
						},
					},
				},
				bson.D{
					{
						"redundantField", "",
					},
				},
			},
		},
	}, bson.D{
		{
			"$set", bson.M{"redundantField": "This field will deleted"},
		},
	})
	//_, _ = tUser.UpdateMany(ctx, bson.D{}, bson.D{
	//	{
	//		"$unset", bson.M{"redundantField": ""},
	//	},
	//})

	// 创建索引
	c_index, err := tUser.Indexes().List(ctx, nil)
	if err != nil {
		log.Panicln(err)
	}
	defer c_index.Close(ctx)
	if !findIndex(c_index, "name_1_age_1") {
		key, err := tUser.Indexes().CreateOne(ctx, mongo.IndexModel{
			Keys: bson.D{{"name", 1}, {"age", 1}},
		}, nil)
		if err != nil {
			log.Panicln(err)
		}
		log.Println("create index:", key)
	}

	if findIndex(c_index, "name_1") {
		_, err := tUser.Indexes().DropOne(ctx, "name_1", nil)
		if err != nil {
			log.Panicln(err)
		}
		log.Println("drop index:", "name_1")
	}
	// Retrieves and prints the number of documents in the collection
	// that match the filter
	count, err := tUser.EstimatedDocumentCount(ctx, nil)
	if err != nil {
		log.Panicln(err)
	}
	log.Println("count documents:", count)
	// 关闭连接
	defer log.Println("Disconnect MongoDB.", client.Disconnect(ctx))
	// ping测试连接是否可用
	fmt.Println(client.Ping(ctx, readpref.Primary()))
}
