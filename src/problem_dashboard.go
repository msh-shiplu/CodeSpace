package main

import (
	"html/template"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
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

type AnswerStatInfo struct {
	Answer     string
	Count      int
	Percent    float64
	NumStudent int
}

type StudentStatInfo struct {
	NumActive       int
	NumWorking      int
	NumSubmitted    int
	NumWaitingGrade int
	NumWaitingHelp  int
}

type SubmissionStatInfo struct {
	NumSubmission      int
	NumHelpRequest     int
	NumGraded          int
	NumNotGraded       int
	NumGradedCorrect   int
	NumGradedIncorrect int
}

type DashBoardInfo struct {
	StudentInfo    []*DashBoardStudentInfo
	ProblemName    string
	Code           string
	IsActive       bool
	ProblemID      int
	StudentStat    *StudentStatInfo
	SubmissionStat *SubmissionStatInfo
	AnswerStats    []*AnswerStatInfo
	CorrectAnswer  string
	UserID         int
	UserRole       string
	Password       string
	Username       string
}

func getName(uid int, role string) string {
	name := ""
	if role == "teacher" {
		name = TeacherIdToName[uid]
	} else {
		name = Students[uid].Name
	}
	return name
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
	rows.Close()
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
	rows.Close()
	return active, help, sub - correct - incorrect, correct, incorrect
}

func getStudentAndSubmissionStatsForProblem(problemID int) (StudentStatInfo, SubmissionStatInfo) {
	rows, err := Database.Query("select active, working, submission, help_request, graded_correct, graded_incorrect from problem_statistics where problem_id = ?", problemID)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	var studentStat StudentStatInfo
	var subStat SubmissionStatInfo
	if rows.Next() {
		rows.Scan(&studentStat.NumActive, &studentStat.NumWorking, &subStat.NumSubmission, &subStat.NumHelpRequest, &subStat.NumGradedCorrect, &subStat.NumGradedIncorrect)
	}
	rows.Close()
	subStat.NumNotGraded = subStat.NumSubmission - subStat.NumGradedCorrect - subStat.NumGradedIncorrect
	subStat.NumGraded = subStat.NumSubmission - subStat.NumNotGraded

	rows, err = Database.Query("select count(distinct student_id) from submission where problem_id=?", problemID)
	if err != nil {
		log.Fatal(err)
	}
	if rows.Next() {
		rows.Scan(&studentStat.NumSubmitted)
	}
	rows.Close()

	rows, err = Database.Query("select count(distinct student_id) from submission where problem_id=? and verdict is NULL", problemID)
	if err != nil {
		log.Fatal(err)
	}
	if rows.Next() {
		rows.Scan(&studentStat.NumWaitingGrade)
	}
	rows.Close()

	rows, err = Database.Query("select count(distinct c.student_id) from message m, code_snapshot c where m.snapshot_id=c.id and c.problem_id=? and m.id not in (select message_id from message_feedback)", problemID)
	if err != nil {
		log.Fatal(err)
	}
	if rows.Next() {
		rows.Scan(&studentStat.NumWaitingHelp)
	}
	rows.Close()

	return studentStat, subStat
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

func getLatestSubmissionTime(problemID int) map[int]time.Time {
	var latestSubmissions = make(map[int]time.Time)
	rows, err := Database.Query("select student_id, max(code_submitted_at) from submission where problem_id=? group by student_id", problemID)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	var studentID int
	var submissionTimeStr string
	var submissionTime time.Time
	layout := "2006-01-02 15:04:05-07:00"
	for rows.Next() {
		rows.Scan(&studentID, &submissionTimeStr)
		submissionTime, _ = time.Parse(layout, submissionTimeStr)
		latestSubmissions[studentID] = submissionTime
	}
	rows.Close()
	return latestSubmissions
}

func getNumStudentForAnswer(answer string) int {
	rows, err := Database.Query("select count(DISTINCT student_id) as cnt from submission where answer is not NULL and answer = ?", answer)
	numStudent := 0
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	if rows.Next() {
		rows.Scan(&numStudent)
	}
	rows.Close()
	return numStudent
}

func problemDashboardHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	problemID, _ := strconv.Atoi(r.FormValue("problem_id"))
	role := r.FormValue("role")
	password := r.FormValue("password")
	students := getAllStudents()
	rows, err := Database.Query("select student_id, max(last_updated_at) from code_snapshot where problem_id=? group by student_id", problemID)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	// currentStudents := getCurrentStudents()
	var lastUpdateMap = make(map[int]time.Time)

	var studentID int
	var lastUpdateString string
	var lastUpdate time.Time
	layout := "2006-01-02 15:04:05-07:00"

	_, ok := HelpEligibleStudents[problemID][uid]
	for rows.Next() {
		rows.Scan(&studentID, &lastUpdateString)
		if role == "teacher" || uid == studentID || (PeerTutorAllowed && ok) {
			lastUpdate, _ = time.Parse(layout, lastUpdateString)
			lastUpdateMap[studentID] = lastUpdate
		}
	}
	rows.Close()
	rows, err = Database.Query("select filename, problem_description, answer, problem_ended_at from problem where id=?", problemID)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	var code string
	var problemEndedAt time.Time
	var filename string
	var answer string
	if rows.Next() {
		rows.Scan(&filename, &code, &answer, &problemEndedAt)
	}
	rows.Close()
	latestSubmissionTime := getLatestSubmissionTime(problemID)
	rows, err = Database.Query("select student_id, coding_stat, help_stat, submission_stat, tutoring_stat from student_status where problem_id=?", problemID)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	var codingStat, submissionStat, helpStat, tutoringStat string
	var studentInfo []*DashBoardStudentInfo
	for rows.Next() {
		rows.Scan(&studentID, &codingStat, &helpStat, &submissionStat, &tutoringStat)
		if role == "teacher" || uid == studentID || (PeerTutorAllowed && ok) {
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
	}
	rows.Close()

	sort.SliceStable(studentInfo, func(i, j int) bool {
		if studentInfo[i].SubmissionStat == "submitted" && studentInfo[j].SubmissionStat == "submitted" {
			return latestSubmissionTime[studentInfo[i].StudentID].Before(latestSubmissionTime[studentInfo[j].StudentID])
		}
		if studentInfo[i].SubmissionStat == "submitted" {
			return true
		}
		if studentInfo[j].SubmissionStat == "submitted" {
			return false
		}
		if studentInfo[i].HelpStat == "Asked for help" {
			return true
		}
		if studentInfo[j].HelpStat == "Asked for help" {
			return false
		}
		return true
	})

	stStat, subStat := getStudentAndSubmissionStatsForProblem(problemID)

	var answerStats []*AnswerStatInfo
	if role != "student" {
		rows, err = Database.Query("select answer, count(*) as cnt from submission where problem_id = ? and answer is not NULL and LENGTH(answer)>0 group by answer", problemID)
		defer rows.Close()
		if err != nil {
			log.Fatal(err)
		}
		var ans string
		var c int
		var total int
		for rows.Next() {
			rows.Scan(&ans, &c)
			// ans = strings.TrimSpace(ans)
			if ans != "" {
				answerStats = append(answerStats, &AnswerStatInfo{
					Answer: ans,
					Count:  c,
				})
			}
			total += c
		}
		rows.Close()
		for i, answer := range answerStats {
			answerStats[i].Percent = float64(answer.Count) * 100.0 / float64(total)
			answerStats[i].Percent = math.Round(answerStats[i].Percent*100) / 100

			answerStats[i].NumStudent = getNumStudentForAnswer(answer.Answer)
		}
	}

	dashBoardData := &DashBoardInfo{
		StudentInfo:    studentInfo,
		ProblemID:      problemID,
		ProblemName:    filename,
		Code:           code,
		IsActive:       problemEndedAt.IsZero(),
		StudentStat:    &stStat,
		SubmissionStat: &subStat,
		AnswerStats:    answerStats,
		CorrectAnswer:  strings.TrimSpace(answer),
		UserID:         uid,
		UserRole:       role,
		Password:       password,
		Username:       getName(uid, role),
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
