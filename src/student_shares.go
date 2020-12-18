//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

//-----------------------------------------------------------------------------------
func student_sharesHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	content, filename := r.FormValue("content"), r.FormValue("filename")
	answer := r.FormValue("answer")
	test_cases := r.FormValue("testcases")
	priority, _ := strconv.Atoi(r.FormValue("priority"))
	sid := int64(0)
	correct_answer := ""
	complete := false
	var err error
	msg := "Your submission will be looked at soon."

	pid := 0
	prob, ok := ActiveProblems[filename]
	if ok {
		if !prob.Active {
			msg = "Problem is no longer active. But the teacher will look at your submission."
		} else {
			pid = prob.Info.Pid
			if _, ok := prob.Attempts[uid]; !ok {
				ActiveProblems[filename].Attempts[uid] = prob.Info.Attempts
			}
			if ActiveProblems[filename].Attempts[uid] == 0 {
				fmt.Fprintf(w, "This is not submitted because either you have reached the submission limit or your solution was previously graded correctly.")
				return
			}

			// Decrement attempts **only if students are not asking for help**
			if priority < 2 {
				ActiveProblems[filename].Attempts[uid] -= 1
				if ActiveProblems[filename].Attempts[uid] <= 3 {
					msg += fmt.Sprintf(" You have %d attempt(s) left.", ActiveProblems[filename].Attempts[uid])
				}
			}

			// Autograding if possible
			correct_answer = ActiveProblems[filename].Info.Answer
			if answer != "" {
				scoring_mesg := ""
				if correct_answer == answer {
					scoring_mesg = add_or_update_score("correct", pid, uid, 0, -1)
					ActiveProblems[filename].Attempts[uid] = 0 // This prevents further submission
					complete = true
				} else if ActiveProblems[filename].Info.ExactAnswer {
					scoring_mesg = add_or_update_score("incorrect", pid, uid, 0, -1)
					complete = true
				} else {
					scoring_mesg = "Answer appears to be incorrect. It will be looked at."
				}
				ActiveProblems[filename].Answers = append(ActiveProblems[filename].Answers, answer)
				fmt.Fprintf(w, scoring_mesg)
			}
			var result sql.Result
			if complete {
				result, err = AddSubmissionCompleteSQL.Exec(pid, uid, content, priority, time.Now(), time.Now())
			} else {
				result, err = AddSubmissionSQL.Exec(pid, uid, content, priority, time.Now())
			}
			if err != nil {

				log.Fatal(err)
			}
			sid, _ = result.LastInsertId()

			if test_cases != "" {
				rows, _ := Database.Query("select id from test_case where student_id=? and problem_id=?", uid, pid)
				tc_id := 0
				for rows.Next() {
					rows.Scan(&tc_id)
					break
				}
				rows.Close()
				if tc_id != 0 {
					UpdateTestCaseSQL.Exec(test_cases, time.Now(), tc_id)
				} else {
					AddTestCaseSQL.Exec(pid, uid, test_cases, time.Now())
				}

			}
			if ActiveProblems[filename].Attempts[uid] == 0 {
				if _, ok := HelpEligibleStudents[pid][uid]; !ok {
					HelpEligibleStudents[pid][uid] = true
					// fmt.Fprintf(w, "You are now elligible to help you friends. To help please click on 'Help Friends' button.")
					msg = msg + "\nYou are now elligible to help you friends. To help please click on 'Help Friends' button."
				}
			}

		}
	}
	if !complete {
		SubSem.Lock()
		defer SubSem.Unlock()
		sub := &Submission{
			Sid:      int(sid),
			Uid:      uid,
			Pid:      pid,
			Content:  content,
			Filename: filename,
			Priority: priority,
			At:       time.Now(),
			Name:     r.FormValue("name"),
		}
		WorkingSubs = append(WorkingSubs, sub)
		Submissions[int(sid)] = sub
		fmt.Fprintf(w, msg)
	}
}

//-----------------------------------------------------------------------------------
