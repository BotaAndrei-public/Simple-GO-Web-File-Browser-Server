package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

var FILE_NAME string = "server_notes.log"

func formHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil{
		fmt.Fprintf(w, "ParserFrom() err: %v", err)
		return
	}
	fmt.Fprint(w, "POST request successful!") 
	name := r.FormValue("name")
	address := r.FormValue("address")
	text := r.FormValue("text")

	//Write In File
	file, err := os.OpenFile(FILE_NAME, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		log.Printf("File Error: %v", err)
		return
	}
	defer file.Close()

	timestamp := time.Now().Format(time.DateTime)
	logLine := fmt.Sprintf("[%s], Name: %s | Address: %s | Text: %s\n", timestamp, name, address, text)

	if _, err := file.WriteString(logLine); err != nil {
		log.Printf("[FileName: %s] Error on write!\n[Error] %v ",FILE_NAME, err);
	}

}


func pingHandler (w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/ping"{
		http.Error(w, "404 not found", http.StatusNotFound)
		return	
	}
	
	if r.Method != "GET" {
		http.Error(w, "method is not supported", http.StatusNotFound)
		return
	}

	//TODO
	// Get local PC address
	var myIP string
	localAddress, ok := r.Context().Value(http.LocalAddrContextKey).(net.Addr)
	if ok {
		myIP,_,_ = net.SplitHostPort(localAddress.String())
	}else{
		myIP = "Not detect your IP-X.X.X.X"
	}
	
	//Send IP
	fmt.Fprint(w, "IP-"+myIP)

}


func main() {
	PORT := "127.0.0.1:8080"
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/",fileServer)
	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/ping", pingHandler)
	
	fmt.Printf("Server starting on %v ...\n",PORT)
	if err := http.ListenAndServe(PORT, nil); err != nil {
		log.Fatal(err)
	}
	
}
