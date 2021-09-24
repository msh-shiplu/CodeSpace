package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type CodeSpaceData struct {
	Snapshots     []*Snapshot
	UserID        int
	UserRole      string
	Authenticated bool
	Passcode      string
}

type FeedbackData struct {
	Feedback     string
	FeedbackTime time.Time
}
type SnapshotData struct {
	Snapshot  *Snapshot
	UserID    int
	UserRole  string
	Feedbacks []*FeedbackData
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

func formatTimeDuration(d time.Duration) string {
	m := int(d.Minutes())
	d1 := d - time.Duration(m*60*1000000000)
	s := int(d1.Seconds())
	str := ""
	if m > 0 {
		str = strconv.Itoa(m) + " minute(s) "
	}
	str += strconv.Itoa(s) + " second(s)"
	return str
}

func formatTimeSince(t time.Time) string {
	d := time.Now().Sub(t)
	return formatTimeDuration(d)
}

func codespaceHandler(w http.ResponseWriter, r *http.Request) {
	passcode := r.FormValue("pc")
	uid, _ := strconv.Atoi(r.FormValue("uid"))
	role := r.FormValue("role")
	temp := template.New("")
	ownFuncs := template.FuncMap{"formatTimeSince": formatTimeSince}
	t, err := temp.Funcs(ownFuncs).Parse(CODESPACE_TEMPLATE)
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
		Passcode:      passcode,
	}
	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
}

func getCodeSnapshotHandler(w http.ResponseWriter, r *http.Request) {
	studentID, _ := strconv.Atoi(r.FormValue("student_id"))
	problemID, _ := strconv.Atoi(r.FormValue("problem_id"))
	uid, _ := strconv.Atoi(r.FormValue("uid"))
	role := r.FormValue("role")
	passcode := r.FormValue("pc")
	if passcode != Passcode {
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}
	temp := template.New("")
	ownFuncs := template.FuncMap{"getEditorMode": getEditorMode}
	t, err := temp.Funcs(ownFuncs).Parse(CODE_SNAPSHOT_TEMPLATE)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := Database.Query("select F.feedback, F.given_at from code_snapshot C, snapshot_feedback F where C.id=F.snapshot_id and C.student_id=? and C.problem_id=?", studentID, problemID)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	feedback, givenAt := "", time.Now()
	var feedbacks []*FeedbackData

	for rows.Next() {
		rows.Scan(&feedback, &givenAt)
		feedbacks = append(feedbacks, &FeedbackData{Feedback: feedback, FeedbackTime: givenAt})
	}
	idx := StudentSnapshot[studentID][problemID]
	data := &SnapshotData{
		Snapshot:  Snapshots[idx],
		UserID:    uid,
		UserRole:  role,
		Feedbacks: feedbacks,
	}
	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
}
