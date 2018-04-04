//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"fmt"
	"net/http"
	// "strconv"
)

//-----------------------------------------------------------------

func student_checks_inHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	// attendance is taken automatically by authorization
	fmt.Fprint(w, "Ok")
}

//-----------------------------------------------------------------
func student_periodic_updateHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	submission_stat := Students[uid].SubmissionStatus
	board_stat := 0
	if len(Students[uid].Boards) > 0 {
		board_stat = 1
	}
	Students[uid].SubmissionStatus = 0 // reset status after notifying student
	fmt.Fprintf(w, "%d;%d", submission_stat, board_stat)

	// pid, _ := strconv.Atoi(r.FormValue("pid"))
	// if prob, ok := ActiveProblems[pid]; ok {
	// 	fmt.Println(pid, prob.Active, prob.Attempts[uid])
	// 	if prob.Active && prob.Attempts[uid] >= 0 {
	// 		code = 1
	// 	}
	// 	if !prob.Active {
	// 		update_mesg = "Problem inactive"
	// 		fmt.Println("inactive")
	// 	}
	// 	if prob.Attempts[uid] < 0 {
	// 		update_mesg = "No more attempts"
	// 		fmt.Println("No more attempts")
	// 	}
	// }
}
