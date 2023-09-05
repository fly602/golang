package main

import "log"

var (
	dirPath = "./"
)

func main() {
	log.SetFlags(log.Lshortfile)
	dbTest()
	// uploadPictureByPath(dirPath)
	downloadPicture()
}
