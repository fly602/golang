package fileserver

import (
	"context"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Conn struct {
	client   *mongo.Client
	Ctx      context.Context
	database *mongo.Database
}

// Connect to mongoDB
func Connect(uri string, dbName string) *Conn {
	ctx := context.Background()
	// 使用URI建立连接
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Panicln(err)
	}
	database := client.Database(dbName)

	if database == nil {
		log.Panicln(err)
	}

	log.Println("Connect to mongoDB success")
	return &Conn{
		client:   client,
		Ctx:      ctx,
		database: database,
	}
}

func (c *Conn) DisConnect() {
	c.client.Disconnect(c.Ctx)
}

func (c *Conn) Database() *mongo.Database {
	if c == nil || c.client == nil {
		log.Panic("Conn nil")
	}
	return c.database
}

func (c *Conn) Client() *mongo.Client {
	return c.client
}

// 上传文件，缓存的文件名为：文件名_md5值_id.后缀
func (c *Conn) UploadFile(filename string, file multipart.File) string {
	begin := time.Now()
	bucket, err := gridfs.NewBucket(c.Database())
	if err != nil {
		return ""
	}
	// 将图片数据写入GridFS文件
	id := primitive.NewObjectID()
	fileExt := filepath.Ext(filename)
	fileName := filepath.Base(filename)
	fileID := "_" + md5sum(file) + "_" + id.Hex()
	newFileName := fileName[:len(fileName)-len(fileExt)] + fileID + fileExt

	// 重置文件读取位置
	file.Seek(0, 0)
	err = bucket.UploadFromStreamWithID(id, newFileName, file)
	if err != nil {
		return ""
	}
	log.Printf("上传文件: %s 耗时: %v\n", newFileName, time.Since(begin).Microseconds())

	return newFileName
}

func objectIDFromFileName(filename string) (primitive.ObjectID, error) {
	if filename == "" {
		err := errors.New("filename empty")
		return primitive.NilObjectID, err
	}
	fileExt := filepath.Ext(filename)
	fileName := filepath.Base(filename)
	realname := fileName[:len(fileName)-len(fileExt)]
	if realname == "" {
		err := errors.New("filename illegal")
		return primitive.NilObjectID, err
	}
	n := strings.LastIndex(realname, "_")
	var id string
	if n != -1 {
		id = realname[n+1:]
	} else {
		err := errors.New("filename illegal")
		return primitive.NilObjectID, err
	}
	return primitive.ObjectIDFromHex(id)
}

func (c *Conn) DownloadFile(filename string, w io.Writer) error {
	begin := time.Now()
	bucket, err := gridfs.NewBucket(c.Database())
	if err != nil {
		return err
	}
	id, err := objectIDFromFileName(filename)
	if err != nil {
		return err
	}

	_, err = bucket.DownloadToStream(id, w)
	if err != nil {
		return err
	}
	log.Printf("下载文件: %s 耗时: %v\n", filename, time.Since(begin).Microseconds())
	return nil
}

func (c *Conn) ListFile() []string {
	files := make([]string, 0)
	bucket, err := gridfs.NewBucket(c.Database())
	if err != nil {
		return nil
	}
	// 查询file collection中所有的图片信息
	cur, err := bucket.Find(bson.M{})
	if err != nil {
		log.Fatalln(err)
	}
	defer cur.Close(c.Ctx)
	for cur.Next(c.Ctx) {
		var res bson.M
		err := cur.Decode(&res)
		if err != nil {
			log.Fatalln(err)
		}
		filename := res["filename"].(string)
		if filename != "" {
			files = append(files, res["filename"].(string))
		}
	}
	return files
}

func (c *Conn) DeleteFile(filename string) error {
	bucket, err := gridfs.NewBucket(c.Database())
	if err != nil {
		return nil
	}
	// 查询file collection中所有的图片信息
	id, err := objectIDFromFileName(filename)
	if err != nil {
		return err
	}
	err = bucket.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
