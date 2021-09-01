package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func getEditorMode(filename string) string {
	filename = strings.ToLower(filename)
	if strings.HasSuffix(filename, ".py") {
		return "python"
	}
	if strings.HasSuffix(filename, ".java") {
		return "text/x-java"
	}
	if strings.HasSuffix(filename, ".cpp") || strings.HasSuffix(filename, ".c++") || strings.HasSuffix(filename, ".c") {
		return "text/x-c++src"
	}
	return "text"
}
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
	ownFuncs := template.FuncMap{"getEditorMode": getEditorMode}
	t, err := temp.Funcs(ownFuncs).Parse(CODE_SNAPSHOT_TEMPLATE)
	if err != nil {
		log.Fatal(err)
	}
	idx := StudentSnapshot[studentID][problemID]
	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, Snapshots[idx])
}
