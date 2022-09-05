package main

import (
	"context"
	"flag"
	"log"
	"sync"
	"time"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var Endpoints = []string{"127.0.0.1:23791", "127.0.0.1:23792", "127.0.0.1:23793"}

type ServerDiscovery struct {
	cli        *clientv3.Client
	serverList map[string]string
	lock       sync.Mutex
}

func NewServerDisCovery(endpoint []string) *ServerDiscovery {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoint,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
		log.Printf("connect to etcd failed, err:%v\n", err)
		return nil
	}
	return &ServerDiscovery{
		cli:        cli,
		serverList: make(map[string]string),
	}
}

// watchServer 初始化服务列表和监视
func (s *ServerDiscovery) WatchServer(prefix string) error {
	resp, err := s.cli.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		log.Println("Get key failed,", err)
		return err
	}
	for _, ev := range resp.Kvs {
		s.SetServerList(string(ev.Key), string(ev.Value))
	}
	go s.watcher(prefix)
	return nil

}

func (s *ServerDiscovery) watcher(prefix string) {
	rch := s.cli.Watch(context.Background(), prefix, clientv3.WithPrefix())
	log.Printf("watching prefix %v now......\n", prefix)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case mvccpb.PUT:
				s.SetServerList(string(ev.Kv.Key), string(ev.Kv.Value))
			case mvccpb.DELETE:
				s.DeleteServerList(string(ev.Kv.Key))
			}
		}
	}
}

func (s *ServerDiscovery) SetServerList(key, val string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.serverList[key] = val
	log.Println("put key:", key, "val:", val)

}

func (s *ServerDiscovery) DeleteServerList(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.serverList, key)
	log.Println("del key:", key)
}

func (s *ServerDiscovery) GetServer() []string {
	s.lock.Lock()
	defer s.lock.Unlock()
	addrs := make([]string, 0)
	for _, v := range s.serverList {
		addrs = append(addrs, v)
	}
	return addrs
}

func (s *ServerDiscovery) Close() error {
	return s.cli.Close()
}

func main() {
	prefix := flag.String("prefix", "", "服务器节点或者节点前缀")
	flag.Parse()
	if *prefix == "" {
		flag.Usage()
		return
	}
	ser := NewServerDisCovery(Endpoints)
	defer ser.Close()
	ser.WatchServer(*prefix)
	for {
		select {
		case <-time.Tick(10 * time.Second):
			log.Println(ser.GetServer())
		}
	}
}
