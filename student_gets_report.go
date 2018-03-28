//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"encoding/json"
	// "fmt"
	"net/http"
	"time"
)

type StudentReport struct {
	Points int
	Date   int64
}

//-----------------------------------------------------------------------------------
func student_gets_reportHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	rows, _ := Database.Query("select points, at from score where stid=?", uid)
	defer rows.Close()
	report := make([]*StudentReport, 0)
	var points int
	var t time.Time
	for rows.Next() {
		rows.Scan(&points, &t)
		report = append(report, &StudentReport{Points: points, Date: t.Unix()})
	}
	js, err := json.Marshal(report)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

//-----------------------------------------------------------------------------------
