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
	need_help_with := r.FormValue("need_help_with")
	sid := int64(0)

	var err error
	msg := "your help message has been sent"

	pid := 0
	prob, ok := ActiveProblems[filename]
	snapshotID := 0
	if ok {
		if !prob.Active {
			msg = "Problem is no longer active. But the teacher will look at your submission."
		} else {
			pid = prob.Info.Pid
			if _, ok := prob.Attempts[uid]; !ok {
				ActiveProblems[filename].Attempts[uid] = prob.Info.Attempts
			}
			now := time.Now()
			snapshotID = addCodeSnapshot(uid, pid, content, 0, now)
			var result sql.Result
			result, err = AddHelpSubmissionSQL.Exec(pid, uid, snapshotID, "", need_help_with, now)

			if err != nil {
				log.Fatal(err)
			}
			sid, _ = result.LastInsertId()
		}
	} else {
		msg = "Invalid filename"
	}
	if ok && prob.Active {
		HelpSubSem.Lock()
		defer HelpSubSem.Unlock()
		sub := &HelpSubmission{
			Sid:        int(sid),
			Uid:        uid,
			Pid:        pid,
			Content:    content,
			Filename:   filename,
			At:         time.Now(),
			SnapshotID: snapshotID,
		}
		WorkingHelpSubs = append(WorkingHelpSubs, sub)
		HelpSubmissions[int(sid)] = sub
	}

	fmt.Fprintf(w, msg)

}

//-----------------------------------------------------------------------------------
