package main

import (
	"encoding/json"
	"fmt"
)

type config struct {
	UpdateMode int
	Version    string
	Enable     bool
}

func main() {
	// 示例JSON字符串
	var cfg *config = &config{
		UpdateMode: 1000,
	}
	jsonStr := `{"Version": "1.1.1","Enable":true}`
	if err := json.Unmarshal([]byte(jsonStr), &cfg); err != nil {
		fmt.Println("===>>>to here", err)
		return
	}
	fmt.Printf("===>>>to here:%+v\n", cfg)
}
