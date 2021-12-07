//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

//-----------------------------------------------------------------------------------
func teacher_gets_passcodeHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	fmt.Fprintf(w, Passcode)
}

func student_gets_passcodeHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	fmt.Fprintf(w, Passcode)
}

//-----------------------------------------------------------------------------------
func testHandler(w http.ResponseWriter, r *http.Request) {
	// Show content of boards
	fmt.Println("Students:", len(Students))
	for _, st := range Students {
		// fmt.Printf("Uid: %d has %d pages. Status: %d\n", uid, len(st.Boards), st.SubmissionStatus)
		for i := 0; i < len(st.Boards); i++ {
			b := st.Boards[i]
			fmt.Printf("Attempts: %d, Filename: %s, Pid: %d, Answer: %s, len of content: %d\n",
				b.Attempts, b.Filename, b.Pid, b.Answer, len(b.Content))
		}
	}

	fmt.Printf("WorkingSubs: %d entries", len(WorkingSubs))
	for i := 0; i < len(WorkingSubs); i++ {
		fmt.Println(WorkingSubs[i].Sid, WorkingSubs[i].Uid, WorkingSubs[i].Pid, WorkingSubs[i].Priority)
		fmt.Println(WorkingSubs[i].Content)
	}
	fmt.Println()

	fmt.Println("ActiveProblems:", ActiveProblems)
	for fname, v := range ActiveProblems {
		fmt.Println(fname, v.Active, "Answers:", v.Answers, "Attempts:", v.Attempts)
		fmt.Println(fname, v.Info.Pid, v.Info.Merit, v.Info.Effort, v.Info.Attempts, v.Info.ExactAnswer, v.Info.Answer)
	}
	fmt.Fprintf(w, Passcode)
}

//-----------------------------------------------------------------------------------

func testcase_getsHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	filename := r.FormValue("file_name")
	rows, err := Database.Query("select id from problem where filename=?", filename)
	problem_id := 0
	for rows.Next() {
		rows.Scan(&problem_id)
		break
	}
	rows.Close()
	rows, err = Database.Query("select test_cases from test_case where problem_id=?", problem_id)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	var test_cases = ""
	for rows.Next() {
		var tc = ""
		rows.Scan(&tc)
		if tc != "" {
			test_cases += tc[1 : len(tc)-1]
		}
	}
	fmt.Fprintf(w, "["+test_cases+"]")
}

func logEvent(eventName string, userID int, userType, eventType, otherInfo string) {
	AddUserEventLogSQL.Exec(eventName, userID, userType, eventType, otherInfo, time.Now())
}
