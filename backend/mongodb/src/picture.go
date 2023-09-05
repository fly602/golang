package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	pictureDB = "picture"
	// picturecollectionName = "uos"
	downloadPath = dirPath + "download/"
)

type fileDocInfo struct {
	FileName string
	Length   int64
}

func uploadPictureByPath(dirPath string) {
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

	fileInfos, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, fileInfo := range fileInfos {
		if filepath.Ext(fileInfo.Name()) == ".jpg" {
			filePath := filepath.Join(dirPath, fileInfo.Name())
			imageData, err := ioutil.ReadFile(filePath)
			if err != nil {
				log.Printf("Error reading file %s: %v\n", filePath, err)
				continue
			}

			bucket, _ := gridfs.NewBucket(client.Database(pictureDB))
			// 创建一个新的GridFS文件
			uploadStream, err := bucket.OpenUploadStream(
				fileInfo.Name(), // 替换为您希望存储的图片名称
			)
			if err != nil {
				log.Fatal(err)
			}
			defer uploadStream.Close()

			// 将图片数据写入GridFS文件
			_, err = uploadStream.Write(imageData)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Image saved to MongoDB!")
		}
	}
}

func downloadPicture() {
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

	bucket, _ := gridfs.NewBucket(client.Database(pictureDB))

	collection := bucket.GetFilesCollection()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("file collection name =", collection.Name())

	// 查询file collection中所有的图片信息
	cur, err := bucket.Find(bson.M{})
	if err != nil {
		log.Fatalln(err)
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var res bson.M
		err := cur.Decode(&res)
		if err != nil {
			log.Fatalln(err)
		}
		if _, err := os.Stat(downloadPath); err != nil {
			if os.IsNotExist(err) {
				os.Mkdir(downloadPath, 0755)
			}
		}

		fileName := filepath.Base(res["filename"].(string))
		fileExt := filepath.Ext(fileName)
		fileID := "_" + (res["_id"].(primitive.ObjectID)).Hex()
		fileFlag := ".download"
		newFileName := fileName[:len(fileName)-len(fileExt)] + fileID + fileFlag + fileExt

		file, err := os.Create(downloadPath + newFileName)
		if err != nil {
			log.Fatalln(err)
		}
		_, err = bucket.DownloadToStreamByName(res["filename"].(string), file)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("download file :", file.Name())
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
}
