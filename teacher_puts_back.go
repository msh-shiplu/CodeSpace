//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

//-----------------------------------------------------------------------------------
func teacher_puts_backHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	content, ext := r.FormValue("content"), r.FormValue("ext")
	priority, _ := strconv.Atoi(r.FormValue("priority"))
	pid, _ := strconv.Atoi(r.FormValue("pid"))
	sid, _ := strconv.Atoi(r.FormValue("sid"))
	stid, _ := strconv.Atoi(r.FormValue("stid"))

	SubSem.Lock()
	defer SubSem.Unlock()
	sub := &Submission{
		Sid:      sid,
		Uid:      stid,
		Pid:      pid,
		Content:  content,
		Ext:      ext,
		Priority: priority,
		At:       time.Now(),
	}
	WorkingSubs = append(WorkingSubs, sub)
	fmt.Fprintf(w, "OK")
}

//-----------------------------------------------------------------------------------
