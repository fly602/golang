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

type Student struct {
	Name string
	Age  int
}

var (
	testDB             = "fly"
	testCollectionName = "uos"
)

// 基本功能测试
func dbTest() {
	ctx := context.Background()
	// 使用URI建立连接
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://admin:admin@127.0.0.1:27017/"))
	if err != nil {
		log.Panicln(err)
	}
	log.Println(client)

	// 关闭连接
	defer client.Disconnect(ctx)
	// ping测试连接是否可用
	fmt.Println(client.Ping(ctx, readpref.Primary()))

	// 输出database列表
	dbs, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Panicln(err)
	}
	log.Println("databases name=", dbs)
	for _, db := range dbs {
		collections, err := client.Database(testDB).ListCollectionNames(ctx, bson.M{})
		if err != nil {
			log.Panicln(err)
		}
		for _, col := range collections {
			log.Println(db, "has", col)
		}
	}

	s1 := Student{"小红", 12}
	s2 := Student{"小兰", 10}
	s3 := Student{"小黄", 11}

	var collection *mongo.Collection
	collection = client.Database(testDB).Collection(testCollectionName)
	// 如果不存在，则创建
	if collection == nil {
		err = client.Database(testDB).CreateCollection(ctx, testCollectionName)
		if err != nil {
			log.Panicln(err)
		}
		log.Println("collection uos created")
		collection = client.Database(testDB).Collection(testCollectionName)
	}

	// 创建索引
	indexModel := mongo.IndexModel{
		Keys: bson.D{{"age", -1}},
	}
	collection.Indexes().CreateOne(ctx, indexModel)

	// 添加一条记录
	res1, err := collection.InsertOne(ctx, s1)
	if err != nil {
		log.Panicln(err)
	}
	log.Println("Inserted a single document:", res1.InsertedID)

	// 添加多条记录
	students := []interface{}{s2, s3}
	res2, err := collection.InsertMany(ctx, students)
	if err != nil {
		log.Panicln(err)
	}
	log.Println("Inserted a single document:", res2.InsertedIDs)

	// 查询一条记录
	filter := bson.D{
		{"$or",
			bson.A{
				bson.D{{"name", "小兰"}},
				bson.D{{"name", "小黄"}},
			},
		},
		// {
		// 	"$and",
		// 	bson.A{
		// 		bson.D{{"age", 11}},
		// 	},
		// },
	}
	res3 := collection.FindOne(ctx, filter)
	if res3 == nil {
		log.Println("name = 小兰 not found")
	} else {
		resm := bson.M{}
		res3.Decode(&resm)
		log.Println(resm)
	}

	// 查询多个
	// 将选项传递给Find()
	findOptions := options.Find()
	// 对查询结果进行排序
	// findOptions.SetSort(bson.D{{"name", 1}, {"age", 1}})
	// 设置查询结果条数
	findOptions.SetLimit(3)

	cur, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		log.Panic("find err")
	}
	for cur.Next(ctx) {
		// 创建一个值，将单个文档解码为该值
		var elem Student
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("find multiple documents", elem)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// 完成后关闭游标
	cur.Close(ctx)

	// 修改文档
	updateFilter := bson.D{
		{
			"$or", bson.A{
				bson.D{{"name", "小兰"}},
				bson.D{{"name", "小黄"}},
			},
		},
	}

	updater := bson.D{
		{
			"$set", bson.D{{"age", 15}},
		},
	}
	updateRes, err := collection.UpdateMany(ctx, updateFilter, updater)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("update success", updateRes)

	// 删除文档和索引
	collection.DeleteMany(ctx, bson.D{{}})
	collection.Indexes().DropAll(ctx)

	// 删除集合
	err = collection.Drop(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// 删除数据库
	err = client.Database(testDB).Drop(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
