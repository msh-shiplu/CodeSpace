//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"
)

//-----------------------------------------------------------------------------------
func studentAskHelpHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	content, filename := r.FormValue("content"), r.FormValue("filename")
	sid := int64(0)

	var err error
	msg := "Your help request will be looked at soon."

	pid := 0
	prob, ok := ActiveProblems[filename]

	if ok {
		if !prob.Active {
			msg = "Problem is no longer active. But the teacher will look at your submission."
		} else {
			pid = prob.Info.Pid
			if _, ok := prob.Attempts[uid]; !ok {
				ActiveProblems[filename].Attempts[uid] = prob.Info.Attempts
			}
			var result sql.Result
			result, err = AddHelpSubmissionSQL.Exec(pid, uid, content, time.Now())

			if err != nil {
				log.Fatal(err)
			}
			sid, _ = result.LastInsertId()
		}
	}

	HelpSubSem.Lock()
	defer HelpSubSem.Unlock()
	sub := &HelpSubmission{
		Sid:      int(sid),
		Uid:      uid,
		Pid:      pid,
		Content:  content,
		Filename: filename,
		At:       time.Now(),
	}
	WorkingHelpSubs = append(WorkingHelpSubs, sub)
	HelpSubmissions[int(sid)] = sub
	fmt.Fprintf(w, msg)

}

//-----------------------------------------------------------------------------------
