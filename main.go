package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
)


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
	localAddress, ok := r.Context().Value(http.LocalAddrContextKey).(net.Addr)
	if ok {
		//myIP,_,_ = net.SplitHostPort()
	}

}


func main() {
	PORT := "80"
	fileServer := http.FileServer(http.Dir("./static/index.html"))
	http.Handle("/",fileServer)
	http.HandleFunc("/from", formHandler)
	http.HandleFunc("/ping", pingHandler)
	
	fmt.Println("Server starting on %v ...",PORT)
	if err := http.ListenAndServe(PORT, nil); err != nil {
		log.Fatal(err)
	}
	
}
