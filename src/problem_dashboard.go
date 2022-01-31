package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

type DashBoardStudentInfo struct {
	StudentName    string
	LastUpdatedAt  time.Time
	CodingStat     string
	HelpStat       string
	SubmissionStat string
	TutoringStat   string
}

type DashBoardInfo struct {
	StudentInfo        []*DashBoardStudentInfo
	ProblemName        string
	NumHelpRequest     int
	NumGradedCorrect   int
	NumGradedIncorrect int
	NumNotGraded       int
}

func getCurrentStudents() []int {
	rows, err := Database.Query("select student_id from attendance where DATE(attendance_at) = ?", time.Now().Format("2022-01-18"))
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	var studentID int
	var currentStudents []int
	for rows.Next() {
		rows.Scan(&studentID)
		currentStudents = append(currentStudents, studentID)
	}
	rows.Close()
	return currentStudents
}

func problemDashboardHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	problemID, _ := strconv.Atoi(r.FormValue("problem_id"))
	role := r.FormValue("role")
	if role != "teacher" {
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}
	rows, err := Database.Query("select student_id, max(last_updated_at) from code_snapshot where problem_id=? group by student_id", problemID)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	// currentStudents := getCurrentStudents()
	var lastUpdateMap = make(map[int]time.Time)

	var studentID int
	var lastUpdate time.Time
	for rows.Next() {
		rows.Scan(&studentID, &lastUpdate)
		lastUpdateMap[studentID] = lastUpdate
	}
	rows.Close()
	rows, err = Database.Query("select student_id, coding_stat, help_stat, submission_stat, tutoring_stat from student_status where problem_id=?", problemID)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	var codingStat, submissionStat, helpStat, tutoringStat string
	var studentInfo []*DashBoardStudentInfo
	for rows.Next() {
		rows.Scan(&studentID, &codingStat, &helpStat, &submissionStat, &tutoringStat)

		studentInfo = append(studentInfo, &DashBoardStudentInfo{
			StudentName:    Students[studentID].Name,
			LastUpdatedAt:  lastUpdateMap[studentID],
			CodingStat:     codingStat,
			HelpStat:       helpStat,
			SubmissionStat: submissionStat,
			TutoringStat:   tutoringStat,
		})
	}
	rows.Close()
	rows, err = Database.Query("Select count(*) from code_explanation where problem_id = ?", problemID)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	var numHelpReq int
	if rows.Next() {
		rows.Scan(&numHelpReq)
	}
	rows.Close()
	rows, err = Database.Query("Select count(*) from submission where problem_id = ? and verdict='correct'", problemID)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	var numGradedCorrect int
	if rows.Next() {
		rows.Scan(&numGradedCorrect)
	}
	rows.Close()
	rows, err = Database.Query("Select count(*) from submission where problem_id = ? and verdict='incorrect'", problemID)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	var numGradedIncorrect int
	if rows.Next() {
		rows.Scan(&numGradedIncorrect)
	}
	rows.Close()
	rows, err = Database.Query("Select count(*) from submission where problem_id = ? and verdict is NULL", problemID)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	var numNotGraded int
	if rows.Next() {
		rows.Scan(&numNotGraded)
	}
	rows.Close()
	rows, err = Database.Query("Select filename from problem where id = ?", problemID)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	var problemName string
	if rows.Next() {
		rows.Scan(&problemName)
	}
	rows.Close()
	dashBoardData := &DashBoardInfo{
		StudentInfo:        studentInfo,
		ProblemName:        problemName,
		NumHelpRequest:     numHelpReq,
		NumGradedCorrect:   numGradedCorrect,
		NumGradedIncorrect: numGradedIncorrect,
		NumNotGraded:       numNotGraded,
	}
	temp := template.New("")
	ownFuncs := template.FuncMap{"formatTimeSince": formatTimeSince}
	t, err := temp.Funcs(ownFuncs).Parse(PROBLEM_DASHBOARD_TEMPLATE)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, dashBoardData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
}
