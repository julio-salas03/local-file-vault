package uploads

import (
	"fmt"
	"io"
	"local-file-vault/api"
	"local-file-vault/auth"
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

func GetFilesFromFolder(folder string) ([]map[string]interface{}, error) {
	entries, err := os.ReadDir(fmt.Sprintf("uploads/%s", folder))

	var files []map[string]interface{}

	if err != nil {
		return files, fmt.Errorf("couldn't read files from %s", folder)
	}

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
			"owner":    folder,
		}

		files = append(files, fileData)
	}

	return files, nil
}

func GetFiles(w http.ResponseWriter, r *http.Request) {

	var allowedFolders []string

	allowedFolders = append(allowedFolders, "shared")

	if authCookie, err := r.Cookie(auth.AuthCookieName); err == nil && authCookie.Value != "" {
		username, err := auth.GetUserFromAuthCookie(authCookie.Value)
		if err == nil && username != "" {
			allowedFolders = append(allowedFolders, username)
		}
	}

	var files []map[string]interface{}

	for i := 0; i < len(allowedFolders); i++ {
		folder := allowedFolders[i]
		_files, err := GetFilesFromFolder(folder)

		if err != nil {
			fmt.Println(err)
			continue
		}

		files = append(files, _files...)
	}

	response := api.Response{
		Message: fmt.Sprintf("Successfully retrieved %d files", len(files)),
		Data: map[string]interface{}{
			"files": files,
		},
	}

	api.WriteResponse(w, response)
}
