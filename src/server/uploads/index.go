package uploads

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

func HandleFileUpload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("file")

	if err != nil {
		fmt.Println("Error retrieving file from form-data")
		fmt.Println(err)
		return
	}

	defer file.Close()

	fmt.Printf("Upload time: %v\n", time.Now())
	fmt.Printf("File Name: %v\n", handler.Filename)
	fmt.Printf("File Size: %v\n", handler.Size)
	fmt.Printf("MIME Type: %v\n", handler.Header.Get("Content-Type"))

	tempFile, err := os.Create(filepath.Join("uploads/shared", filepath.Base(handler.Filename)))

	if err != nil {
		fmt.Println(err)
		return
	}

	defer tempFile.Close()

	fileBytes, err := io.ReadAll(file)

	if err != nil {
		fmt.Println(err)
		return
	}

	tempFile.Write(fileBytes)
	fmt.Fprint(w, "Successfully uploaded file")
}

func GetFiles(w http.ResponseWriter, r *http.Request) {
	entries, err := os.ReadDir("uploads/shared")

	if err != nil {
		log.Fatal(err)
	}

	var files []map[string]interface{}

	for _, entry := range entries {

		if entry.IsDir() {
			continue
		}
		file, err := entry.Info()

		if err != nil {
			fmt.Println(err)
			continue
		}

		download, err := url.JoinPath("api/file/shared", file.Name())

		if err != nil {
			fmt.Println(err)
			continue
		}

		var fileData = map[string]interface{}{
			"name":     file.Name(),
			"size":     file.Size(),
			"lastmod":  file.ModTime().UTC(),
			"download": download,
		}

		files = append(files, fileData)
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(files); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
