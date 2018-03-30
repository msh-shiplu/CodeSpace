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
	Points   int
	Filename string
	Date     int64
}

//-----------------------------------------------------------------------------------
func student_gets_reportHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	rows, err := Database.Query("select score.points, score.at, problem.filename from score join problem on problem.id == score.pid where stid=?", uid)
	defer rows.Close()
	if err != nil {
		panic(err)
	}
	report := make([]*StudentReport, 0)
	var points int
	var t time.Time
	var filename string
	for rows.Next() {
		// fmt.Println(rows)
		rows.Scan(&points, &t, &filename)
		report = append(report, &StudentReport{Points: points, Filename: filename, Date: t.Unix()})
	}
	js, _ := json.Marshal(report)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

//-----------------------------------------------------------------------------------
