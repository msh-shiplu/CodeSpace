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
	if need_help_with == "" {
		need_help_with = "None."
	}

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
			snapshotID = addCodeSnapshot(uid, pid, content, 0, now, "at_ask_for_help")
			var result sql.Result
			// result, err = AddHelpSubmissionSQL.Exec(pid, uid, snapshotID, "", need_help_with, now)
			result, err = AddMessageSQL.Exec(snapshotID, need_help_with, uid, "student", now, 0)
			if err != nil {
				log.Fatal(err)
			}
			sid, _ = result.LastInsertId()
			_, err = IncProblemStatHelpSQL.Exec(pid)
			if err != nil {
				log.Fatal(err)
			}
			addOrUpdateStudentStatus(uid, pid, "", "Asked for help", "", "")
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
			Content:    need_help_with,
			Filename:   filename,
			At:         time.Now(),
			SnapshotID: snapshotID,
			Snapshot:   content,
		}
		WorkingHelpSubs = append(WorkingHelpSubs, sub)
		HelpSubmissions[int(sid)] = sub
	}

	fmt.Fprintf(w, msg)

}

//-----------------------------------------------------------------------------------
