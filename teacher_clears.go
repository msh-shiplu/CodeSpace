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
// When problems are deactivated, no new submissions are possibile.
//-----------------------------------------------------------------------------------
func teacher_deactivates_problemsHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	active_pids := make([]int, 0)
	for pid, prob := range ActiveProblems {
		if prob.Active && len(prob.Answers) > 0 {
			active_pids = append(active_pids, pid)
		}
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
// All submissions are cleared.  All boards are cleared.
//-----------------------------------------------------------------------------------
func teacher_clearsHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	WorkingSubs = make([]*Submission, 0)
	for stid, _ := range Students {
		Students[stid].Boards = make([]*Board, 0)
		Students[stid].Status = ""
	}
	fmt.Fprintf(w, "Done.")
}
