//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"fmt"
	"net/http"
	"strconv"
)

//-----------------------------------------------------------------------------------
func teacher_puts_backHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	sid, _ := strconv.Atoi(r.FormValue("sid"))
	SubSem.Lock()
	defer SubSem.Unlock()
	if _, ok := Submissions[sid]; ok {
		WorkingSubs = append(WorkingSubs, Submissions[sid])
		fmt.Fprintf(w, "Submission has been put back into the queue.")
	} else {
		fmt.Fprintf(w, "Unknown submission.")
	}
}

// //-----------------------------------------------------------------------------------
