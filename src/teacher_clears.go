//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//-----------------------------------------------------------------------------------
// When problems are deactivated, boards cleared, no new submissions are possibile.
//-----------------------------------------------------------------------------------
func teacher_deactivates_problemsHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	active_pids := make([]int, 0)
	for pid, prob := range ActiveProblems {
		if prob.Active {
			prob.Active = false
			if len(prob.Answers) > 0 {
				active_pids = append(active_pids, pid)
			}
		}
	}
	for stid, _ := range Students {
		Students[stid].Boards = make([]*Board, 0)
		Students[stid].SubmissionStatus = 0
	}
	js, err := json.Marshal(active_pids)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	} else {
		panic(err)
	}
}

//-----------------------------------------------------------------------------------
// Clear submissions, boards, statuses, and set all problems inactive.
//-----------------------------------------------------------------------------------
func teacher_clears_submissionsHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	WorkingSubs = make([]*Submission, 0)
	fmt.Fprintf(w, "Done.")
}
