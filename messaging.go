//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type TeacherMessage struct {
	P1             int
	P2             int
	ActiveProblems int
}

//-----------------------------------------------------------------------------------
func teacher_messagesHandler(w http.ResponseWriter, r *http.Request) {
	passcode := r.FormValue("pc")
	if passcode == Passcode {
		t := template.New("")
		t, err := t.Parse(TEACHER_MESSAGING_TEMPLATE)
		if err == nil {
			priority := []int{0, 0, 0}
			for i := 0; i < len(WorkingSubs); i++ {
				priority[WorkingSubs[i].Priority]++
			}
			data := &TeacherMessage{
				P1:             priority[1],
				P2:             priority[2],
				ActiveProblems: len(ActiveProblems),
			}
			w.Header().Set("Content-Type", "text/html")
			t.Execute(w, data)
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Fprint(w, "Unauthorized")
	}
}

//-----------------------------------------------------------------------------------
func student_messagesHandler(w http.ResponseWriter, r *http.Request) {
	stid, err := strconv.Atoi(r.FormValue("stid"))
	if err != nil {
		fmt.Fprintf(w, "Error")
	}
	mesg, ok := MessageBoards[stid]
	if ok {
		t := template.New("")
		t, err := t.Parse(STUDENT_MESSAGING_TEMPLATE)
		if err == nil {
			data := struct{ Message string }{mesg}
			w.Header().Set("Content-Type", "text/html")
			t.Execute(w, data)
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Fprint(w, "Error")
	}
}
