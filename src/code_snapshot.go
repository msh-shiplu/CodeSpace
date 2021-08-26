package main

import (
	"net/http"
	"strconv"
	"time"
)

func codeSnapshotHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	code := r.FormValue("code")
	problemID, _ := strconv.Atoi(r.FormValue("problem_id"))
	studentID, _ := strconv.Atoi(r.FormValue("uid"))
	addOrUpdateCodeSnapshot(studentID, problemID, -1, code, time.Now())
}
