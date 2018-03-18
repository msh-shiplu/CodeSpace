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
	http.HandleFunc("/teacher_broadcast", Authorize(teacher_broadcastHandler))
	http.HandleFunc("/teacher_get_passcode", Authorize(teacher_get_passcodeHandler))
	http.HandleFunc("/setup_new_teacher", AuthorizeLocalhost(setup_new_teacherHandler))
	http.HandleFunc("/register_teacher", register_teacherHandler)
	http.HandleFunc("/register_student", register_studentHandler)
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
