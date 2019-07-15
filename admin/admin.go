package main

import (
	//"bufio"
	"encoding/csv"
	"fmt"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"os"
)

func dataFetch() []string {
	// Open the file
	csvfile, err := os.Open("pendingSessions")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	// Parse the file
	r := csv.NewReader(csvfile)
	//r := csv.NewReader(bufio.NewReader(csvfile))
	dataArray := []string{}

	// Iterate through the records
	for {
		// Read each record from csv
		sessionData, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		data := fmt.Sprintf("SessionID: \t%s\nIP Address: \t%s\nPublickey: \t%s\nUserName: \t%s\n", sessionData[0], sessionData[1], sessionData[2], sessionData[3])
		dataArray = append(dataArray, data)
	}
	return dataArray
}

func adminPage(w http.ResponseWriter, r *http.Request) {
	parsedData := dataFetch
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	fmt.Println("Admin page accessed from: " + ip)
	tpl := template.Must(template.ParseGlob("admin.html"))
	tpl.Execute(w, parsedData)
}

func main() {
	http.HandleFunc("/", adminPage)
	http.ListenAndServe("192.168.100.1:8001", nil)
}
