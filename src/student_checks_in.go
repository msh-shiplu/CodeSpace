//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

//-----------------------------------------------------------------

func student_checks_inHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	// attendance is taken automatically by authorization when this handler is called.
	// Next: return student attendance report
	rows, err := Database.Query("select attendance_at from attendance where student_id=?", uid)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	dates := make([]int64, 0)
	var t time.Time
	for rows.Next() {
		rows.Scan(&t)
		dates = append(dates, t.Unix())
	}
	js, _ := json.Marshal(dates)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

//-----------------------------------------------------------------
func student_periodic_updateHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	submission_stat := Students[uid].SubmissionStatus
	board_stat := 0
	if len(Students[uid].Boards) > 0 {
		board_stat = 1
	}
	Students[uid].SubmissionStatus = 0 // reset status after notifying student
	fmt.Fprintf(w, "%d;%d", submission_stat, board_stat)
}
