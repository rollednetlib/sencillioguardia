package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	//	"net/http/cookiejar"
	"encoding/csv"
	"html/template"
	"os"
	"strconv"
	"strings"
	"time"
)

func genSessionString() string {
	rand.Seed(time.Now().UnixNano())
	bytes := make([]byte, 128)
	for i := 0; i < 128; i++ {
		bytes[i] = byte(65 + rand.Intn(25))
	}
	return string(bytes)
}

func genSessionID() string {
	sessionString := genSessionString()
	f, err := os.OpenFile("sessionjar", os.O_CREATE, 0660)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), sessionString) {
			sessionString := genSessionString()
			fmt.Println(sessionString)
		}
	}
	jar, err := os.OpenFile("sessionjar", os.O_RDWR|os.O_APPEND, 0660)
	jar.WriteString(sessionString + "\n")
	defer jar.Close()
	return sessionString
}

type cookedPage struct {
	SessionID     string
	SessionStatus string
}

func wwwExchange(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("sessionID")
	if cookie != nil {
		t, _ := template.ParseFiles("cookedSession.html")
		p := cookedPage{
			SessionID:     cookie.Value,
			SessionStatus: "Pending",
		}
		t.Execute(w, p)
	} else {
		serverPublicKey, err := ioutil.ReadFile("pub")
		if err != nil {
			panic(err)
		}
		switch r.Method {
		case "GET":
			http.ServeFile(w, r, "requestForm.html")
		case "POST":
			if err := r.ParseForm(); err != nil {
				panic(err)
			}
			clientPublicKey := r.FormValue("publickey")
			clientUserName := r.FormValue("username")

			clientSessionID := genSessionID()
			expiration := time.Now().Add(28 * 24 * time.Hour)
			http.SetCookie(w, &http.Cookie{
				Name:    "sessionID",
				Value:   clientSessionID,
				Expires: expiration})
			fmt.Fprintf(w, string(serverPublicKey)+"\n")
			fmt.Fprintf(w, clientSessionID+"\n")
			//			databaseCreate, err := os.OpenFile("database", os.O_CREATE, 0660)
			if err != nil {
				panic(err)
			}
			ip, _, _ := net.SplitHostPort(r.RemoteAddr)
			database, err := os.OpenFile("pendingSessions", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
			if err != nil {
				panic(err)
			}
			database.WriteString(clientSessionID + "," + ip + "," + clientPublicKey + "," + clientUserName + ",pending," + time.Now().String() + "\n")
		}
	}
}

func exchange(w http.ResponseWriter, r *http.Request) {
	/*	if r.URL.Path != "/" {
		http.Error(w, "404", http.StatusNotFound)
		return
	} */
	if r.URL.Path == "/clear" {
		//		fmt.Fprint(w, "Eliminar cookies")
		http.SetCookie(w, &http.Cookie{
			Name:    "sessionID",
			MaxAge:  -1,
			Expires: time.Now(),
		})
		fmt.Println("Eliminar session")
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	fmt.Println(ip)
	//Sort connections based on bin or web
	switch r.Header.Get("User-Agent") {
	case "sencillioguard-binaryagent-0.0.1":
		//		binExchange(w, r)
	default:
		wwwExchange(w, r)
	}
}

/*
func main() {
	http.HandleFunc("/", exchange)
	http.ListenAndServe("192.168.100.1:8000", nil)
	clientSession := genSessionID()
	fmt.Println(clientSession)
}

package main

import (
	"encoding/csv"
	"html/template"
	"net/http"
	"os"
	"strconv"
)
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
	}
	return data
}

func adminPage(w http.ResponseWriter, r *http.Request) {
	//t := template.New("t")
	t, err := template.ParseFiles("adminpage.html")
	if err != nil {
		panic(err)
	}

	err = t.Execute(w, fetchData())
	if err != nil {
		panic(err)
	}
}

/*
func main() {
	http.HandleFunc("/", exchange)
	http.ListenAndServe("192.168.100.1:8000", nil)
	clientSession := genSessionID()
	fmt.Println(clientSession)
}

func main() {
	http.HandleFunc("/", adminPage)
	http.ListenAndServe("192.168.100.1:8003", nil)
}
*/

func readConfig() (string, string, string, string) {
	f, err := os.Open("sg.conf")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	lines, err := csv.NewReader(f).ReadAll()

	var publicKey string
	var privateKey string
	var publicBind string
	var adminBind string
	for _, each := range lines {
		if each[0] == "publicKey" {
			publicKey = each[1]
		}
		if each[0] == "privateKey" {
			privateKey = each[1]
		}
		if each[0] == "publicBind" {
			publicBind = each[1]
		}
		if each[0] == "adminBind" {
			adminBind = each[1]
		}
	}
	return publicKey, privateKey, publicBind, adminBind
}

func main() {

	publicKey, _, publicBind, adminBind := readConfig()
	finish := make(chan bool)

	fmt.Println("publicKey: " + publicKey)
	serverPub := http.NewServeMux()
	serverPub.HandleFunc("/", exchange)

	serverAdmin := http.NewServeMux()
	serverAdmin.HandleFunc("/", adminPage)

	go func() {
		http.ListenAndServe(publicBind, serverPub)
	}()

	go func() {
		http.ListenAndServe(adminBind, serverAdmin)
	}()

	<-finish
}

/*
func foo8001(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Listening on 8001: foo "))
}

func bar8001(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Listening on 8001: bar "))
}

func foo8002(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Listening on 8002: foo "))
}

func bar8002(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Listening on 8002: bar "))
}
*/
