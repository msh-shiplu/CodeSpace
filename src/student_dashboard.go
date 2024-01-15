package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type FeedbackDashBaord struct {
	Name            string
	Role            string
	Feedback        string
	FeedbackID      int
	CurrentUserVote string
	Downvote        int
	Upvote          int
	GivenAt         time.Time
}

type MessageDashBoard struct {
	ID         int
	Name       string
	Role       string
	Message    string
	Type       int // 0 = help request, 1 = unsolicited
	Event      string
	GivenAt    time.Time
	Code       string
	SnapshotID int
	Feedbacks  []*FeedbackDashBaord
}

type FeedbackProvisionDashBoard struct {
	StudentName  string
	ProblemName  string
	Status       DashBoardStudentInfo
	LastSnapshot *Snapshot
	Messages     []*MessageDashBoard
	StudentID    int
	ProblemID    int
	UserID       int
	UserRole     string
	Password     string
	Username     string
}

type SubmissionInfo struct {
	ID          int
	Code        string
	Grade       string
	SubmittedAt time.Time
	SnapshotID  int
}

type SubmissionDashboard struct {
	Submissions []*SubmissionInfo
	StudentName string
	ProblemName string
	StudentID   int
	ProblemID   int
	UserID      int
	UserRole    string
	Password    string
	Username    string
}

type TemplateDate struct {
	Feedback   	   FeedbackProvisionDashBoard
	Submission     SubmissionDashboard
	Status         DashBoardStudentInfo
	ChatgptaServer string
	UserID         int
	UserRole       string
	Password       string
	Username       string
	CourseName	   string
}

func getCurrentUserVote(feedbackID int, userID int, userRole string) string {
	var vote string
	row, err := Database.Query("select useful from message_back_feedback where message_feedback_id=? and author_id=? and author_role=?", feedbackID, userID, userRole)
	defer row.Close()
	if err != nil {
		log.Fatal(err)
	}
	for row.Next() {
		row.Scan(&vote)
	}
	row.Close()
	return vote
}

func getMessageFeedbacks(messageID int, userID int, userRole string) []*FeedbackDashBaord {
	rows, err := Database.Query("select id, feedback, author_id, author_role, given_at from message_feedback where message_id = ?", messageID)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	var feedback, authorRole string
	var authorID, feedbackID int
	var givenAt time.Time

	feedbacks := make([]*FeedbackDashBaord, 0)

	for rows.Next() {
		rows.Scan(&feedbackID, &feedback, &authorID, &authorRole, &givenAt)
		name := ""
		if authorRole == "teacher" {
			name = getTeacherName(authorID)
		} else {
			name = getStudentName(authorID)
		}
		feedbacks = append(feedbacks, &FeedbackDashBaord{
			Name:            name,
			Role:            authorRole,
			Feedback:        feedback,
			FeedbackID:      feedbackID,
			CurrentUserVote: getCurrentUserVote(feedbackID, userID, userRole),
			Downvote:        getBackFeedbackCount(feedbackID, "no"),
			Upvote:          getBackFeedbackCount(feedbackID, "yes"),
			GivenAt:         givenAt,
		})
	}
	rows.Close()
	return feedbacks
}

func getLatestSnapshot(studentID int, problemID int) *Snapshot {
	rows, err := Database.Query("select id, code, max(last_updated_at) from code_snapshot where problem_id = ? and student_id=? group by problem_id, student_id", problemID, studentID)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	var ID int
	var code string
	var lastUpdate time.Time
	if rows.Next() {
		rows.Scan(&ID, &code, &lastUpdate)
	}
	rows.Close()
	return &Snapshot{
		ID:          ID,
		ProblemName: getProblemNameFromID(problemID),
		Code:        code,
		LastUpdated: lastUpdate,
	}
}

func getTeacherName(authorID int) string {
	rows, err := Database.Query("select name from teacher where id = ?", authorID)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	var name string
	if rows.Next() {
		rows.Scan(&name)
	}
	return name
}

func getStudentName(studentID int) string {
	rows, err := Database.Query("select name from student where id = ?", studentID)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	var name string
	if rows.Next() {
		rows.Scan(&name)
	}
	return name
}

func getBackFeedbackCount(feedbackID int, backFeedbackType string) int {
	vote := 0
	rows, err := Database.Query("select count(*) from message_back_feedback where useful = ? and message_feedback_id = ?", backFeedbackType, feedbackID)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		rows.Scan(&vote)
	}
	return vote
}

func studentDashboardFeedbackProvisionHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	role := r.FormValue("role")
	problemID, _ := strconv.Atoi(r.FormValue("problem_id"))
	studentID, _ := strconv.Atoi(r.FormValue("student_id"))
	temp := template.New("")
	ownFuncs := template.FuncMap{"getEditorMode": getEditorMode}
	t, err := temp.Funcs(ownFuncs).Parse(FEEDBACK_PROVISION_TEMPLATE)
	if err != nil {
		log.Fatal(err)
	}
	students := getAllStudents()
	var messages = make([]*MessageDashBoard, 0)
	_, ok := HelpEligibleStudents[problemID][uid]
	if role == "teacher" || uid == studentID || (PeerTutorAllowed && ok) {
		rows, err := Database.Query("select M.id, M.snapshot_id, M.message, M.author_id, M.author_role, M.given_at, M.type, C.Code, C.event from message M, code_snapshot C where M.snapshot_id = C.id and C.problem_id = ? and C.student_id = ?", problemID, studentID)
		defer rows.Close()
		if err != nil {
			log.Fatal(err)
		}
		var snapshotID, authorID, messageType, messageID int
		var message, authorRole, code, event string
		var givenAt time.Time
		for rows.Next() {
			rows.Scan(&messageID, &snapshotID, &message, &authorID, &authorRole, &givenAt, &messageType, &code, &event)
			name := ""
			if authorRole == "teacher" {
				name = getTeacherName(authorID)
			} else {
				name = students[authorID]
			}
			messages = append(messages, &MessageDashBoard{
				ID:         messageID,
				Name:       name,
				Role:       authorRole,
				Message:    message,
				Type:       messageType,
				Event:      event,
				GivenAt:    givenAt,
				SnapshotID: snapshotID,
				Code:       code,
				Feedbacks:  getMessageFeedbacks(messageID, uid, role),
			})
		}
	} else {
		http.Error(w, "You are not authorized to access!", http.StatusUnauthorized)
	}

	// Sort the messages by descding order of GivenAt
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].GivenAt.After(messages[j].GivenAt)
	})

	// TODO(shiplu): sort the messages array
	// sort.Slice(helpRequests, func(i, j int) bool { return helpRequests[i].GivenAt.Before(helpRequests[j].GivenAt) })
	latestSnapshot := &Snapshot{}
	if _, ok := StudentSnapshot[studentID][problemID]; ok {
		latestSnapshot = Snapshots[StudentSnapshot[studentID][problemID]]
	} else {
		latestSnapshot = getLatestSnapshot(studentID, problemID)
	}

	// Get student status
	rows, err := Database.Query("select coding_stat, help_stat, submission_stat, tutoring_stat from student_status where problem_id=? and student_id=?", problemID, studentID)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	var codingStat, submissionStat, helpStat, tutoringStat string
	// var studentInfo []*DashBoardStudentInfo
	studentStats := DashBoardStudentInfo{}
	for rows.Next() {
		rows.Scan(&codingStat, &helpStat, &submissionStat, &tutoringStat)
		if role == "teacher" || uid == studentID || (PeerTutorAllowed && ok) {
			studentStats.CodingStat = codingStat
			studentStats.HelpStat = helpStat
			studentStats.SubmissionStat = submissionStat
			studentStats.TutoringStat = tutoringStat
		}
	}
	rows.Close()

	data := &FeedbackProvisionDashBoard{
		StudentName:  students[studentID],
		ProblemName:  latestSnapshot.ProblemName,
		LastSnapshot: latestSnapshot,
		Messages:     messages,
		StudentID:    studentID,
		ProblemID:    problemID,
		UserID:       uid,
		UserRole:     role,
		Password:     r.FormValue("password"),
		Status:       studentStats,
		Username:     getName(uid, role),
	}
	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
}

func studentDashboardSubmissionHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	role := r.FormValue("role")
	problemID, _ := strconv.Atoi(r.FormValue("problem_id"))
	studentID, _ := strconv.Atoi(r.FormValue("student_id"))
	temp := template.New("")
	ownFuncs := template.FuncMap{"getEditorMode": getEditorMode}
	t, err := temp.Funcs(ownFuncs).Parse(SUBMISSION_VIEW_TEMPLATE)
	if err != nil {
		log.Fatal(err)
	}

	var submissions = make([]*SubmissionInfo, 0)
	_, ok := HelpEligibleStudents[problemID][uid]
	if role == "teacher" || uid == studentID || (PeerTutorAllowed && ok) {
		rows, err := Database.Query("select id, snapshot_id, student_code, code_submitted_at, verdict from submission where student_id = ? and problem_id = ?", studentID, problemID)
		defer rows.Close()
		if err != nil {
			log.Fatal(err)
		}
		var snapshotID, submissionID int
		var verdict, code string
		var submittedAt time.Time
		for rows.Next() {
			verdict = ""
			rows.Scan(&submissionID, &snapshotID, &code, &submittedAt, &verdict)

			submissions = append(submissions, &SubmissionInfo{
				ID:          submissionID,
				SnapshotID:  snapshotID,
				Code:        code,
				Grade:       verdict,
				SubmittedAt: submittedAt,
			})
		}
		sort.SliceStable(submissions, func(i, j int) bool {
			if submissions[i].Grade == "" && submissions[j].Grade == "" {
				return submissions[i].SubmittedAt.Before(submissions[j].SubmittedAt)
			}
			if submissions[i].Grade == "" {
				return true
			}
			if submissions[j].Grade == "" {
				return false
			}
			return submissions[i].SubmittedAt.Before(submissions[j].SubmittedAt)
		})
	} else {
		http.Error(w, "You are not authorized to access!", http.StatusUnauthorized)
	}
	// TODO(shiplu): sort the messages array
	// sort.Slice(helpRequests, func(i, j int) bool { return helpRequests[i].GivenAt.Before(helpRequests[j].GivenAt) })
	data := &SubmissionDashboard{
		StudentName: getStudentName(studentID),
		ProblemName: getProblemNameFromID(problemID),
		Submissions: submissions,
		StudentID:   studentID,
		ProblemID:   problemID,
		UserID:      uid,
		UserRole:    role,
		Password:    r.FormValue("password"),
	}
	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
}

func hasMessageBackFeedbackHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	feedbackID, _ := strconv.Atoi(r.FormValue("feedback_id"))
	userRole := r.FormValue("role")
	fmt.Print(feedbackID, userRole, uid)
	row, err := Database.Query("select * from message_back_feedback where message_feedback_id=? and author_id = ? and author_role = ?", feedbackID, uid, userRole)
	defer row.Close()
	if err != nil {
		log.Fatal(err)
	}
	if row.Next() {
		fmt.Fprint(w, "yes")
	} else {
		fmt.Fprint(w, "no")
	}
}

func studentDashboardCodeSpaceHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	role := r.FormValue("role")
	problemID, _ := strconv.Atoi(r.FormValue("problem_id"))
	studentID, _ := strconv.Atoi(r.FormValue("student_id"))
	temp := template.New("")
	ownFuncs := template.FuncMap{"getEditorMode": getEditorMode}
	t, err := temp.Funcs(ownFuncs).Parse(CODE_SNAPSHOT_TAB_TEMPLATE)
	if err != nil {
		log.Fatal(err)
	}
	students := getAllStudents()

	// Get latest snapshot from DB
	latestSnapshot := &Snapshot{}
	if _, ok := StudentSnapshot[studentID][problemID]; ok {
		latestSnapshot = Snapshots[StudentSnapshot[studentID][problemID]]
	} else {
		latestSnapshot = getLatestSnapshot(studentID, problemID)
	}

	// Get all student messages from DB
	var messages = make([]*MessageDashBoard, 0)
	_, ok := HelpEligibleStudents[problemID][uid]
	if role == "teacher" || uid == studentID || (PeerTutorAllowed && ok) {
		rows, err := Database.Query("select M.id, M.snapshot_id, M.message, M.author_id, M.author_role, M.type, C.Code, C.event, M.given_at from message M, code_snapshot C where M.snapshot_id = C.id and C.problem_id = ? and C.student_id = ?", problemID, studentID)
		defer rows.Close()
		if err != nil {
			log.Fatal(err)
		}
		var snapshotID, authorID, messageType, messageID int
		var message, authorRole, code, event string
		var givenAt time.Time
		for rows.Next() {
			rows.Scan(&messageID, &snapshotID, &message, &authorID, &authorRole, &messageType, &code, &event, &givenAt)
			name := ""
			if authorRole == "teacher" {
				name = getTeacherName(authorID)
			} else {
				name = students[authorID]
			}
			messages = append(messages, &MessageDashBoard{
				ID:         messageID,
				Name:       name,
				Role:       authorRole,
				Message:    message,
				Type:       messageType,
				Event:      event,
				GivenAt:    givenAt,
				SnapshotID: snapshotID,
				Code:       code,
				Feedbacks:  getMessageFeedbacks(messageID, uid, role),
			})
		}
	} else {
		http.Error(w, "You are not authorized to access!", http.StatusUnauthorized)
	}

	// Sort the messages by descding order of GivenAt
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].GivenAt.After(messages[j].GivenAt)
	})

	feedback := &FeedbackProvisionDashBoard{
		StudentName:  students[studentID],
		ProblemName:  latestSnapshot.ProblemName,
		LastSnapshot: latestSnapshot,
		Messages:     messages,
		StudentID:    studentID,
		ProblemID:    problemID,
		UserID:       uid,
		UserRole:     role,
		Password:     r.FormValue("password"),
	}

	// Get all submissions from DB.
	var submissions = make([]*SubmissionInfo, 0)
	_, ok = HelpEligibleStudents[problemID][uid]
	if role == "teacher" || uid == studentID || (PeerTutorAllowed && ok) {
		rows, err := Database.Query("select id, snapshot_id, student_code, code_submitted_at, verdict from submission where student_id = ? and problem_id = ?", studentID, problemID)
		defer rows.Close()
		if err != nil {
			log.Fatal(err)
		}
		var snapshotID, submissionID int
		var verdict, code string
		var submittedAt time.Time
		for rows.Next() {
			verdict = ""
			rows.Scan(&submissionID, &snapshotID, &code, &submittedAt, &verdict)

			submissions = append(submissions, &SubmissionInfo{
				ID:          submissionID,
				SnapshotID:  snapshotID,
				Code:        code,
				Grade:       verdict,
				SubmittedAt: submittedAt,
			})
		}
		sort.SliceStable(submissions, func(i, j int) bool {
			if submissions[i].Grade == "" && submissions[j].Grade == "" {
				return submissions[i].SubmittedAt.After(submissions[j].SubmittedAt)
			}
			if submissions[i].Grade == "" {
				return true
			}
			if submissions[j].Grade == "" {
				return false
			}
			return submissions[i].SubmittedAt.After(submissions[j].SubmittedAt)
		})
	} else {
		http.Error(w, "You are not authorized to access!", http.StatusUnauthorized)
	}

	submission := &SubmissionDashboard{
		StudentName: getStudentName(studentID),
		ProblemName: getProblemNameFromID(problemID),
		Submissions: submissions,
		StudentID:   studentID,
		ProblemID:   problemID,
		UserID:      uid,
		UserRole:    role,
		Password:    r.FormValue("password"),
		Username:    getName(uid, role),
	}

	// Get student status
	rows, err := Database.Query("select coding_stat, help_stat, submission_stat, tutoring_stat from student_status where problem_id=? and student_id=?", problemID, studentID)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	var codingStat, submissionStat, helpStat, tutoringStat string
	// var studentInfo []*DashBoardStudentInfo
	studentStats := DashBoardStudentInfo{}
	for rows.Next() {
		rows.Scan(&codingStat, &helpStat, &submissionStat, &tutoringStat)
		if role == "teacher" || uid == studentID || (PeerTutorAllowed && ok) {
			studentStats.CodingStat = codingStat
			studentStats.HelpStat = helpStat
			studentStats.SubmissionStat = submissionStat
			studentStats.TutoringStat = tutoringStat
		}
	}
	rows.Close()

	data := TemplateDate{
		Submission:     *submission,
		Feedback:       *feedback,
		Status:         studentStats,
		ChatgptaServer: Config.ChatgptaServer,
		UserID:         uid,
		UserRole:       role,
		Password:       r.FormValue("password"),
		Username:       getName(uid, role),
		CourseName:		Config.CourseName,
	}

	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
}
