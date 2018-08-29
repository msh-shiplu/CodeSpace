//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	// "encoding/json"
	"fmt"
	// "log"
	"net/http"
)

//-----------------------------------------------------------------------------------
// When problems are deactivated, boards cleared, no new submissions are possibile.
//-----------------------------------------------------------------------------------
func teacher_deactivates_problemsHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	filename := r.FormValue("filename")
	if prob, ok := ActiveProblems[filename]; ok {
		prob.Active = false
		if len(prob.Answers) > 0 {
			fmt.Fprintf(w, "1")
		} else {
			fmt.Fprintf(w, "0")
		}
	} else {
		fmt.Fprintf(w, "-1")
	}
	// for k, p := range ActiveProblems {
	// 	fmt.Println(k, p)
	// }

	// filenames := make([]string, 0)
	// for fname, prob := range ActiveProblems {
	// 	if prob.Active {
	// 		prob.Active = false
	// 		if len(prob.Answers) > 0 {
	// 			filenames = append(filenames, fname)
	// 		}
	// 	}
	// }
	// for stid, _ := range Students {
	// 	Students[stid].Boards = make([]*Board, 0)
	// 	Students[stid].SubmissionStatus = 0
	// }
	// js, err := json.Marshal(filenames)
	// if err == nil {
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.Write(js)
	// } else {
	// 	log.Fatal(err)
	// }
}

//-----------------------------------------------------------------------------------
// Clear submissions, boards, statuses, and set all problems inactive.
//-----------------------------------------------------------------------------------
func teacher_clears_submissionsHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	WorkingSubs = make([]*Submission, 0)
	fmt.Fprintf(w, "Done.")
}
