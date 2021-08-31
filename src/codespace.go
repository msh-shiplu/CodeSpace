package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func codespaceHandler(w http.ResponseWriter, r *http.Request) {
	temp := template.New("")
	t, err := temp.Parse(CODESPACE_TEMPLATE)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, Snapshots)
}

func getCodeSnapshotHandler(w http.ResponseWriter, r *http.Request) {
	studentID, _ := strconv.Atoi(r.FormValue("student_id"))
	problemID, _ := strconv.Atoi(r.FormValue("problem_id"))
	temp := template.New("")
	t, err := temp.Parse(CODE_SNAPSHOT_TEMPLATE)
	if err != nil {
		log.Fatal(err)
	}
	idx := StudentSnapshot[studentID][problemID]
	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, Snapshots[idx])
}
