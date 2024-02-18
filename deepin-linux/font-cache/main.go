package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
)

type Family struct {
	Id   string
	Name string

	Styles []string

	Monospace bool
	Show      bool
}

type FamilyHashTable map[string]*Family

var table = make(FamilyHashTable)

func loadCacheFromFile(file string, obj interface{}) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	var r = bytes.NewBuffer(data)
	decoder := gob.NewDecoder(r)
	err = decoder.Decode(obj)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	loadCacheFromFile("family_hash", &table)
	for name, info := range table {
		fmt.Printf("name:%v info:%+v\n", name, info)
	}

}
