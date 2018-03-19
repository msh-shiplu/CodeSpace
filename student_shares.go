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
func student_sharesHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	content, ext, priority := r.FormValue("content"), r.FormValue("ext"), r.FormValue("priority")
	pid, _ := strconv.Atoi(r.FormValue("pid"))
	sid := int64(0)
	if pid > 0 { // only keep in database submissions related to problems
		result, err := AddSubmissionSQL.Exec(pid, uid, content, priority, time.Now())
		if err != nil {
			panic(err)
		}
		sid, _ = result.LastInsertId()
	}

	SubSem.Lock()
	defer SubSem.Unlock()
	sub := &Submission{
		Sid:      int(sid),
		Uid:      uid,
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
