//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path"
	"time"
)

type Configuration struct {
	IP      string
	Port    int
	Address string
	Log     string
	Max     int
}

var Config = &Configuration{}

type Record struct {
	Uid     string
	Address string
}

var Records = make(map[string]*Record)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

//-----------------------------------------------------------------------------
// Return a unique identifier to be used by Record
//-----------------------------------------------------------------------------
func newUid() string {
	n := 10
	for {
		new_uid := make([]rune, n)
		for i := range new_uid {
			new_uid[i] = letterRunes[rand.Intn(len(letterRunes))]
		}
		if _, ok := Records[string(new_uid)]; !ok {
			return string(new_uid)
		}
	}
}

//-----------------------------------------------------------------------------
func informIPAddress() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err)
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && ipnet.IP.IsGlobalUnicast() {
			return ipnet.IP.String()
		}
	}
	return ""
}

//-----------------------------------------------------------------------------
func writeLog(filename, message string) {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)
	// log.Println(time.Now().Format("Mon Jan _2 15:04:05 2006"), " ", message)
	log.Println(message)
}

//-----------------------------------------------------------------------------
// Use config.json which is presumed to exist in the current directory.
//-----------------------------------------------------------------------------
func initConfig() *Configuration {
	var pwd, filename string
	var err error
	var file *os.File
	pwd, err = os.Getwd()
	if err != nil {
		log.Fatal("Current directory is inaccessible.")
	}
	filename = path.Join(pwd, "config.json")
	file, err = os.Open(filename)
	if err != nil {
		log.Fatal("Could not open " + filename)
	}
	decoder := json.NewDecoder(file)
	config := &Configuration{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}
	if config.IP == "" {
		config.IP = informIPAddress()
	}
	config.Address = fmt.Sprintf("%s:%d", config.IP, config.Port)
	return config
}

//-----------------------------------------------------------------------------
func askHandler(w http.ResponseWriter, r *http.Request) {
	who := r.FormValue("who")
	rec, ok := Records[who]
	if !ok {
		fmt.Fprintf(w, "unknown")
	} else {
		fmt.Fprintf(w, "http://"+rec.Address)
	}
}

//-----------------------------------------------------------------------------
func tellHandler(w http.ResponseWriter, r *http.Request) {
	if len(Records) > Config.Max {
		fmt.Fprintf(w, "max_record_exceeded")
		return
	}

	who, uid, address := r.FormValue("who"), r.FormValue("uid"), r.FormValue("address")
	rec, ok := Records[who]
	if !ok || uid == "" {
		// Trusted party is new. Generate a new uid. Inform the party. Save record
		uid = newUid()
		Records[who] = &Record{
			Uid:     uid,
			Address: address,
		}
		fmt.Println("Registered", who, uid)
		writeLog(Config.Log, fmt.Sprintf("Register %s %s at %s.", who, uid, address))
		fmt.Fprintf(w, uid)
	} else if rec.Uid != uid {
		fmt.Println("Unmatched", who, uid, rec.Uid)
		writeLog(Config.Log, fmt.Sprintf("%s provided incorrect uid (%s).", who, uid))
		fmt.Fprintf(w, "unmatched")
	} else {
		writeLog(Config.Log, fmt.Sprintf("%s informed new address (%s).", who, address))
		fmt.Fprintf(w, "ok")
	}
}

//-----------------------------------------------------------------------------
func main() {
	Config = initConfig()
	fmt.Printf("Name server is running on http://%s\n", Config.Address)
	fmt.Println("Log is written to", Config.Log)

	http.HandleFunc("/tell", tellHandler)
	http.HandleFunc("/ask", askHandler)

	rand.Seed(time.Now().UnixNano())
	err := http.ListenAndServe(Config.Address, nil)
	if err != nil {
		log.Fatal(err)
	}
}
