//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

//-----------------------------------------------------------------------------------
// Return a submission by index or priority
//-----------------------------------------------------------------------------------
func teacher_getsHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	index, _ := strconv.Atoi(r.FormValue("index"))
	priority, _ := strconv.Atoi(r.FormValue("priority"))

	BoardsSem.Lock()
	defer BoardsSem.Unlock()

	selected := &Submission{}

	if index >= 0 {
		// Try to select by index first
		selected = WorkingSubs[index]
		WorkingSubs = append(WorkingSubs[:index], WorkingSubs[index+1:]...)
	} else if priority > 0 {
		// Try to select by priority: 1 (got it), 2 (help me)
		for i := 0; i < len(WorkingSubs); i++ {
			if WorkingSubs[i].Priority == priority {
				selected = WorkingSubs[i]
				WorkingSubs = append(WorkingSubs[:i], WorkingSubs[i+1:]...)
				Students[selected.Uid].SubmissionStatus = 1
				break
			}
		}
	} else {
		// Try to select the first highest priority
		first_sub_w_priority := []int{-1, -1, -1}
		for i := 0; i < len(WorkingSubs); i++ {
			p := WorkingSubs[i].Priority
			if first_sub_w_priority[p] == -1 {
				first_sub_w_priority[p] = i
			}
		}
		for i := len(first_sub_w_priority) - 1; i > 0; i-- {
			if first_sub_w_priority[i] != -1 {
				j := first_sub_w_priority[i]
				selected = WorkingSubs[j]
				WorkingSubs = append(WorkingSubs[:j], WorkingSubs[j+1:]...)
				Students[selected.Uid].SubmissionStatus = 1
				break
			}
		}
	}
	js, err := json.Marshal(selected)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

//-----------------------------------------------------------------------------------
func teacher_gets_queueHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	js, err := json.Marshal(WorkingSubs)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

//-----------------------------------------------------------------------------------
