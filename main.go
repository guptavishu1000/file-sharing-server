package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"io" // Add this import for the 'io' package
)

const (
	uploadDirectory = "uploads/"
	port            = "8080"
)

func main() {
	http.HandleFunc("/", fileHandler)
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir(uploadDirectory))))

	fmt.Printf("Server is listening on :%s\n", port)
	http.ListenAndServe(":"+port, nil)
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "index.html")
	} else if r.Method == http.MethodPost {
		r.ParseMultipartForm(10 << 20) // 10 MB limit

		file, handler, err := r.FormFile("file")
		if err != nil {
			fmt.Println("Error reading file:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		filePath := filepath.Join(uploadDirectory, handler.Filename)
		out, err := os.Create(filePath)
		if err != nil {
			fmt.Println("Error creating file:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			fmt.Println("Error copying file:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "File uploaded successfully.")
	}
}
