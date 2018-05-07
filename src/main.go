//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	// "io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"time"
)

//-----------------------------------------------------------------
func init_handlers() {
	http.HandleFunc("/test", Authorize(testHandler))

	// Analytics
	// http.HandleFunc("/learning_report", learning_reportHandler)
	http.HandleFunc("/analyze_submissions", analyze_submissionsHandler)
	http.HandleFunc("/view_activities", view_activitiesHandler)
	http.HandleFunc("/view_tags", view_tagsHandler)
	http.HandleFunc("/report_tag", report_tagHandler)

	// Others
	http.HandleFunc("/student_periodic_update", Authorize(student_periodic_updateHandler))

	http.HandleFunc("/student_gets_report", Authorize(student_gets_reportHandler))
	http.HandleFunc("/student_checks_in", Authorize(student_checks_inHandler))
	http.HandleFunc("/student_shares", Authorize(student_sharesHandler))
	http.HandleFunc("/student_gets", Authorize(student_getsHandler))
	http.HandleFunc("/view_bulletin_board", view_bulletin_boardHandler)
	http.HandleFunc("/remove_bulletin_page", remove_bulletin_pageHandler)
	http.HandleFunc("/bulletin_board_data", bulletin_board_dataHandler)
	http.HandleFunc("/view_answers", view_answersHandler)
	http.HandleFunc("/complete_registration", complete_registrationHandler)

	http.HandleFunc("/teacher_adds_bulletin_page", Authorize(teacher_adds_bulletin_pageHandler))
	http.HandleFunc("/teacher_clears_submissions", Authorize(teacher_clears_submissionsHandler))
	http.HandleFunc("/teacher_deactivates_problems", Authorize(teacher_deactivates_problemsHandler))
	http.HandleFunc("/teacher_shares", Authorize(teacher_sharesHandler))
	http.HandleFunc("/teacher_grades", Authorize(teacher_gradesHandler))
	http.HandleFunc("/teacher_puts_back", Authorize(teacher_puts_backHandler))
	http.HandleFunc("/teacher_gets", Authorize(teacher_getsHandler))
	http.HandleFunc("/teacher_broadcasts", Authorize(teacher_broadcastsHandler))
	http.HandleFunc("/teacher_gets_passcode", Authorize(teacher_gets_passcodeHandler))

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintf(w, "pong") })
}

//-----------------------------------------------------------------
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

//-----------------------------------------------------------------
func init_config(filename string) *Configuration {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
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

//-----------------------------------------------------------------
func inform_name_server() {
	nameserver := fmt.Sprintf("%s/tell?who=%s&address=%s", Config.NameServer, Config.CourseId, Config.Address)
	_, err := http.Get(nameserver)
	if err != nil {
		log.Fatal("Unable to contact with name server.")
	}
	// defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(">", string(body))
}

//-----------------------------------------------------------------
func main() {
	rand.Seed(time.Now().UnixNano())
	config_file, teacher_file, student_file := "", "", ""
	flag.StringVar(&config_file, "c", config_file, "json-formatted configuration file.")
	flag.StringVar(&teacher_file, "add_teacher", teacher_file, "teacher file.")
	flag.StringVar(&student_file, "add_student", student_file, "student file.")
	flag.Parse()
	if config_file == "" {
		flag.Usage()
		os.Exit(1)
	}
	Config = init_config(config_file)
	if Config.NameServer != "" {
		inform_name_server()
	}
	init_database(Config.Database)
	if teacher_file != "" {
		add_multiple(teacher_file, "teacher")
	}
	if student_file != "" {
		add_multiple(student_file, "student")
	}
	init_handlers()
	load_teachers()
	fmt.Println("*********************************************")
	fmt.Printf("*   GEM (%s)\n", VERSION)
	fmt.Printf("*   Server address: %s\n", Config.Address)
	fmt.Println("*********************************************\n")
	err := http.ListenAndServe(Config.Address, nil)
	if err != nil {
		log.Fatal("Unable to serve gem server at " + Config.Address)
	}
}
