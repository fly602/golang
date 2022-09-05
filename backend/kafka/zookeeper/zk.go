package zookeeper

import (
	"log"
	"sync"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

func Create(conn *zk.Conn, path string, value string) {
	acls := zk.WorldACL(zk.PermAll)
	_, err := conn.Create(path, []byte(value), 0, acls)
	if err != nil {
		log.Println("Create failed,", err)
	}
}

func Get(conn *zk.Conn, path string) {
	val, stat, err := conn.Get(path)
	if err != nil {
		log.Printf("Get [%s] failed,%s\n", path, err.Error())
		return
	}
	log.Printf("Get [%s] value=%s, stat=%+v\n", path, string(val), stat)
}

func Update(conn *zk.Conn, path string, value string) {
}

func WatchKid(conn *zk.Conn, path string) {
	for {
		ret, _, event, err := conn.ChildrenW(path)
		if err != nil {
			log.Printf("ChildrenW [%s] failed,%s\n", path, err.Error())
			return
		}
		log.Printf("Ready to watch [%s],%+v\n", path, ret)
		select {
		case val := <-event:
			log.Printf("Watched [%s] event=%+v", path, val)
		}
	}
}

func WatchKid2(conn *zk.Conn, path string) {
	for {
		ret, _, event, err := conn.GetW(path)
		if err != nil {
			log.Printf("GetW [%s] failed,%s\n", path, err.Error())
			return
		}
		log.Printf("Ready to GetW [%s],%+v\n", path, string(ret))
		select {
		case val := <-event:
			log.Printf("GetW [%s] event=%+v", path, val)
		}
	}
}

func GetChild(conn *zk.Conn, path string) []string {
	kids, _, err := conn.Children(path)
	if err != nil {
		log.Printf("GetKids [%s] failed,%s\n", path, err.Error())
		return nil
	}
	log.Printf("GetKids [%s],kids=%v\n", path, kids)
	return kids
}

func Check(conn *zk.Conn, path string) {
	ret, _, err := conn.Exists(path)
	if err != nil {
		log.Printf("CheckExists [%s] failed,%s\n", path, err.Error())
		return
	}
	log.Println("Check", path, " isexist=", ret)
}

func CheckW(conn *zk.Conn, path string) {
	for {
		ret, _, event, err := conn.ExistsW(path)
		if err != nil {
			log.Printf("CheckExists [%s] failed,%s\n", path, err.Error())
			return
		}
		log.Println("CheckW", path, " isexist=", ret)
		select {
		case val := <-event:
			log.Printf("CheckW [%s] event=%+v", path, val)
		}
	}
}

func Conn() (*zk.Conn, error) {
	hosts := []string{"127.0.0.1:21811", "127.0.0.1:21812", "127.0.0.1:21813"}
	conn, _, err := zk.Connect(hosts, time.Second*5)
	if err != nil {
		log.Fatal("Zookeeper connect failed,", err)
	}
	Create(conn, "/sanguo", "history")
	Create(conn, "/sanguo/shuguo", "liubei")
	Create(conn, "/sanguo/wuguo", "shunquan")
	Create(conn, "/sanguo/weiguo", "caocao")

	// Get
	Get(conn, "/sanguo")
	Get(conn, "/sanguo/shuguo")
	Get(conn, "/sanguo/wuguo")
	Get(conn, "/sanguo/weiguo")

	kids := GetChild(conn, "/sanguo")
	Check(conn, "/")
	go CheckW(conn, "/")
	go WatchKid(conn, "/")
	go WatchKid2(conn, "/")
	for _, kid := range kids {
		Check(conn, "/sanguo"+"/"+kid)
		go CheckW(conn, "/sanguo"+"/"+kid)
		go WatchKid(conn, "/sanguo"+"/"+kid)
		go WatchKid2(conn, "/sanguo"+"/"+kid)
	}
	zklock := zk.NewLock(conn, "/sanguo/lock", zk.WorldACL(zk.PermAll))

	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		go func() {
			zklock.Lock()
			time.Sleep(time.Second * 10)
			defer zklock.Unlock()
			wg.Done()
		}()
		wg.Add(1)

	}
	wg.Wait()
	return conn, err
}
