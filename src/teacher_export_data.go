package main

import (
	"encoding/csv"
	"log"
	"net/http"
	"os"
	"strconv"
)

func exportPointsHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	role := r.FormValue("role")
	if role == "student" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	rows, err := Database.Query("select student_id, filename, score from problem P, Score S where P.id = S.problem_id")
	if err != nil {
		log.Fatal(err)
	}
	var studentID, score int
	var filename string
	var data = make(map[string]map[int]int)
	for rows.Next() {
		rows.Scan(&studentID, &filename, &score)
		if data[filename] == nil {
			data[filename] = make(map[int]int)
		}
		data[filename][studentID] = score
	}
	csvData := make([][]string, 1)
	var i = 1
	csvData[0] = make([]string, len(data)+1)
	for k := range data {
		csvData[0][i] = k
		i++
	}
	for sid := range Students {
		dt := make([]string, len(data)+1)
		dt[0] = getName(sid, "student")
		i = 1
		for _, v := range data {
			s, ok := v[sid]
			if ok {
				dt[i] = strconv.Itoa(s)
			}
			i++
		}
		csvData = append(csvData, dt)
	}
	file, err := os.Create("score.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)

	defer writer.Flush()

	writer.WriteAll(csvData)

	w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote("score.csv"))
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, file.Name())
}
