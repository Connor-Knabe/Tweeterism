package main

import (
//	"encoding/json"

	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)


func ReadWebPage(webaddr string,RmtIP string) []byte {
	client := &http.Client{}

	reqest, err := http.NewRequest("GET", webaddr, nil)

	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(0)
	}

	reqest.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*SEPERATETHEM/*;q=0.8")
	reqest.Header.Add("Connection", "keep-alive")
	reqest.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:12.0) Gecko/20100101 Firefox/12.0")
	reqest.Header.Add("X-Forwarded-For",RmtIP)
	response, err := client.Do(reqest)

	defer response.Body.Close()

	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(0)
	}
	bodyByte, _ := ioutil.ReadAll(response.Body)
	return bodyByte
}

func returnobj(Inq string,RmtIP string) []byte {

	ResponseByte := ReadWebPage("http://maps.googleapis.com/maps/api/geocode/json?sensor=true&address=" + Inq,RmtIP)


	return ResponseByte
}

func WebInquery(w http.ResponseWriter, req *http.Request) {

	var RmtIP = req.Header.Get("X-Forwarded-For")
	if RmtIP == "" {
		RmtIP = req.RemoteAddr
	}
	log.Println("\t", RmtIP, "\t", req.Method, "\t", req.URL.Path)
	if req.URL.Path == "/inq/" || req.URL.Path == "/inq" {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Write([]byte("It Works!"))

		return

	}
	ResKey := strings.TrimLeft(req.URL.Path, "/inq/")
	InqRes := returnobj(ResKey,RmtIP)

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Write(InqRes)

	return
}

func FileHandler(w http.ResponseWriter, req *http.Request) {

	var RmtIP = req.Header.Get("X-Forwarded-For")
	if RmtIP == "" {
		RmtIP = req.RemoteAddr
	}
	log.Println("\t", RmtIP, "\t", req.Method, "\t", req.URL.Path)
	var Fs = http.FileServer(http.Dir("./log/"))
	Fs.ServeHTTP(w, req)
}



func main() {
	//SetLog("./log/log.txt")

	http.HandleFunc("/", FileHandler)
	http.HandleFunc("/inq/", WebInquery)
	http.HandleFunc("/inq", WebInquery)
	http.ListenAndServe(":5010", nil)

} 