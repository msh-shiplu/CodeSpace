//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"fmt"
	"net/http"
	"strconv"
	// "strings"
	"time"
)

//-----------------------------------------------------------------------------------
// func remove_header(content, ext string) string {
// 	prefix := "//"
// 	if ext != "java" && ext != "c++" && ext != "c" && ext != ".go" {
// 		prefix = "#"
// 	}
// 	content = strings.Trim(content, "\n ")
// 	if strings.HasPrefix(content, prefix) {
// 		items := strings.SplitN(content, "\n", 2)
// 		if len(items) > 1 {
// 			return items[1]
// 		}
// 	}
// 	return content
// }

//-----------------------------------------------------------------------------------
func teacher_gradesHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	content, ext, decision := r.FormValue("content"), r.FormValue("ext"), r.FormValue("decision")
	changed := r.FormValue("changed")
	pid, _ := strconv.Atoi(r.FormValue("pid"))
	stid, _ := strconv.Atoi(r.FormValue("stid"))
	mesg, student_mesg := "", ""
	if decision == "dismiss" {
		MessageBoards[stid] = "Teacher looked at but did not grade your submission."
		fmt.Fprintf(w, "Submission dismissed.")
		return
	}
	if changed == "True" {
		AddFeedbackSQL.Exec(uid, stid, content, time.Now())
		mesg = "Feedback saved to student's board."
		student_mesg += "You've got feedback."
		BoardsSem.Lock()
		defer BoardsSem.Unlock()
		b := &Board{
			Content:      content,
			Answer:       "",
			Attempts:     100,
			Ext:          ext,
			Pid:          pid,
			StartingTime: time.Now(),
		}
		Boards[stid] = append(Boards[stid], b)
	}
	score_id, current_points, current_attempts := 0, 0, 0
	// change this to query by score id
	rows, _ := Database.Query("select id, points, attempts from score where pid=? and stid=?", pid, stid)
	for rows.Next() {
		rows.Scan(&score_id, &current_points, &current_attempts)
		break
	}
	rows.Close()
	merit, effort := 0, 0
	rows, _ = Database.Query("select merit, effort from problem where id=?", pid)
	for rows.Next() {
		rows.Scan(&merit, &effort)
		break
	}
	rows.Close()
	if decision == "correct" {
		if score_id == 0 {
			_, err := AddScoreSQL.Exec(pid, stid, uid, merit, 1)
			if err != nil {
				panic(err)
			}
		} else {
			_, err := UpdateScoreSQL.Exec(uid, merit, current_attempts+1, score_id)
			if err != nil {
				panic(err)
			}
		}
		mesg = "Problem graded correct.\n" + mesg
		student_mesg += " Your submission was correct."
		next_pid, ok := NextProblem[int64(pid)]
		if ok {
			new_content, new_answer, new_ext, new_merit, new_effort, new_attempts := "", "", "", 0, 0, 0
			rows, _ = Database.Query("select content, answer, ext, merit, effort, attempts from problem where id=?", next_pid)
			for rows.Next() {
				rows.Scan(&new_content, &new_answer, &new_ext, &new_merit, &new_effort, &new_attempts)
				break
			}
			rows.Close()

			// this is done to be consistent.
			// new_content = remove_header(new_content, new_ext)

			b := &Board{
				Content:      new_content,
				Answer:       new_answer,
				Attempts:     new_attempts,
				Ext:          new_ext,
				Pid:          int(next_pid),
				StartingTime: time.Now(),
			}
			Boards[stid] = append(Boards[stid], b)
			mesg = "Next problem added to student's board\n" + mesg
			student_mesg += " A new problem is added to your board."
		}
	} else {
		if score_id == 0 {
			_, err := AddScoreSQL.Exec(pid, stid, uid, effort, 1)
			if err != nil {
				panic(err)
			}
		} else {
			_, err := UpdateScoreSQL.Exec(uid, current_points, current_attempts+1, score_id)
			if err != nil {
				panic(err)
			}
		}
		mesg = "Problem graded incorrect.\n" + mesg
		student_mesg += " Your submssion was not correct. Try again."
	}
	MessageBoards[stid] = student_mesg
	fmt.Fprintf(w, mesg)
}

//-----------------------------------------------------------------------------------
