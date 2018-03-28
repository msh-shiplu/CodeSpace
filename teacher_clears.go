//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"fmt"
	"net/http"
)

//-----------------------------------------------------------------------------------
// When problems are deactivated, no new submissions are possibile.
//-----------------------------------------------------------------------------------
func teacher_deactivates_problemsHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	ActiveProblems = make(map[int]*ProblemFormat)
	fmt.Fprintf(w, "Done.")
}

//-----------------------------------------------------------------------------------
// All submissions are cleared.  All boards are cleared.
//-----------------------------------------------------------------------------------
func teacher_clearsHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	WorkingSubs = make([]*Submission, 0)
	for stid, _ := range Boards {
		Boards[stid] = make([]*Board, 0)
	}
	for stid, _ := range MessageBoards {
		MessageBoards[stid] = "White board is empty."
	}
	fmt.Fprintf(w, "Done.")
}
