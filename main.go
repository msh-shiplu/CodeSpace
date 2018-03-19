//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"math/rand"
	"net"
	"net/http"
	"path/filepath"
	"time"
)

//-----------------------------------------------------------------
func init_handlers() {
	http.HandleFunc("/test", Authorize(testHandler))
	http.HandleFunc("/student_shares", Authorize(student_sharesHandler))
	http.HandleFunc("/student_gets", Authorize(student_getsHandler))
	http.HandleFunc("/teacher_broadcasts", Authorize(teacher_broadcastsHandler))
	http.HandleFunc("/teacher_gets_passcode", Authorize(teacher_gets_passcodeHandler))
	http.HandleFunc("/teacher_adds_ta", AuthorizeLocalhost(teacher_adds_taHandler))
	http.HandleFunc("/teacher_registers", teacher_registersHandler)
	http.HandleFunc("/student_registers", student_registersHandler)
}

//-----------------------------------------------------------------
func informIPAddress() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err.Error() + "\n")
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && ipnet.IP.IsGlobalUnicast() {
			return ipnet.IP.String()
		}
	}
	return ""
}

//-----------------------------------------------------------------
func main() {
	server := informIPAddress()
	port := "8080"
	fmt.Println("*********************************************")
	fmt.Printf("*   GEM (%s)\n", VERSION)
	fmt.Printf("*   Server address: %s:%s\n", server, port)
	fmt.Println("*********************************************\n")
	rand.Seed(time.Now().UnixNano())
	db_name := filepath.Join(".", "gem.sqlite3")
	flag.StringVar(&db_name, "db", db_name, "user database (sqlite).")
	flag.Parse()
	init_handlers()
	init_database(db_name)
	load_teachers()
	load_students()
	err := http.ListenAndServe("0.0.0.0:"+port, nil)
	if err != nil {
		panic(err.Error() + "\n")
	}
}
