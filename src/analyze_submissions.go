//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"net/http"
	"sort"
	"time"
)

//-----------------------------------------------------------------------------------
type SubmissionData struct {
	Flag      string
	Start     int64
	At        int64
	Completed int64
}

//-----------------------------------------------------------------------------------
func analyze_submissionsHandler(w http.ResponseWriter, r *http.Request) {
	// if r.FormValue("pc") != Passcode {
	// 	fmt.Fprintf(w, "Unauthorized")
	// 	return
	// }
	pid := r.FormValue("pid")
	records := make(map[int][]*SubmissionData)
	var sid, priority int
	var start, at, completed time.Time

	row, _ := Database.Query("select at from problem where id=?", pid)
	for row.Next() {
		row.Scan(&start)
	}
	row.Close()

	rows, _ := Database.Query("select sid, priority, at, completed from submission where pid=?", pid)
	for rows.Next() {
		rows.Scan(&sid, &priority, &at, &completed)
		if _, ok := records[sid]; !ok {
			records[sid] = make([]*SubmissionData, 0)
		}
		flag := "unknown"
		if priority == 1 {
			flag = "Got it!"
		} else if priority == 2 {
			flag = "Help!"
		}
		records[sid] = append(
			records[sid],
			&SubmissionData{
				Flag:      flag,
				Start:     start.UnixNano(),
				At:        at.UnixNano(),
				Completed: completed.UnixNano(),
			})
	}
	for sid, _ := range records {
		sort.Slice(records[sid], func(i, j int) bool {
			return records[sid][i].At < records[sid][j].At
		})
	}
	rows.Close()
	w.Header().Set("Content-Type", "text/html")
	t, _ := template.New("").Parse(ANALYZE_SUBMISSIONS_TEMPLATE)
	err := t.Execute(w, records)
	if err != nil {
		fmt.Println(err)
	}
}
