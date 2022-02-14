package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

type DashBoardStudentInfo struct {
	StudentID      int
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
	ProblemID          int
	NumActive          int
	NumHelpRequest     int
	NumGradedCorrect   int
	NumGradedIncorrect int
	NumNotGraded       int
	UserID             int
	UserRole           string
	Password           string
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

func getAllStudents() map[int]string {
	rows, err := Database.Query("select id, name from student")
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	var ID int
	var name string
	var students = make(map[int]string)
	for rows.Next() {
		rows.Scan(&ID, &name)
		students[ID] = name
	}
	return students
}

func getProblemStats(problemID int) (int, int, int, int, int) {
	rows, err := Database.Query("select active, submission, help_request, graded_correct, graded_incorrect from problem_statistics where problem_id = ?", problemID)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	var active, sub, help, correct, incorrect int
	if rows.Next() {
		rows.Scan(&active, &sub, &help, &correct, &incorrect)
	}
	return active, help, sub - correct - incorrect, correct, incorrect
}

func getProblemNameFromID(problemID int) string {
	rows, err := Database.Query("Select filename from problem where id = ?", problemID)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	var problemName string
	if rows.Next() {
		rows.Scan(&problemName)
	}
	rows.Close()
	return problemName
}

func problemDashboardHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	problemID, _ := strconv.Atoi(r.FormValue("problem_id"))
	role := r.FormValue("role")
	password := r.FormValue("password")
	students := getAllStudents()
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
			StudentID:      studentID,
			StudentName:    students[studentID],
			LastUpdatedAt:  lastUpdateMap[studentID],
			CodingStat:     codingStat,
			HelpStat:       helpStat,
			SubmissionStat: submissionStat,
			TutoringStat:   tutoringStat,
		})
	}
	rows.Close()

	nActive, nHelp, nNotGraded, nCorrect, nIncorrect := getProblemStats(problemID)

	dashBoardData := &DashBoardInfo{
		StudentInfo:        studentInfo,
		ProblemID:          problemID,
		ProblemName:        getProblemNameFromID(problemID),
		NumActive:          nActive,
		NumHelpRequest:     nHelp,
		NumGradedCorrect:   nCorrect,
		NumGradedIncorrect: nIncorrect,
		NumNotGraded:       nNotGraded,
		UserID:             uid,
		UserRole:           role,
		Password:           password,
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
