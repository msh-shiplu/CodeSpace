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
	sid, _ := strconv.Atoi(r.FormValue("sid"))
	changed := r.FormValue("changed")
	// stid, _ := strconv.Atoi(r.FormValue("stid"))
	// pid, _ := strconv.Atoi(r.FormValue("pid"))
	mesg := ""

	sub, ok := Submissions[sid]
	if !ok {
		fmt.Fprintf(w, "Unknown submission cannot be graded.")
		return
	}
	stid := sub.Uid
	if changed == "True" {
		// If the original file is changed, there's feedback.  Copy it to whiteboard.
		if prob, ok := ActiveProblems[sub.Filename]; ok {
			AddFeedbackSQL.Exec(uid, stid, content, time.Now())
			mesg = "Feedback saved to student's board."
			BoardsSem.Lock()
			defer BoardsSem.Unlock()
			b := &Board{
				Content:      content,
				Answer:       prob.Info.Answer,
				Attempts:     0, // This tells the client this is an existing problem
				Filename:     sub.Filename,
				Pid:          sub.Pid,
				StartingTime: time.Now(),
				Type:         "feedback",
			}
			Students[stid].Boards = append(Students[stid].Boards, b)
		}
	}

	// If submission is dismissed, do not take that attempt away from the student.
	if decision == "dismissed" {
		Students[stid].SubmissionStatus = 2
		ActiveProblems[sub.Filename].Attempts[stid] += 1
		fmt.Fprintf(w, "Submission dismissed.")
	} else if decision == "ungraded" {
		Students[stid].SubmissionStatus = 5
		fmt.Fprintf(w, mesg)
	} else {
		// Update score based on the grading decision
		scoring_mesg := add_or_update_score(decision, sub.Pid, sub.Uid, uid)
		mesg = scoring_mesg + "\n" + mesg
		if decision == "correct" {
			Students[sub.Uid].SubmissionStatus = 4
			ActiveProblems[sub.Filename].Attempts[stid] = 0 // This prevents further submission.
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
