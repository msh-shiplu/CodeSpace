//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"fmt"
	"net/http"
	"strconv"
	// "strings"
	"time"
)

//-----------------------------------------------------------------------------------
func teacher_gradesHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	content, decision := r.FormValue("content"), r.FormValue("decision")
	changed := r.FormValue("changed")
	pid, _ := strconv.Atoi(r.FormValue("pid"))
	stid, _ := strconv.Atoi(r.FormValue("stid"))
	mesg, student_mesg := "", ""
	if decision == "dismissed" {
		MessageBoards[stid] = "Teacher looked at but did not grade your submission."
		fmt.Fprintf(w, "Submission dismissed.")
		return
	}
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
			Boards[stid] = append(Boards[stid], b)
		}
	}
	scoring_mesg := add_or_update_score(decision, pid, stid, uid)
	mesg = scoring_mesg + "\n" + mesg
	student_mesg += scoring_mesg
	MessageBoards[stid] = student_mesg
	fmt.Fprintf(w, mesg)
}

//-----------------------------------------------------------------------------------
