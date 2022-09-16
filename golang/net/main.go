package main

import "net/http"

func main() {
	http.Handle("/file", http.StripPrefix("/file", http.FileServer(http.Dir("./static/"))))
	http.ListenAndServe(":30001", nil)
}
