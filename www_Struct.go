package main

import ("fmt"
	"net"
	"net/http"
	"time"
	"io/ioutil"
	"log"
	"os"
	"math/rand"
)

func genSessionID() string {
	rand.Seed(time.Now().UnixNano())
	bytes := make([]byte, 48)
	for i := 0; i < 48; i++ {
		bytes[i] = byte(65 + rand.Intn(25))
	}
	return string(bytes)
}

func init_Session(w http.ResponseWriter, r *http.Request){
	ip,_,_ := net.SplitHostPort(r.RemoteAddr)
	dt := time.Now()
	sessionID := genSessionID()
	publicKey, err := ioutil.ReadFile("pub")
	if err != nil {
		log.Fatal("Failed reading publicKey file")
	}

	fmt.Fprintf(w, string(publicKey) + "\n")

	f, err := os.OpenFile("log", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	fmt.Println("Recieved connection from: " + ip)

	f.WriteString(ip + "\t\t" + dt.Format("01-02-2000 15:04:04 Mon") + "\t" + sessionID + "\n\n")
	defer f.Close()

	err := r.ParseForm()
	name := r.FormValue("publickey")
	fmt.Println(name)
}

func main() {
	http.HandleFunc("/", init_Session)
	http.ListenAndServe("192.168.100.1:8000", nil)
	fmt.Println("Recieved!")
}
