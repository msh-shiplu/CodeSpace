//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	// "strings"
	// "time"
)

//-----------------------------------------------------------------------------------
func student_shareHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	content, ext, level := r.FormValue("content"), r.FormValue("ext"), r.FormValue("level")
	pid, _ := strconv.Atoi(r.FormValue("pid"))
	fmt.Println(pid, ext, level, "\n", content)
	fmt.Fprintf(w, "OK")
}

//-----------------------------------------------------------------------------------
func student_get_boardcontentHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	var js []byte
	var err error

	BoardsSem.Lock()
	defer BoardsSem.Unlock()

	if board, ok := Boards[uid]; ok {
		js, err = json.Marshal(board)
		Boards[uid] = []*Board{}
		if err == nil {
			fmt.Println(string(js))
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			return
		}
	}
	fmt.Println(err.Error())
	js, err = json.Marshal([]*Board{})
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
