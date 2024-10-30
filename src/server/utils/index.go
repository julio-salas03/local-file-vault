package utils

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func ServeOptimizedFile(filename string, w http.ResponseWriter, r *http.Request) {
	acceptedEncodings := r.Header.Values("Accept-Encoding")
	var filepath strings.Builder
	filepath.WriteString(filename)

	for i := 0; i < len(acceptedEncodings); i++ {
		if strings.Contains(acceptedEncodings[i], "br") {
			w.Header().Add("Content-Encoding", "br")
			filepath.WriteString(".br")
			break
		}
	}

	bytes, err := os.ReadFile(filepath.String())

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal Server Error")
		return
	}

	w.Write(bytes)
}
