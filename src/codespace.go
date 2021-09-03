package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type CodeSpaceData struct {
	Snapshots     []*Snapshot
	UserID        int
	UserRole      string
	Authenticated bool
}

type SnapshotData struct {
	Snapshot *Snapshot
	UserID   int
	UserRole string
}

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
	passcode := r.FormValue("pc")
	uid, _ := strconv.Atoi(r.FormValue("uid"))
	role := r.FormValue("role")
	temp := template.New("")
	t, err := temp.Parse(CODESPACE_TEMPLATE)
	if err != nil {
		log.Fatal(err)
	}
	var snapshots []*Snapshot
	if role == "student" {
		for _, s := range Snapshots {
			if _, ok := HelpEligibleStudents[s.ProblemID][uid]; ok {
				snapshots = append(snapshots, s)
			}
		}
	} else {
		snapshots = Snapshots
	}
	data := &CodeSpaceData{
		Snapshots:     snapshots,
		UserID:        uid,
		UserRole:      role,
		Authenticated: passcode == Passcode,
	}
	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, data)
}

func getCodeSnapshotHandler(w http.ResponseWriter, r *http.Request) {
	studentID, _ := strconv.Atoi(r.FormValue("student_id"))
	problemID, _ := strconv.Atoi(r.FormValue("problem_id"))
	uid, _ := strconv.Atoi(r.FormValue("uid"))
	role := r.FormValue("role")
	temp := template.New("")
	ownFuncs := template.FuncMap{"getEditorMode": getEditorMode}
	t, err := temp.Funcs(ownFuncs).Parse(CODE_SNAPSHOT_TEMPLATE)
	if err != nil {
		log.Fatal(err)
	}
	idx := StudentSnapshot[studentID][problemID]
	data := &SnapshotData{
		Snapshot: Snapshots[idx],
		UserID:   uid,
		UserRole: role,
	}
	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, data)
}
