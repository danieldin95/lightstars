package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	fmt.Printf("%d", 100 << 20)
	http.HandleFunc("/upload", uploadHandler)
	http.ListenAndServe(":5050", nil)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s user-agent=%s", r.Method, r.Proto, r.URL.String(), r.UserAgent())
	if r.Method != http.MethodPost {
		http.Error(w, "only POST method is allowed", http.StatusBadRequest)
		return
	}
	// no more than 50MiB of memory
	r.ParseMultipartForm(50 << 20)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		http.Error(w, "failed to read uploadfile form field", http.StatusBadRequest)
		log.Printf("failed to read uploadfile form field: %v", err)
		return
	}
	defer file.Close()

	tempFile, err := os.OpenFile("temp/"+handler.Filename, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	io.Copy(tempFile, file)
}
