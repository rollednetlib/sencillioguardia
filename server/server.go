package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	//	"net/http/cookiejar"
	"html/template"
	"os"
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
			database, err := os.OpenFile("database", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
			if err != nil {
				panic(err)
			}
			database.WriteString(clientSessionID + ":" + ip + ":" + clientPublicKey + ":" + clientUserName + ":" + time.Now().String() + "\n")
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
	case "sencillioguard-0.0.1":
		//		binExchange(w, r)
	default:
		wwwExchange(w, r)
	}
}

func main() {
	http.HandleFunc("/", exchange)
	http.ListenAndServe("192.168.100.1:8000", nil)
	clientSession := genSessionID()
	fmt.Println(clientSession)
}
