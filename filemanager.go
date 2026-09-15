package main

import (
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const storageDir = "./MySharedFolder"

type FileItem struct{
	Name string
	IsDir bool
}

type ManagerData struct{
	CurrentDir string
	ParentDir string
	Files []FileItem
}

func FileManagerViewHandler(w http.ResponseWriter, r *http.Request){
	//os.MkdirAll(storageDir, os.ModePerm)
	err := os.MkdirAll(storageDir, 0755)
	if err != nil {
		http.Error(w, "Dir storage Error", http.StatusInternalServerError)
		return
	}

	requestedDir := r.URL.Query().Get("dir")
	targetPath := filepath.Join(storageDir, requestedDir)

	if !strings.HasPrefix(filepath.Clean(targetPath), filepath.Clean(storageDir)){
		http.Error(w, "Unauthorized Access!", http.StatusForbidden)
		return
	}

	//TODO
	entries, err := os.ReadDir(targetPath)
	if err != nil {
		http.Error(w, "Cloud not open directory", http.StatusInternalServerError)
		return
	}
	


}

