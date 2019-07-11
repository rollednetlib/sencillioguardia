package main

import (
	"bufio"
	"fmt"
	"math/rand"
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

func main() {
	clientSession := genSessionID()
	fmt.Println(clientSession)
}
