//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func create_tables() {
	execSQL := func(s string) {
		sql_stmt, err := Database.Prepare(s)
		if err != nil {
			log.Fatal(err)
		}
		sql_stmt.Exec()
	}
	execSQL("create table if not exists student (id integer primary key, name text unique, password text)")
	execSQL("create table if not exists teacher (id integer primary key, name text unique, password text)")
	execSQL("create table if not exists attendance (id integer primary key, student_id integer, attendance_at timestamp)")
	execSQL("create table if not exists tag (id integer primary key, topic_description text unique)")
	execSQL("create table if not exists problem (id integer primary key, teacher_id integer, problem_description blob, answer text, filename text, merit integer, effort integer, attempts integer, topic_id integer, tag integer, problem_uploaded_at timestamp, problem_ended_at timestamp)")
	execSQL("create table if not exists submission (id integer primary key, problem_id integer, student_id integer, student_code blob, submission_category integer, code_submitted_at timestamp, completed timestamp, verdict text, attempt_number integer)")
	execSQL("create table if not exists score (id integer primary key, problem_id integer, student_id integer, teacher_id integer, score integer, graded_submission_number integer, score_given_at timestamp, unique(problem_id,student_id))")
	execSQL("create table if not exists feedback (id integer primary key, teacher_id integer, student_id integer, feedback text, feedback_given_at timestamp, submission_id integer)")
	execSQL("create table if not exists test_case (id integer primary key, problem_id integer, student_id integer, test_cases text, added_at timestamp)")
	execSQL("create table if not exists code_explanation (id integer primary key, problem_id integer, student_id integer, snapshot_id integer, trying_what text, need_help_with text, code_submitted_at timestamp)")
	execSQL("create table if not exists help_message (id integer primary key, code_explanation_id integer, student_id integer, message text, given_at timestamp, useful text, updated_at timestamp)")
	execSQL("create table if not exists code_snapshot (id integer primary key, student_id integer, problem_id integer, code blob, last_updated_at timestamp, status int default 0)") // 0 = not submitted, 1 = submitted but not graded, 2 = submitted and incorrect, 3 = submitted and correct
	execSQL("create table if not exists snapshot_feedback (id integer primary key, snapshot_id integer, feedback text, author_id integer, author_role string, given_at timestamp)")
	execSQL("create table if not exists snapshot_back_feedback (id integer primary key, snapshot_feedback_id integer, author_id integer, author_role string, is_helpful string, given_at timestamp)")
	execSQL("create table if not exists message (id integer primary key, snapshot_id integer, message text, author_id integer, author_role string, given_at timestamp, type integer)")
	execSQL("create table if not exists message_feedback (id integer primary key, message_id integer, feedback text, author_id integer, author_role string, given_at timestamp)")
	execSQL("create table if not exists help_eligible (id integer primary key, problem_id integer, student_id integer, became_eligible_at timestamp)")
	execSQL("create table if not exists user_event_log (id integer primary key, name string, user_id integer, user_type string, event_type string, referral_info string, event_time timestamp)")
	execSQL("create table if not exists student_status (id integer primary key, student_id integer, problem_id integer, coding_stat string, help_stat string, submission_stat string, tutoring_stat string, last_updated_at timestamp)")
	execSQL("create table if not exists problem_statistics (id integer primary key, problem_id integer not null, active integer default 0, submission integer default 0, help_request integer default 0, graded_correct integer default 0, graded_incorrect integer default 0)")
	// foreign key example: http://www.sqlitetutorial.net/sqlite-foreign-key/
}

//-----------------------------------------------------------------
func init_database(db_name string) {
	var err error
	prepare := func(s string) *sql.Stmt {
		stmt, err := Database.Prepare(s)
		if err != nil {
			log.Fatal(err)
		}
		return stmt
	}

	Database, err = sql.Open("sqlite3", db_name)
	if err != nil {
		log.Fatal(err)
	}
	create_tables()
	AddStudentSQL = prepare("insert into student (name, password) values (?, ?)")
	AddTeacherSQL = prepare("insert into teacher (name, password) values (?, ?)")
	AddProblemSQL = prepare("insert into problem (teacher_id, problem_description, answer, filename, merit, effort, attempts, topic_id, tag, problem_uploaded_at) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	AddSubmissionSQL = prepare("insert into submission (problem_id, student_id, student_code, submission_category, attempt_number, code_submitted_at) values (?, ?, ?, ?, ?, ?)")
	AddSubmissionCompleteSQL = prepare("insert into submission (problem_id, student_id, student_code, submission_category, attempt_number, code_submitted_at, completed) values (?, ?, ?, ?, ?, ?, ?)")
	CompleteSubmissionSQL = prepare("update submission set completed=?, verdict=? where id=?")
	AddScoreSQL = prepare("insert into score (problem_id, student_id, teacher_id, score, graded_submission_number, score_given_at) values (?, ?, ?, ?, ?, ?)")
	AddFeedbackSQL = prepare("insert into feedback (teacher_id, student_id, feedback, feedback_given_at, submission_id) values (?, ?, ?, ?, ?)")
	UpdateScoreSQL = prepare("update score set teacher_id=?, score=?, graded_submission_number=? where id=?")
	AddAttendanceSQL = prepare("insert into attendance (student_id, attendance_at) values (?, ?)")
	AddTagSQL = prepare("insert into tag (topic_description) values (?)")
	AddTestCaseSQL = prepare("insert into test_case (problem_id, student_id, test_cases, added_at) values (?, ?, ?, ?)")
	UpdateTestCaseSQL = prepare("update test_case set test_cases=?, added_at=? where id=?")
	AddHelpSubmissionSQL = prepare("insert into code_explanation (problem_id, student_id, snapshot_id, trying_what, need_help_with, code_submitted_at) values(?, ?, ?, ?, ?, ?)")
	AddHelpMessageSQL = prepare("insert into help_message (code_explanation_id, student_id, message, given_at) values (?, ?, ?, ?)")
	UpdateHelpMessageSQL = prepare("update help_message set useful=?, updated_at=? where id=?")
	AddCodeSnapshotSQL = prepare("insert into code_snapshot (student_id, problem_id, code, status, last_updated_at) values(?, ?, ?, ?, ?)")
	AddSnapShotFeedbackSQL = prepare("insert into snapshot_feedback (snapshot_id, feedback, author_id, author_role, given_at) values(?, ?, ?, ?, ?)")
	AddSnapshotBackFeedbackSQL = prepare("insert into snapshot_back_feedback (snapshot_feedback_id, author_id, author_role, is_helpful, given_at) values(?, ?, ?, ?, ?)")
	UpdateSnapshotBackFeedbackSQL = prepare("update snapshot_back_feedback set is_helpful=?, given_at=? where snapshot_feedback_id=? and author_id=? and author_role=?")
	UpdateProblemEndTimeSQL = prepare("update problem set problem_ended_at=? where id=?")
	AddHelpEligibleSQL = prepare("insert into help_eligible (problem_id, student_id, became_eligible_at) values(?, ?, ?)")
	AddUserEventLogSQL = prepare("insert into user_event_log (name, user_id, user_type, event_type, referral_info, event_time) values(?, ?, ?, ?, ?, ?)")
	AddStudentStatusSQL = prepare("insert into student_status (student_id, problem_id, coding_stat, help_stat, submission_stat, tutoring_stat, last_updated_at) values(?, ?, ?, ?, ?, ?, ?)")
	UpdateStudentCodingStatSQL = prepare("update student_status set coding_stat = ?, last_updated_at = ? where student_id = ? and problem_id = ?")
	UpdateStudentHelpStatSQL = prepare("update student_status set help_stat = ?, last_updated_at = ? where student_id = ? and problem_id = ?")
	UpdateStudentSubmissionStatSQL = prepare("update student_status set submission_stat = ?, last_updated_at = ? where student_id = ? and problem_id = ?")
	UpdateStudentTutoringStatSQL = prepare("update student_status set tutoring_stat = ?, last_updated_at = ? where student_id = ? and problem_id = ?")
	AddMessageSQL = prepare("insert into message (snapshot_id, message, author_id, author_role, given_at, type) values (?, ?, ?, ?, ?, ?)")
	AddMessageFeedbackSQL = prepare("insert into message_feedback (message_id, feedback, author_id, author_role, given_at) values(?, ?, ?, ?, ?)")
	AddProblemStatisticsSQL = prepare("insert into problem_statistics (problem_id, active, submission, help_request, graded_correct, graded_incorrect) values (?, 0, 0, 0, 0, 0)")
	IncProblemStatActiveSQL = prepare("update problem_statistics set active = active + 1 where problem_id = ?")
	IncProblemStatSubmissionSQL = prepare("update problem_statistics set submission = submission + 1 where problem_id = ?")
	IncProblemStatHelpSQL = prepare("update problem_statistics set help_request = help_request + 1 where problem_id = ?")
	IncProblemStatGradedCorrectSQL = prepare("update problem_statistics set graded_correct = graded_correct + 1 where problem_id = ?")
	IncProblemStatGradedIncorrectSQL = prepare("update problem_statistics set graded_incorrect = graded_incorrect + 1 where problem_id = ?")
	// Initialize passcode for current session and default board
	Passcode = RandStringRunes(12)
	Students[0] = &StudenInfo{
		Boards: make([]*Board, 0),
	}
}

//-----------------------------------------------------------------
// Add or update score based on a decision. If decision is "correct"
// a new problem, if there's one, is added to student's board.
//-----------------------------------------------------------------
func add_or_update_score(decision string, pid, student_id, teacher_id, partial_credits int) string {
	mesg := ""

	// Find score information for this student (student_id) for this problem (pid)

	score_id, current_points, current_attempts, current_tid := 0, 0, 0, 0
	rows, _ := Database.Query("select id, score, graded_submission_number, teacher_id from score where problem_id=? and student_id=?", pid, student_id)
	for rows.Next() {
		rows.Scan(&score_id, &current_points, &current_attempts, &current_tid)
		break
	}
	rows.Close()

	// Find merit points and effort points for this problem (pid)
	merit, effort := 0, 0
	rows, _ = Database.Query("select merit, effort from problem where id=?", pid)
	for rows.Next() {
		rows.Scan(&merit, &effort)
		break
	}
	rows.Close()

	// Determine points for this student
	points, teacher := 0, teacher_id
	if decision == "correct" {
		points = merit
		mesg = "Answer is correct."
	} else {
		if partial_credits < merit {
			points = partial_credits
		} else {
			points = effort
		}

		// If the problem was previously graded correct, this submission
		// does not reduce it.  Grading is asynchronous.
		if points < current_points {
			points = current_points
			teacher = current_tid
		}
		mesg = "Answer is incorrect."
	}
	// m := add_next_problem_to_board(pid, student_id, decision)
	// mesg = mesg + m

	// Add a new score or update a current score for this student & problem
	if score_id == 0 {
		_, err := AddScoreSQL.Exec(pid, student_id, teacher_id, points, current_attempts+1, time.Now())
		if err != nil {
			mesg = fmt.Sprintf("Unable to add score: %d %d %d", pid, student_id, teacher_id)
			writeLog(Config.LogFile, mesg)
			return mesg
		}
	} else {
		_, err := UpdateScoreSQL.Exec(teacher, points, current_attempts+1, score_id)
		if err != nil {
			mesg = fmt.Sprintf("Unable to update score: %d %d", teacher, score_id)
			writeLog(Config.LogFile, mesg)
			return mesg
		}
	}
	return mesg
}

func addOrUpdateStudentStatus(studentID int, problemID int, codingStat string, helpStat string, submissionStat string, tutoringStat string) {
	rows, err := Database.Query("select * from student_status where student_id = ? and problem_id = ?", studentID, problemID)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	if rows.Next() {
		now := time.Now()
		if codingStat != "" {
			UpdateStudentCodingStatSQL.Exec(codingStat, now, studentID, problemID)
		}
		if helpStat != "" {
			UpdateStudentHelpStatSQL.Exec(helpStat, now, studentID, problemID)
		}
		if submissionStat != "" {
			UpdateStudentSubmissionStatSQL.Exec(submissionStat, now, studentID, problemID)
		}
		if tutoringStat != "" {
			UpdateStudentTutoringStatSQL.Exec(tutoringStat, now, studentID, problemID)
		}
	} else {
		AddStudentStatusSQL.Exec(studentID, problemID, codingStat, helpStat, submissionStat, tutoringStat, time.Now())
	}
	rows.Close()
}

//-----------------------------------------------------------------
func init_teacher(id int, password string) {
	Teacher[id] = password
}

//-----------------------------------------------------------------
// initialize once per session
//-----------------------------------------------------------------
func init_student(student_id int, name string, password string) {
	AddAttendanceSQL.Exec(student_id, time.Now())

	BoardsSem.Lock()
	defer BoardsSem.Unlock()

	Students[student_id] = &StudenInfo{
		Name:                  name,
		Password:              password,
		Boards:                make([]*Board, 0),
		SubmissionStatus:      make([]*StudentSubmissionStatus, 0),
		SnapShotFeedbackQueue: make([]*SnapShotFeedback, 0),
		ThankStatus:           0,
	}

	// Student[student_id] = password
	// MessageBoards[student_id] = ""
	// Boards[student_id] = make([]*Board, 0)

	for i := 0; i < len(Students[0].Boards); i++ {
		b := &Board{
			Content:      Students[0].Boards[i].Content,
			Answer:       Students[0].Boards[i].Answer,
			Attempts:     Students[0].Boards[i].Attempts,
			Filename:     Students[0].Boards[i].Filename,
			Pid:          Students[0].Boards[i].Pid,
			StartingTime: time.Now(),
		}
		Students[student_id].Boards = append(Students[student_id].Boards, b)
	}
}

//-----------------------------------------------------------------
func load_and_authorize_student(student_id int, password string) bool {
	var pw, name string
	found := false
	rows, _ := Database.Query("select name, password from student where id=?", student_id)
	for rows.Next() {
		rows.Scan(&name, &pw)
		found = true
		break
	}
	rows.Close()
	if !found || pw != password {
		return false
	}
	init_student(student_id, name, password)
	return true
}

//-----------------------------------------------------------------
func load_teachers() {
	rows, _ := Database.Query("select id, password from teacher")
	defer rows.Close()
	var password string
	var id int
	for rows.Next() {
		rows.Scan(&id, &password)
		Teacher[id] = password
	}
	Passcode = RandStringRunes(20)
}
