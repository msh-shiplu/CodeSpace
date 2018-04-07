//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"math/rand"
	"net"
	"net/http"
	"os"
	// "path/filepath"
	"time"
)

//-----------------------------------------------------------------
func init_handlers() {
	http.HandleFunc("/test", Authorize(testHandler))
	http.HandleFunc("/student_periodic_update", Authorize(student_periodic_updateHandler))

	http.HandleFunc("/student_gets_report", Authorize(student_gets_reportHandler))
	http.HandleFunc("/student_checks_in", Authorize(student_checks_inHandler))
	http.HandleFunc("/student_shares", Authorize(student_sharesHandler))
	http.HandleFunc("/student_gets", Authorize(student_getsHandler))
	http.HandleFunc("/student_registers", student_registersHandler)
	// http.HandleFunc("/show_student_messages", student_messagesHandler)
	http.HandleFunc("/view_bulletin_board", view_bulletin_boardHandler)
	http.HandleFunc("/remove_bulletin_page", remove_bulletin_pageHandler)
	http.HandleFunc("/bulletin_board_data", bulletin_board_dataHandler)
	http.HandleFunc("/view_answers", view_answersHandler)

	http.HandleFunc("/teacher_adds_bulletin_page", Authorize(teacher_adds_bulletin_pageHandler))
	http.HandleFunc("/teacher_clears", Authorize(teacher_clearsHandler))
	http.HandleFunc("/teacher_deactivates_problems", Authorize(teacher_deactivates_problemsHandler))
	http.HandleFunc("/teacher_shares", Authorize(teacher_sharesHandler))
	http.HandleFunc("/teacher_grades", Authorize(teacher_gradesHandler))
	http.HandleFunc("/teacher_puts_back", Authorize(teacher_puts_backHandler))
	http.HandleFunc("/teacher_gets", Authorize(teacher_getsHandler))
	http.HandleFunc("/teacher_broadcasts", Authorize(teacher_broadcastsHandler))
	http.HandleFunc("/teacher_gets_passcode", Authorize(teacher_gets_passcodeHandler))
	http.HandleFunc("/teacher_completes_registration", teacher_completes_registrationHandler)
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
func init_config(filename string) *Configuration {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(file)
	config := &Configuration{}
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}
	if config.IP == "" {
		config.IP = informIPAddress()
	}
	config.Address = fmt.Sprintf("%s:%d", config.IP, config.Port)
	fmt.Println(config)
	return config
}

//-----------------------------------------------------------------
func main() {
	rand.Seed(time.Now().UnixNano())
	new_teacher, new_ta, config_file := "", "", ""
	flag.StringVar(&config_file, "config", config_file, "configuration file.")
	flag.StringVar(&new_teacher, "add_teacher", new_teacher, "add a new teacher.")
	flag.StringVar(&new_ta, "add_ta", new_ta, "add a new teaching assistant.")
	flag.Parse()
	Config = init_config(config_file)
	init_database(Config.Database)
	if new_teacher != "" {
		add_teacher(new_teacher)
	} else if new_ta != "" {
		add_teacher(new_ta)
	} else {
		init_handlers()
		load_teachers()
		fmt.Println("*********************************************")
		fmt.Printf("*   GEM (%s)\n", VERSION)
		fmt.Printf("*   Server address: %s\n", Config.Address)
		fmt.Println("*********************************************\n")
		err := http.ListenAndServe(Config.Address, nil)
		if err != nil {
			panic(err.Error() + "\n")
		}
	}
}
