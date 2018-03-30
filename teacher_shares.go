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
	for stid, _ := range Boards {
		b := &Board{
			Content:      content,
			Answer:       "",
			Attempts:     -1,
			Filename:     filename,
			Pid:          0,
			StartingTime: time.Now(),
		}
		Boards[stid] = append(Boards[stid], b)
		MessageBoards[stid] = "Teacher shared some material with you."
	}
	fmt.Fprintf(w, "Content shared to all students.")

}
