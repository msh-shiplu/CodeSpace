//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"fmt"
	"net/http"
)

//-----------------------------------------------------------------------------------
func studentGetHelpCode(w http.ResponseWriter, r *http.Request, who string, uid int) {
	filename := r.FormValue("filename")

	pid := 0
	prob, ok := ActiveProblems[filename]

	if ok {
		if !prob.Active {
			fmt.Fprint(w, "This problem is not active.")
			return
		}
		pid = prob.Info.Pid
		if _, ok := HelpEligibleStudents[pid][uid]; !ok {
			fmt.Fprint(w, "You are not elligible to help for this problem.")
		}
		if ok {
			fmt.Fprintf(w, "You are elligible to help in this problem.")
		}

	}

}

//-----------------------------------------------------------------------------------
