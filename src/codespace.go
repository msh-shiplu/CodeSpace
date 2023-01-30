package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type CodeSpaceData struct {
	Snapshots []*Snapshot
	UserID    int
	UserRole  string
	Password  string
}

type FeedbackData struct {
	FeedbackID      int
	Feedback        string
	FeedbackTime    time.Time
	Upvote          int
	Downvote        int
	CurrentUserVote string
	GivenBy         string
	Code            string
}
type SnapshotData struct {
	Snapshot       *Snapshot
	HelpRequestIDs []int
	UserID         int
	UserRole       string
	Feedbacks      []*FeedbackData
	Password       string
}

type HelpRequest struct {
	ID          int
	StudentName string
	Explanation string
	NumReply    int
	GivenAt     time.Time
	SnapshotID  int
	Snapshot    string
	ProblemName string
	UserID      int
	UserRole    string
	Password    string
}

type HelpRequestListData struct {
	HelpRequests  []*HelpRequest
	NumHelpNeeded int
	UserID        int
	UserRole      string
	Password      string
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
		str = strconv.Itoa(m) + " min "
	}
	str += strconv.Itoa(s) + " sec"
	return str
}

func formatTimeSince(t time.Time) string {
	d := time.Now().Sub(t)
	return formatTimeDuration(d)
}

func getVoteCount(feedbackID int, voteType string) int {
	vote := 0
	rows, err := Database.Query("select count(*) from snapshot_back_feedback where is_helpful = ? and snapshot_feedback_id = ?", voteType, feedbackID)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		rows.Scan(&vote)
	}
	return vote
}

func add(x int, y int) int {
	return x + y
}

func codespaceHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
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
			if s.StudentID == uid {
				snapshots = append(snapshots, s)
			} else if _, ok := HelpEligibleStudents[s.ProblemID][uid]; ok {
				snapshots = append(snapshots, s)
			}
		}
	} else {
		for _, s := range Snapshots {
			snapshots = append(snapshots, s)
		}
	}
	sort.Slice(snapshots, func(i, j int) bool { return snapshots[i].LastUpdated.After(snapshots[j].LastUpdated) })
	data := &CodeSpaceData{
		Snapshots: snapshots,
		UserID:    uid,
		UserRole:  role,
		Password:  r.FormValue("password"),
	}
	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
}

func getCodeSnapshotHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	snapshotID := 0
	studentID := 0
	problemID := 0
	if r.FormValue("snapshot_id") != "" {
		snapshotID, _ = strconv.Atoi(r.FormValue("snapshot_id"))
		rows, err := Database.Query("select student_id, problem_id from code_snapshot where id = ?", snapshotID)
		defer rows.Close()
		if err != nil {
			log.Fatal(err)
		}
		if rows.Next() {
			rows.Scan(&studentID, &problemID)
		}
		rows.Close()
	} else {
		studentID, _ = strconv.Atoi(r.FormValue("student_id"))
		problemID, _ = strconv.Atoi(r.FormValue("problem_id"))
	}
	role := r.FormValue("role")
	temp := template.New("")
	ownFuncs := template.FuncMap{"getEditorMode": getEditorMode, "add": add}
	t, err := temp.Funcs(ownFuncs).Parse(CODE_SNAPSHOT_TEMPLATE)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := Database.Query("select id from code_explanation where student_id = ? and problem_id = ?", studentID, problemID)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	var helpID int
	var helpIDs []int
	for rows.Next() {
		rows.Scan(&helpID)
		helpIDs = append(helpIDs, helpID)
	}
	rows.Close()

	rows, err = Database.Query("select F.id, F.feedback, F.author_id, F.author_role, F.given_at, C.code from code_snapshot C, snapshot_feedback F where C.id=F.snapshot_id and C.student_id=? and C.problem_id=? order by F.given_at desc", studentID, problemID)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	feedbackID, feedback, authorID, authorRole, givenAt, code := 0, "", 0, "", time.Now(), ""
	var feedbacks []*FeedbackData

	upvote, downvote := 0, 0
	currentUserVote := ""
	for rows.Next() {
		rows.Scan(&feedbackID, &feedback, &authorID, &authorRole, &givenAt, &code)

		upvote = getVoteCount(feedbackID, "yes")
		downvote = getVoteCount(feedbackID, "no")
		rows2, err := Database.Query("select is_helpful from snapshot_back_feedback where snapshot_feedback_id=? and author_id=? and author_role=?", feedbackID, uid, role)
		defer rows2.Close()
		if err != nil {
			log.Fatal(err)
		}
		for rows2.Next() {
			rows2.Scan(&currentUserVote)
		}
		rows2.Close()
		if authorRole == "teacher" {
			rows2, err = Database.Query("select name from teacher where id=?", authorID)
		} else {
			rows2, err = Database.Query("select name from student where id=?", authorID)
		}
		defer rows2.Close()
		if err != nil {
			log.Fatal(err)
		}
		authorName := ""
		if rows2.Next() {
			rows2.Scan(&authorName)
		}
		rows2.Close()
		feedbacks = append(feedbacks, &FeedbackData{
			FeedbackID:      feedbackID,
			Feedback:        feedback,
			FeedbackTime:    givenAt,
			Upvote:          upvote,
			Downvote:        downvote,
			CurrentUserVote: currentUserVote,
			GivenBy:         authorName,
			Code:            code,
		})
		currentUserVote = ""

	}
	idx := StudentSnapshot[studentID][problemID]
	data := &SnapshotData{
		Snapshot:       Snapshots[idx],
		HelpRequestIDs: helpIDs,
		UserID:         uid,
		UserRole:       role,
		Feedbacks:      feedbacks,
		Password:       r.FormValue("password"),
	}
	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
	otherInfo := struct {
		Referral  string
		ProblemID int
	}{
		Referral:  r.URL.String(),
		ProblemID: problemID,
	}
	elog, _ := json.Marshal((otherInfo))
	logEvent("willing to help", uid, role, "click", string(elog))
}

func getNumberOfReply(snapshotID int) int {
	rows, err := Database.Query("select count(*) from snapshot_feedback where snapshot_id = ?", snapshotID)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	c := 0
	if rows.Next() {
		rows.Scan(&c)
	}
	return c
}

func helpRequestListHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	role := r.FormValue("role")
	temp := template.New("")
	ownFuncs := template.FuncMap{"formatTimeSince": formatTimeSince}
	t, err := temp.Funcs(ownFuncs).Parse(HELP_REQUEST_LIST_TEMPLATE)
	if err != nil {
		log.Fatal(err)
	}
	var helpRequests []*HelpRequest
	var helpNeededCount int
	var numReply int
	if role == "student" {
		for _, s := range HelpSubmissions {
			if _, ok := HelpEligibleStudents[s.Pid][uid]; ok || s.Uid == uid {
				numReply = getNumberOfReply(s.SnapshotID)
				helpRequests = append(helpRequests, &HelpRequest{
					ID:          s.Sid,
					NumReply:    numReply,
					StudentName: Students[s.Uid].Name,
					GivenAt:     s.At,
				})
				if numReply == 0 {
					helpNeededCount++
				}
			}
		}
	} else {
		for _, s := range HelpSubmissions {
			numReply = getNumberOfReply(s.SnapshotID)
			helpRequests = append(helpRequests, &HelpRequest{
				ID:          s.Sid,
				NumReply:    numReply,
				StudentName: Students[s.Uid].Name,
				GivenAt:     s.At,
			})
			if numReply == 0 {
				helpNeededCount++
			}
		}
	}
	sort.Slice(helpRequests, func(i, j int) bool { return helpRequests[i].GivenAt.Before(helpRequests[j].GivenAt) })
	data := &HelpRequestListData{
		HelpRequests:  helpRequests,
		NumHelpNeeded: helpNeededCount,
		UserID:        uid,
		UserRole:      role,
		Password:      r.FormValue("password"),
	}
	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
}

func viewHelpRequestHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	requestID, _ := strconv.Atoi(r.FormValue("request_id"))
	role := r.FormValue("role")
	pw := r.FormValue("password")
	temp := template.New("")
	ownFuncs := template.FuncMap{"getEditorMode": getEditorMode, "formatTimeSince": formatTimeSince}
	t, err := temp.Funcs(ownFuncs).Parse(HELP_REQUEST_VIEW_TEMPLATE)
	if err != nil {
		log.Fatal(err)
	}
	data := &HelpRequest{}
	problemID := -1
	for _, s := range HelpSubmissions {
		if s.Sid == requestID {
			data = &HelpRequest{
				StudentName: Students[s.Uid].Name,
				Explanation: s.Content,
				GivenAt:     s.At,
				SnapshotID:  s.SnapshotID,
				Snapshot:    s.Snapshot,
				ProblemName: s.Filename,
				UserID:      uid,
				UserRole:    role,
				Password:    pw,
			}
			problemID = s.Pid
		}
	}
	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
	otherInfo := struct {
		Referral  string
		ProblemID int
	}{
		Referral:  r.URL.String(),
		ProblemID: problemID,
	}
	elog, _ := json.Marshal(otherInfo)
	logEvent("willing to help", uid, role, "click", string(elog))

}

func setPeerTutorHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	turnOn, _ := strconv.Atoi(r.FormValue("turn_on"))
	if turnOn == 1 {
		PeerTutorAllowed = true
	} else if turnOn == 0 {
		PeerTutorAllowed = false
	}
}

func peerTutorHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	filename := r.FormValue("filename")
	if prob, ok := ActiveProblems[filename]; ok {
		if eligible, ok := HelpEligibleStudents[prob.Info.Pid][uid]; ok {
			if PeerTutorAllowed && eligible {
				fmt.Fprint(w, "redirect")
				return
			}
		}
	}
	fmt.Fprintf(w, "You are not eligible for peer tutoring for this problem!")
}
