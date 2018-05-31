//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	// "fmt"
	"net/http"
	// "strconv"
	// "time"
)

//-----------------------------------------------------------------------------------
func teacher_puts_backHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	// priority, _ := strconv.Atoi(r.FormValue("priority"))
	// pid, _ := strconv.Atoi(r.FormValue("pid"))
	// sid, _ := strconv.Atoi(r.FormValue("sid"))
	// stid, _ := strconv.Atoi(r.FormValue("stid"))

	// SubSem.Lock()
	// defer SubSem.Unlock()
	// if prob, ok := ActiveProblems[pid]; ok && prob.Active {
	// 	sub := &Submission{
	// 		Sid:      sid,
	// 		Uid:      stid,
	// 		Pid:      pid,
	// 		Content:  prob.Info.Description,
	// 		Filename: prob.Info.Filename,
	// 		Priority: priority,
	// 		At:       time.Now(),
	// 	}
	// 	WorkingSubs = append(WorkingSubs, sub)
	// 	fmt.Fprintf(w, "OK")
	// } else {
	// 	fmt.Fprintf(w, "This problem is not active.  Cannot put back.")
	// }
}

// //-----------------------------------------------------------------------------------
