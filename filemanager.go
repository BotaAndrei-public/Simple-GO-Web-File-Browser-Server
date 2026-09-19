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

	
	entries, err := os.ReadDir(targetPath)
	if err != nil {
		http.Error(w, "Cloud not open directory", http.StatusInternalServerError)
		return
	}

	var files []FileItem
	for _, entry := range entries{
		files = append(files, FileItem{
			Name: entry.Name(),
			IsDir: entry.IsDir(),
		})
	}
	
	parentDir := ""
	if requestedDir != ""{
		parentDir = filepath.Dir(requestedDir)
		if parentDir == "."{
			parentDir = ""
		}
	}
	
	data := ManagerData{
		CurrentDir: requestedDir,
		ParentDir: parentDir,
		Files: files,
	}

	tmpl, err := template.ParseFiles("templates/manager.html")
	if err != nil {
		http.Error(w, "Html Template Error: " + err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}


// Upload Files
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" { return }

	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, "Error reading multipart stream", http.StatusBadRequest)
		return
	}

	currentDir := ""
	var filename string
	var srcFile io.Reader

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, "Stream reading failed", http.StatusInternalServerError)
			return
		}
		
	//	***
	//	if part.FormName() == "dir" {
	//		buf := new(strings.Builder)
	//		io.Copy(buf, part)
	//		currentDir = buf.String()
	//}
	
	// Read the target 
	if part.FormName() == "dir" {
		buf := new(strings.Builder)

		maxLength := io.LimitReader(part, 1<<20)

		_,err = io.Copy(buf,maxLength)
		if err != nil {
			http.Error(w, "Filed text too large or read error", http.StatusBadRequest)
			return
		}
		currentDir = buf.String()
	}

	//File stream
	if part.FormName() == "myFile"{
		filename = part.FileName()
		srcFile = part
		break
	}

	}


	if filename == "" || srcFile == nil{
		http.Error(w, "No file probider", http.StatusBadRequest)
		return
	}

	savePath := filepath.Join(storageDir, currentDir, filename)
	dst, err := os.Create(savePath)
	if err != nil{
		http.Error(w, "Disk error creating file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	
	_, err = io.Copy(dst, srcFile)
	if err != nil {
		http.Error(w, "Upload failed", http.StatusInternalServerError)
		return
	}

	http.Redirect(w,r,"/files?dir="+currentDir, http.StatusSeeOther)

}

// Delete File or Folder
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" { return }
	currentDir := r.FormValue("dir")
	filename := r.FormValue("filename")

	fullPath := filepath.Join(storageDir, currentDir, filename)
	os.RemoveAll(fullPath)

	http.Redirect(w, r, "/files?dir="+currentDir, http.StatusSeeOther)
}

// Download Files
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	fileParam := r.URL.Query().Get("file")
	fullPath := filepath.Join(storageDir, fileParam)

	if !strings.HasPrefix(filepath.Clean(fullPath), filepath.Clean(storageDir)) {
		http.Error(w, "Forbidden Access", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(fullPath))
	http.ServeFile(w, r, fullPath)
}


// Rename File or Folder
func RenameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" { return }
	currentDir := r.FormValue("dir")
	oldName := r.FormValue("oldName")
	newName := r.FormValue("newName")

	oldPath := filepath.Join(storageDir, currentDir, oldName)
	newPath := filepath.Join(storageDir, currentDir, newName)

	os.Rename(oldPath, newPath)

	http.Redirect(w, r, "/files?dir="+currentDir, http.StatusSeeOther)
}

// Create New Folder (Mkdir)
func MkdirHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" { return }
	currentDir := r.FormValue("dir")
	folderName := r.FormValue("folderName")

	newFolderPath := filepath.Join(storageDir, currentDir, folderName)
	os.MkdirAll(newFolderPath, os.ModePerm)

	http.Redirect(w, r, "/files?dir="+currentDir, http.StatusSeeOther)
}
