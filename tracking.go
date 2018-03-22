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

//-----------------------------------------------------------------------------------
func student_trackingHandler(w http.ResponseWriter, r *http.Request) {
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
