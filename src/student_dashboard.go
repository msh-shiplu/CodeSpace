package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

type FeedbackDashBaord struct {
	Name     string
	Role     string
	Feedback string
	GivenAt  time.Time
}

type MessageDashBoard struct {
	Name       string
	Role       string
	Message    string
	Type       int // 0 = help request, 1 = unsolicited
	GivenAt    time.Time
	Code       string
	SnapshotID int
	Feedbacks  []*FeedbackDashBaord
}

type FeedbackProvisionDashBoard struct {
	StudentName    string
	ProblemName    string
	Code           string
	LastSnapshotAt time.Time
	Messages       []*MessageDashBoard
	UserID         int
	UserRole       string
	Password       string
}

func getMessageFeedbacks(messageID int) []*FeedbackDashBaord {
	rows, err := Database.Query("select feedback, author_id, author_role, given_at from message_feedback, where message_id = ?", messageID)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	var feedback, authorRole string
	var authorID int
	var givenAt time.Time

	feedbacks := make([]*FeedbackDashBaord, 0)

	for rows.Next() {
		rows.Scan(&feedback, &authorID, &authorRole, &givenAt)
		name := ""
		if authorRole == "teacher" {
			name = Teacher[authorID]
		} else {
			name = Students[authorID].Name
		}
		feedbacks = append(feedbacks, &FeedbackDashBaord{
			Name:     name,
			Role:     authorRole,
			Feedback: feedback,
			GivenAt:  givenAt,
		})
	}
	rows.Close()
	return feedbacks
}

func studentDashboardFeedbackProvisionHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	role := r.FormValue("role")
	problemID, _ := strconv.Atoi(r.FormValue("problem_id"))
	studentID, _ := strconv.Atoi(r.FormValue("student_id"))
	temp := template.New("")
	ownFuncs := template.FuncMap{"formatTimeSince": formatTimeSince}
	t, err := temp.Funcs(ownFuncs).Parse(FEEDBACK_PROVISION_TEMPLATE)
	if err != nil {
		log.Fatal(err)
	}
	var messages = make([]*MessageDashBoard, 0)
	_, ok := HelpEligibleStudents[problemID][uid]
	if role == "teacher" || uid == studentID || (PeerTutorAllowed && ok) {
		rows, err := Database.Query("select M.id, M.snapshot_id, M.message, M.author_id, M.author_role, M.given_at, M.type, C.Code from message M, code_snapshot C where M.snapshot_id = C.id and C.problem_id = ? and C.student_id = ?", problemID, studentID)
		defer rows.Close()
		if err != nil {
			log.Fatal(err)
		}
		var snapshotID, authorID, messageType, messageID int
		var message, authorRole, code string
		var givenAt time.Time
		for rows.Next() {
			rows.Scan(&messageID, &snapshotID, &message, &authorID, &authorRole, &givenAt, &messageType, &code)
			name := ""
			if authorRole == "teacher" {
				name = Teacher[authorID]
			} else {
				name = Students[authorID].Name
			}
			messages = append(messages, &MessageDashBoard{
				Name:       name,
				Role:       authorRole,
				Message:    message,
				Type:       messageType,
				GivenAt:    givenAt,
				SnapshotID: snapshotID,
				Code:       code,
				Feedbacks:  getMessageFeedbacks(messageID),
			})
		}
	} else {
		http.Error(w, "You are not authorized to access!", http.StatusUnauthorized)
	}
	// TODO(shiplu): sort the messages array
	// sort.Slice(helpRequests, func(i, j int) bool { return helpRequests[i].GivenAt.Before(helpRequests[j].GivenAt) })
	latestSnapshot := Snapshots[StudentSnapshot[studentID][problemID]]
	data := &FeedbackProvisionDashBoard{
		StudentName:    Students[studentID].Name,
		ProblemName:    latestSnapshot.ProblemName,
		Code:           latestSnapshot.Code,
		LastSnapshotAt: latestSnapshot.LastUpdated,
		Messages:       messages,
		UserID:         uid,
		UserRole:       role,
		Password:       r.FormValue("password"),
	}
	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
}