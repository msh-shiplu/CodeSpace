//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

//-----------------------------------------------------------------------------------
func student_sharesHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	content, filename := r.FormValue("content"), r.FormValue("filename")
	answer := r.FormValue("answer")
	priority, _ := strconv.Atoi(r.FormValue("priority"))
	pid, _ := strconv.Atoi(r.FormValue("pid"))
	sid := int64(0)
	correct_answer := ""
	msg := "Your submission will be looked at soon."
	if pid > 0 { // only keep in database submissions related to problems
		active_problem, ok := ActiveProblems[pid]
		if !ok {
			fmt.Fprintf(w, "This is not a known problem.")
			return
		}
		if !active_problem.Active {
			fmt.Fprintf(w, "This problem is no longer active.")
			return
		}
		if _, ok := ActiveProblems[pid].Attempts[uid]; !ok {
			ActiveProblems[pid].Attempts[uid] = ActiveProblems[pid].Info.Attempts
		}
		ActiveProblems[pid].Attempts[uid] -= 1
		if ActiveProblems[pid].Attempts[uid] < 0 {
			fmt.Fprintf(w, "Submission limit reached. Not submitted.")
			return
		}

		result, err := AddSubmissionSQL.Exec(pid, uid, content, priority, time.Now())
		if err != nil {
			panic(err)
		}
		sid, _ = result.LastInsertId()
		correct_answer = ActiveProblems[pid].Info.Answer
		if answer != "" {
			ActiveProblems[pid].Answers = append(ActiveProblems[pid].Answers, answer)
			if correct_answer == answer {
				// Auto-grading: set tid to 0
				scoring_mesg := add_or_update_score("correct", pid, uid, 0)
				fmt.Fprintf(w, scoring_mesg)
				return
			}
		}
		if ActiveProblems[pid].Attempts[uid] == 0 {
			msg += add_next_problem_to_board(pid, uid)
		}
	}
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
	}
	WorkingSubs = append(WorkingSubs, sub)
	fmt.Fprintf(w, msg)
}

//-----------------------------------------------------------------------------------
