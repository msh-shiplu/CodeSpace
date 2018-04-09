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
func teacher_gradesHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	content, decision := r.FormValue("content"), r.FormValue("decision")
	changed := r.FormValue("changed")
	pid, _ := strconv.Atoi(r.FormValue("pid"))
	stid, _ := strconv.Atoi(r.FormValue("stid"))
	sid, _ := strconv.Atoi(r.FormValue("sid"))
	mesg, student_mesg := "", ""

	// If the original file is changed, there's feedback.  Copy it to whiteboard.
	if changed == "True" {
		if prob, ok := ActiveProblems[pid]; ok {
			AddFeedbackSQL.Exec(uid, stid, content, time.Now())
			mesg = "Feedback saved to student's board."
			student_mesg += "You've got feedback."
			BoardsSem.Lock()
			defer BoardsSem.Unlock()
			b := &Board{
				Content:      content,
				Answer:       prob.Info.Answer,
				Attempts:     0, // This tells the client this is an existing problem
				Filename:     prob.Info.Filename,
				Pid:          pid,
				StartingTime: time.Now(),
			}
			Students[stid].Boards = append(Students[stid].Boards, b)
		}
	}

	// If submission is dismissed, do not take that attempt away from the student.
	if decision == "dismissed" {
		Students[stid].SubmissionStatus = 2
		ActiveProblems[pid].Attempts[stid] += 1
		fmt.Fprintf(w, "Submission dismissed.")
	} else {
		// Update score based on the grading decision
		scoring_mesg := add_or_update_score(decision, pid, stid, uid)
		mesg = scoring_mesg + "\n" + mesg
		student_mesg += scoring_mesg
		if decision == "correct" {
			Students[stid].SubmissionStatus = 4
			ActiveProblems[pid].Attempts[stid] = 0 // This prevents further submission.
		} else {
			Students[stid].SubmissionStatus = 3
		}

		// Update submission complete time
		_, err := CompleteSubmissionSQL.Exec(time.Now(), sid)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, mesg)
	}
}

//-----------------------------------------------------------------------------------
