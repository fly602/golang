package main

import (
	"log"

	"github.com/godbus/dbus"
)

const (
	fingerPrint        = "com.huawei.Fingerprint"
	fingerPath         = "/com/huawei/Fingerprint"
	fingerSearchDevice = "com.huawei.Fingerprint.SearchDevice"
)

func main() {
	log.Println("Text launch ...")
	systemConn, err := dbus.SystemBus()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("sysobj ...")
	systemConnObj := systemConn.Object(fingerPrint, fingerPath)
	var val bool
	log.Println("call SearchDevice ...")
	err = systemConnObj.Call(fingerSearchDevice, 0).Store(&val)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("get value=", val)
}
