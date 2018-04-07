//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"fmt"
	"net/http"
	"time"
)

//-----------------------------------------------------------------------------------
func teacher_sharesHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	content, filename := r.FormValue("content"), r.FormValue("filename")
	for stid := range Students {
		b := &Board{
			Content:      content,
			Answer:       "",
			Attempts:     -1,
			Filename:     filename,
			Pid:          0,
			StartingTime: time.Now(),
		}
		Students[stid].Boards = append(Students[stid].Boards, b)
	}
	fmt.Fprintf(w, "Content shared to all students.")

}
