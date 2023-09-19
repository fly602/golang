package webcache

import (
	"context"
	"io/ioutil"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

func CleanWebCache(cli *redis.Client) {
	keys, err := cli.Keys(context.Background(), "webpage:*").Result()
	if err != nil {
		log.Println("get web cache err")
		return
	}
	for _, key := range keys {
		_, err = cli.Del(context.Background(), key).Result()
		if err != nil {
			log.Println("clean web cache err")
			return
		}

	}
}

func GetWebPageContent(cli *redis.Client, url string) ([]byte, error) {
	// 生成用于在 Redis 中存储缓存的键
	cacheKey := "webpage:" + url

	// 尝试从 Redis 中获取缓存的网页内容
	cachedContent, err := cli.Get(context.Background(), cacheKey).Bytes()
	if err == redis.Nil {
		// 如果缓存不存在，则需要从源获取网页内容
		log.Println("Cache miss. Fetching data from source...")

		// 模拟从源获取网页内容的操作
		// 这里可以使用 HTTP 请求或其他方式获取网页内容
		sourceContent, err := loadWebFile(url)
		if err != nil {
			return nil, err
		}

		// 将获取的网页内容存储到 Redis 中，并设置过期时间
		err = cli.Set(context.Background(), cacheKey, sourceContent, 24*time.Hour).Err()
		if err != nil {
			return nil, err
		}

		return sourceContent, nil
	} else if err != nil {
		return nil, err
	}

	// 如果缓存存在，则直接返回缓存的网页内容
	log.Println("Cache hit. Returning cached data...")
	return cachedContent, nil
}

func loadWebFile(url string) ([]byte, error) {
	// 这里可以实现从源获取网页内容的逻辑
	// 例如使用 HTTP 请求库获取网页内容
	// 返回获取到的网页内容字符串

	return ioutil.ReadFile(url)
}
