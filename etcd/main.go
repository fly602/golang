package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

// etcd client put/get demo
// use etcd/clientv3

var Endpoints = []string{"127.0.0.1:23791", "127.0.0.1:23792", "127.0.0.1:23793"}

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   Endpoints,
		DialTimeout: 5 * time.Second,
	})
	if cli == nil || err != nil {
		// handle error!
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	fmt.Println("connect to etcd success")
	defer cli.Close()
	ctx, cancel := context.WithCancel(context.Background())
	if ctx == nil {
		log.Fatal("Ctx error", err)
	}
	defer cancel()
	res, err := cli.Get(ctx, "/root/key", clientv3.WithLastCreate()...)
	if err != nil {
		log.Fatal("Get key err,", err)
	}
	for _, ev := range res.Kvs {
		log.Println("Get key=", ev)
	}
	watcher := clientv3.NewWatcher(cli)
	watchch := watcher.Watch(context.TODO(), "num")
	go func() {
		for res := range watchch {
			//log.Printf("Wathcer response=%+v\n", res)
			for _, event := range res.Events {
				log.Printf("[%v]: %v=%v, Create=%v, Revision=%v\n",
					event.Type, string(event.Kv.Key), string(event.Kv.Value), event.Kv.CreateRevision, event.Kv.ModRevision)
			}
			//time.Sleep(time.Second)
		}

	}()
	defer watcher.Close()
	log.Println("Put num...")
	for i := 0; i < 100; i++ {
		cli.Put(context.TODO(), "num", fmt.Sprintf("%d", i))
	}
	time.Sleep(time.Second * 5)

	log.Println("Get key...")
	txn := cli.Txn(context.Background())
	_, err = txn.If(clientv3.Compare(clientv3.Version("key"), "=", 2)).Then(clientv3.OpGet("key")).Else(clientv3.OpPut("key", "aaaa")).Commit()
	if err != nil {
		log.Fatalln("txn Commit err,", err)
	}

	// 分布式锁
	var glob = 0
	g := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		g.Add(1)
		go func() {
			session, err := concurrency.NewSession(cli)
			if err != nil {
				log.Println("Get session err")
			}
			mtx := concurrency.NewMutex(session, "/root/lock")
			mtx.Lock(context.Background())
			glob++
			log.Println("glob++ =", glob)
			mtx.Unlock(context.Background())
			g.Done()
		}()
	}
	g.Wait()
	log.Println("get glob=", glob)
}
