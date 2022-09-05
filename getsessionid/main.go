package main

import (
	"log"
	"os/exec"
)

func GetSessionID() string {
	cmd := exec.Command("cat", "/proc/self/sessionid")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal("GetSessionID failed")
	}
	log.Println("GetSessionID", string(out))
	return string(out)
}

func main() {
	GetSessionID()
}
