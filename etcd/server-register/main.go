package main

import (
	"context"
	"flag"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type ServerRegister struct {
	cli       *clientv3.Client
	leaseID   clientv3.LeaseID
	keepalive <-chan *clientv3.LeaseKeepAliveResponse
	key       string
	val       string
}

var Endpoints = []string{"127.0.0.1:23791", "127.0.0.1:23792", "127.0.0.1:23793"}

// 设置租约
func (s *ServerRegister) putKeyWithLease(lease int64) error {
	// 设置租约时间
	resp, err := s.cli.Grant(context.Background(), lease)
	if err != nil {
		return err
	}
	// 注册服务，并绑定租约
	_, err = s.cli.Put(context.Background(), s.key, s.val, clientv3.WithLease(resp.ID))
	if err != nil {
		return err
	}

	leaseRespChan, err := s.cli.KeepAlive(context.Background(), resp.ID)
	if err != nil {
		return err
	}

	s.leaseID = resp.ID
	log.Println("Lease id=", resp.ID)
	s.keepalive = leaseRespChan
	log.Printf("Put key:%v val:%v success!\n", s.key, s.val)
	return nil
}

// 监听续租情况
func (s *ServerRegister) ListenRespLeaseChan() {
	for leaseKeepResp := range s.keepalive {
		log.Println("续约成功", leaseKeepResp)
	}
	log.Println("续约关闭")
}

// 注销服务
func (s *ServerRegister) Close() error {
	// 撤销租约
	if _, err := s.cli.Revoke(context.Background(), s.leaseID); err != nil {
		log.Println("撤销服务失败,", err)
		return err
	}
	log.Println("撤销服务")
	return s.cli.Close()
}

func NewServerRegister(endpoint []string, key, val string, lease int64) (*ServerRegister, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoint,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
		log.Printf("connect to etcd failed, err:%v\n", err)
		return nil, err
	}
	log.Println("connect to etcd success")

	ser := &ServerRegister{
		cli: cli,
		key: key,
		val: val,
	}

	err = ser.putKeyWithLease(lease)
	if err != nil {
		log.Printf("putKeyWithLease failed, err:%v\n", err)
		return nil, err
	}
	return ser, nil
}

func main() {
	node := flag.String("node", "", "服务器etcd节点")
	addr := flag.String("addr", "", "服务器地址")
	flag.Parse()
	if *node == "" || *addr == "" {
		flag.Usage()
		return
	}

	ser, err := NewServerRegister(Endpoints, *node, *addr, 5)
	if err != nil {
		log.Fatalln("NewServer failed,", err)
	}
	// 监听续租相应的chan
	go ser.ListenRespLeaseChan()
	defer ser.Close()

	select {}
}
