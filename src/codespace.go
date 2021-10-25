package main

import (
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
	Snapshot  *Snapshot
	UserID    int
	UserRole  string
	Feedbacks []*FeedbackData
	Password  string
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
	fmt.Printf("before\n%+v\n%+v\n", snapshots[0], snapshots[1])
	sort.Slice(snapshots, func(i, j int) bool { return snapshots[i].LastUpdated.After(snapshots[j].LastUpdated) })
	fmt.Printf("after\n%+v\n%+v\n", snapshots[0], snapshots[1])
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
	studentID, _ := strconv.Atoi(r.FormValue("student_id"))
	problemID, _ := strconv.Atoi(r.FormValue("problem_id"))
	role := r.FormValue("role")
	temp := template.New("")
	ownFuncs := template.FuncMap{"getEditorMode": getEditorMode}
	t, err := temp.Funcs(ownFuncs).Parse(CODE_SNAPSHOT_TEMPLATE)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := Database.Query("select F.id, F.feedback, F.author_id, F.author_role, F.given_at, C.code from code_snapshot C, snapshot_feedback F where C.id=F.snapshot_id and C.student_id=? and C.problem_id=? order by F.given_at desc", studentID, problemID)
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
		Snapshot:  Snapshots[idx],
		UserID:    uid,
		UserRole:  role,
		Feedbacks: feedbacks,
		Password:  r.FormValue("password"),
	}
	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
}
