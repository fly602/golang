package main

import (
	"fmt"

	"github.com/olivere/elastic"
)

func main() {
	client, err := elastic.NewClient(elastic.SetURL("http://127.0.0.1:49200"), elastic.SetSniff(false))
	if err != nil {
		fmt.Println("connect es err:", err)
	}
	fmt.Println(client)
}
