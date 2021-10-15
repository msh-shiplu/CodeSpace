package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

func getFeedbackData(rows *sql.Rows, uid int, role string) []*FeedbackData {
	var feedbacks []*FeedbackData
	feedbackID, feedback, authorID, authorRole, givenAt, code := 0, "", 0, "", time.Now(), ""

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
	return feedbacks
}

func studentViewsFeedbackHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	spid := r.FormValue("pid")
	pid := -1
	nextPid := -1
	filename := ""
	if spid == "" {
		rows, err := Database.Query("select id, filename from problem order by id desc limit 2")
		defer rows.Close()
		if err != nil {
			log.Fatal(err)
		}
		if rows.Next() {
			rows.Scan(&pid, &filename)
		}
		tmp := ""
		if rows.Next() {
			rows.Scan(&nextPid, &tmp)
		}
		rows.Close()
	} else {
		pid, _ = strconv.Atoi(spid)
		rows, err := Database.Query("select id from problem where id<? order by id desc limit 1", pid)
		defer rows.Close()
		if err != nil {
			log.Fatal(err)
		}
		if rows.Next() {
			rows.Scan(&nextPid)
		}
		rows.Close()

		rows2, err2 := Database.Query("select filename from problem where id = ?", pid)
		defer rows2.Close()
		if err2 != nil {
			log.Fatal(err2)
		}
		if rows2.Next() {
			rows2.Scan(&filename)
		}
		rows2.Close()
	}
	viewType := r.FormValue("viewtype")
	role := r.FormValue("role")
	var rows *sql.Rows
	var err error
	if viewType == "forme" {
		rows, err = Database.Query("select F.id, F.feedback, F.author_id, F.author_role, F.given_at, C.code from code_snapshot C, snapshot_feedback F where C.id=F.snapshot_id and C.student_id=? and C.problem_id = ? order by F.given_at desc", uid, pid)
		defer rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	} else if viewType == "all" {
		active := false
		for _, prob := range ActiveProblems {
			if prob.Info.Pid == pid {
				if prob.Active {
					active = true
				}
				break
			}
		}
		if active == true {
			if _, ok := HelpEligibleStudents[pid][uid]; ok {
				rows, err = Database.Query("select F.id, F.feedback, F.author_id, F.author_role, F.given_at, C.code from code_snapshot C, snapshot_feedback F where C.id=F.snapshot_id and C.problem_id = ? order by F.given_at desc", pid)
			} else {
				rows, err = Database.Query("select F.id, F.feedback, F.author_id, F.author_role, F.given_at, C.code from code_snapshot C, snapshot_feedback F where C.id=F.snapshot_id and C.student_id=? and C.problem_id = ? order by F.given_at desc", uid, pid)
			}
		} else {
			rows, err = Database.Query("select F.id, F.feedback, F.author_id, F.author_role, F.given_at, C.code from code_snapshot C, snapshot_feedback F where C.id=F.snapshot_id and C.problem_id = ? order by F.given_at desc", pid)
		}
		defer rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("Invalid parameter!")
	}
	feedbacks := getFeedbackData(rows, uid, role)
	data := struct {
		Feedbacks  []*FeedbackData
		ViewType   string
		UserRole   string
		UserID     int
		Password   string
		Filename   string
		CurrentPid int
		NextPid    int
	}{
		Feedbacks:  feedbacks,
		ViewType:   viewType,
		UserRole:   role,
		UserID:     uid,
		Password:   r.FormValue("password"),
		Filename:   filename,
		CurrentPid: pid,
		NextPid:    nextPid,
	}
	temp := template.New("")
	ownFuncs := template.FuncMap{"getEditorMode": getEditorMode}
	t, err := temp.Funcs(ownFuncs).Parse(STUDENT_VIEWS_FEEDBACK_TEMPLATE)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
}

func teacherViewsFeedbackHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	spid := r.FormValue("pid")
	pid := -1
	nextPid := -1
	filename := ""
	if spid == "" {
		rows, err := Database.Query("select id, filename from problem order by id desc limit 2")
		defer rows.Close()
		if err != nil {
			log.Fatal(err)
		}
		if rows.Next() {
			rows.Scan(&pid, &filename)
		}
		tmp := ""
		if rows.Next() {
			rows.Scan(&nextPid, &tmp)
		}
		rows.Close()
	} else {
		pid, _ = strconv.Atoi(spid)
		rows, err := Database.Query("select id from problem where id<? order by id desc limit 1", pid)
		defer rows.Close()
		if err != nil {
			log.Fatal(err)
		}
		if rows.Next() {
			rows.Scan(&nextPid)
		}
		rows.Close()

		rows2, err2 := Database.Query("select filename from problem where id = ?", pid)
		defer rows2.Close()
		if err2 != nil {
			log.Fatal(err2)
		}
		if rows2.Next() {
			rows2.Scan(&filename)
		}
		rows2.Close()
	}
	role := r.FormValue("role")
	rows, err := Database.Query("select F.id, F.feedback, F.author_id, F.author_role, F.given_at, C.code from code_snapshot C, snapshot_feedback F where C.id=F.snapshot_id and C.problem_id = ? order by F.given_at desc", pid)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	feedbacks := getFeedbackData(rows, uid, role)
	data := struct {
		Feedbacks  []*FeedbackData
		UserRole   string
		UserID     int
		Password   string
		Filename   string
		CurrentPid int
		NextPid    int
	}{
		Feedbacks:  feedbacks,
		UserRole:   role,
		UserID:     uid,
		Password:   r.FormValue("password"),
		Filename:   filename,
		CurrentPid: pid,
		NextPid:    nextPid,
	}
	temp := template.New("")
	ownFuncs := template.FuncMap{"getEditorMode": getEditorMode}
	t, err := temp.Funcs(ownFuncs).Parse(TEACHER_VIEWS_FEEDBACK_TEMPLATE)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
}
