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
	http.HandleFunc("/test", testHandler)

	// Analytics
	// http.HandleFunc("/learning_report", learning_reportHandler)
	http.HandleFunc("/analyze_submissions", analyze_submissionsHandler)
	http.HandleFunc("/view_activities", view_activitiesHandler)
	http.HandleFunc("/report", reportHandler)
	http.HandleFunc("/report_tag", report_tagHandler)
	http.HandleFunc("/view_answers", view_answersHandler)
	http.HandleFunc("/statistics", statisticsHandler)

	// Others
	http.HandleFunc("/student_periodic_update", Authorize(student_periodic_updateHandler))

	http.HandleFunc("/student_gets_report", Authorize(student_gets_reportHandler))
	http.HandleFunc("/student_checks_in", Authorize(student_checks_inHandler))
	http.HandleFunc("/student_shares", Authorize(student_sharesHandler))
	http.HandleFunc("/student_gets", Authorize(student_getsHandler))
	http.HandleFunc("/view_bulletin_board", view_bulletin_boardHandler)
	http.HandleFunc("/remove_bulletin_page", remove_bulletin_pageHandler)
	http.HandleFunc("/bulletin_board_data", bulletin_board_dataHandler)
	http.HandleFunc("/complete_registration", complete_registrationHandler)

	http.HandleFunc("/student_ask_help", Authorize(studentAskHelpHandler))
	http.HandleFunc("/student_get_help_code", Authorize(studentGetHelpCode))
	http.HandleFunc("/student_return_without_feedback", Authorize(student_return_without_feedbackHandler))
	http.HandleFunc("/student_send_help_message", Authorize(student_send_help_messageHandler))
	http.HandleFunc("/student_send_thank_you", Authorize(sendThankYouHandler))

	http.HandleFunc("/teacher_get_help_code", Authorize(teacherGetHelpCode))
	http.HandleFunc("/teacher_return_without_feedback", Authorize(teacher_return_without_feedbackHandler))
	http.HandleFunc("/teacher_send_help_message", Authorize(teacher_send_help_messageHandler))

	http.HandleFunc("/teacher_gets_queue", Authorize(teacher_gets_queueHandler))
	http.HandleFunc("/teacher_adds_bulletin_page", Authorize(teacher_adds_bulletin_pageHandler))
	http.HandleFunc("/teacher_clears_submissions", Authorize(teacher_clears_submissionsHandler))
	http.HandleFunc("/teacher_deactivates_problems", Authorize(teacher_deactivates_problemsHandler))
	http.HandleFunc("/teacher_grades", Authorize(teacher_gradesHandler))
	http.HandleFunc("/teacher_puts_back", Authorize(teacher_puts_backHandler))
	http.HandleFunc("/teacher_gets", Authorize(teacher_getsHandler))
	http.HandleFunc("/teacher_broadcasts", Authorize(teacher_broadcastsHandler))
	http.HandleFunc("/teacher_gets_passcode", Authorize(teacher_gets_passcodeHandler))
	http.HandleFunc("/student_gets_passcode", Authorize(student_gets_passcodeHandler))

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintf(w, "pong") })

	http.HandleFunc("/get_testcase", Authorize(testcase_getsHandler))

	http.HandleFunc("/code_snapshot", Authorize(codeSnapshotHandler))
	http.HandleFunc("/get_global_info", Authorize(globalInfoHandler))

	http.HandleFunc("/get_codespace", Authorize(codespaceHandler))
	http.HandleFunc("/get_snapshot", Authorize(getCodeSnapshotHandler))
	http.HandleFunc("/save_snapshot_feedback", Authorize(codeSnapshotFeedbackHandler))
	http.HandleFunc("/get_snapshot_feedback", Authorize(getSnapshotFeedbackHandler))
	http.HandleFunc("/save_snapshot_back_feedback", Authorize(studentSendBackFeedbackHandler))
	http.HandleFunc("/student_views_feedback", Authorize(studentViewsFeedbackHandler))
	http.HandleFunc("/teacher_views_feedback", Authorize(teacherViewsFeedbackHandler))

	http.HandleFunc("/help_requests", Authorize(helpRequestListHandler))
	http.HandleFunc("/view_help_request", Authorize(viewHelpRequestHandler))
	http.HandleFunc("/set_peer_tutor", Authorize(setPeerTutorHandler))

	http.HandleFunc("/view_exercises", Authorize(problemListHandler))
	http.HandleFunc("/problem_dashboard", Authorize(problemDashboardHandler))
	http.HandleFunc("/student_dashboard_feedback_provision", Authorize(studentDashboardFeedbackProvisionHandler))
	http.HandleFunc("/save_message_feedback", Authorize(messageFeedbackHandler))
	http.HandleFunc("/student_dashboard_submissions", Authorize(studentDashboardSubmissionHandler))
}

//-----------------------------------------------------------------
func informIPAddress() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err)
	}
	// for _, a := range addrs {
	// if ipnet, ok := a.(*net.IPNet); ok && ipnet.IP.IsGlobalUnicast() {
	// if ipnet.IP.To4() != nil {
	// return ipnet.IP.String()
	// }
	// }
	// }
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && ipnet.IP.IsGlobalUnicast() {
			ip4 := ipnet.IP.To4()
			if ip4 != nil {
				switch {
				// case ip4[0] == 10:
				case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
				case ip4[0] == 192 && ip4[1] == 168:
				default:
					return ip4.String()
				}
			}
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
		fmt.Println("Error", err)
		log.Fatal("Unable to contact with name server.")
	}
}

//-----------------------------------------------------------------
func main() {
	rand.Seed(time.Now().UnixNano())
	config_file, teacher_file, student_file := "", "", ""
	flag.StringVar(&config_file, "c", config_file, "json-formatted configuration file.")
	flag.StringVar(&teacher_file, "add_teachers", teacher_file, "teacher file.")
	flag.StringVar(&student_file, "add_students", student_file, "student file.")
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
	fmt.Println("**************************************************")
	fmt.Printf("*   Course id:      %s\n", Config.CourseId)
	if Config.NameServer != "" {
		fmt.Printf("*   Server address: %s\n", Config.NameServer)
	} else {
		fmt.Printf("*   Serving at:     %s\n", Config.Address)
	}
	fmt.Printf("*   GEM %s\n", VERSION)
	fmt.Println("**************************************************\n")
	err := http.ListenAndServe(Config.Address, nil)
	if err != nil {
		log.Fatal("Unable to serve gem server at " + Config.Address)
	}
}
