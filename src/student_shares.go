//
// Author: Vinhthuy Phan, 2018
//
package main

import (
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
	priority, _ := strconv.Atoi(r.FormValue("priority"))
	pid, _ := strconv.Atoi(r.FormValue("pid"))
	sid := int64(0)
	correct_answer := ""
	msg := "Your submission will be looked at soon."
	if pid > 0 { // only keep in database submissions related to problems
		if prob, ok := ActiveProblems[pid]; !ok || !prob.Active {
			fmt.Fprintf(w, "This problem is no longer active.")
			return
		}
		if _, ok := ActiveProblems[pid].Attempts[uid]; !ok {
			ActiveProblems[pid].Attempts[uid] = ActiveProblems[pid].Info.Attempts
		}
		if ActiveProblems[pid].Attempts[uid] == 0 {
			fmt.Fprintf(w, "This is not submitted because either you have reached the submission limit or your solution was previously graded correctly.")
			return
		}

		// Decrement attempts
		ActiveProblems[pid].Attempts[uid] -= 1
		if ActiveProblems[pid].Attempts[uid] <= 3 {
			msg += fmt.Sprintf(" You have %d attempt(s) left.", ActiveProblems[pid].Attempts[uid])
		}

		// Add to submission queue
		result, err := AddSubmissionSQL.Exec(pid, uid, content, priority, time.Now())
		if err != nil {
			log.Fatal(err)
		}
		sid, _ = result.LastInsertId()

		// Autograding if possible
		correct_answer = ActiveProblems[pid].Info.Answer
		if answer != "" {
			ActiveProblems[pid].Answers = append(ActiveProblems[pid].Answers, answer)
			if correct_answer == answer {
				scoring_mesg := add_or_update_score("correct", pid, uid, 0)
				ActiveProblems[pid].Attempts[uid] = 0 // This prevents further submission
				fmt.Fprintf(w, scoring_mesg)
				return
			}
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
