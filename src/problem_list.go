package main

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

type ProblemData struct {
	ID                 int
	Filename           string
	UploadedAt         time.Time
	IsActive           bool
	Attendance         int
	NumActive          int
	NumHelpRequest     int
	NumGradedCorrect   int
	NumGradedIncorrect int
	NumNotGraded       int
}

type ProblemListData struct {
	Problems         []*ProblemData
	PeerTutorAllowed bool
	UserID           int
	UserRole         string
	Password         string
	Username         string
}

func problemListHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	role := r.FormValue("role")
	password := r.FormValue("password")
	rows, err := Database.Query("select id, filename, problem_uploaded_at, problem_ended_at from problem")
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	var problemID int
	var filename string
	var problemUploadedAt time.Time
	var problems = make([]*ProblemData, 0)

	var problemIDs = make([]int, 0)
	var filenameList = make([]string, 0)
	var uploadTimeList = make([]time.Time, 0)
	var endTimeList = make([]time.Time, 0)

	for rows.Next() {
		var problemEndedAt time.Time
		rows.Scan(&problemID, &filename, &problemUploadedAt, &problemEndedAt)
		problemIDs = append(problemIDs, problemID)
		filenameList = append(filenameList, filename)
		uploadTimeList = append(uploadTimeList, problemUploadedAt)
		endTimeList = append(endTimeList, problemEndedAt)
	}
	rows.Close()

	for i := 0; i < len(problemIDs); i++ {
		nActive, nHelp, nNotGraded, nCorrect, nIncorrect := getProblemStats(problemIDs[i])
		problems = append(problems, &ProblemData{
			ID:                 problemIDs[i],
			Filename:           filenameList[i],
			UploadedAt:         uploadTimeList[i],
			IsActive:           endTimeList[i].IsZero(),
			Attendance:         len(getCurrentStudents()),
			NumActive:          nActive,
			NumHelpRequest:     nHelp,
			NumGradedCorrect:   nCorrect,
			NumGradedIncorrect: nIncorrect,
			NumNotGraded:       nNotGraded,
		})
	}
	problemListData := &ProblemListData{
		Problems:         problems,
		PeerTutorAllowed: PeerTutorAllowed,
		UserID:           uid,
		UserRole:         role,
		Password:         password,
		Username:         getName(uid, role),
	}
	temp := template.New("")
	t, err := temp.Parse(PROBLEM_LIST_TEMPLATE)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, problemListData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
}
