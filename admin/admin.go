package main

import (
	"encoding/csv"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
)

/*var htmlTemplate = `{{range $index, $element := .}}{{$index}}
{{range $element}}{{.}}
{{end}}
{{end}}`
*/
/*
type mdata struct {
	sessionID, ipAddress, publicKey, userName string
}

type sdata struct {
	session  string
	metadata []mdata
}
*/

func fetchData() map[string]map[string]string {
	f, err := os.Open("pendingSessions")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	lines, err := csv.NewReader(f).ReadAll()

	var data = map[string]map[string]string{}
	for x, line := range lines {
		c := strconv.Itoa(x)
		data[c] = map[string]string{}
		data[c]["sessionID"] = line[0]
		data[c]["ipAddress"] = line[1]
		data[c]["publicKey"] = line[2]
		data[c]["userName"] = line[3]
		data[c]["status"] = line[4]
		fmt.Println(c)
	}
	fmt.Println(data)
	return data
}

func adminPage(w http.ResponseWriter, r *http.Request) {
	//t := template.New("t")
	t, err := template.ParseFiles("adminpage.html")
	if err != nil {
		panic(err)
	}

	fmt.Println(fetchData())
	err = t.Execute(w, fetchData())
	if err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("/", adminPage)
	http.ListenAndServe("192.168.100.1:8003", nil)
}
