package main

import "encoding/base64"

func getEdidBase64(edid [128]byte) string {
	return base64.StdEncoding.EncodeToString(edid)
}

func main() {

}
