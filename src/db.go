// Author: Vinhthuy Phan, 2018
package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func create_tables() {
	execSQL := func(s string) {
		sql_stmt, err := Database.Prepare(s)
		if err != nil {
			log.Fatal(err)
		}
		sql_stmt.Exec()
	}
	execSQL("create table if not exists student (id INT AUTO_INCREMENT NOT NULL, name VARCHAR(100) unique, password VARCHAR(100), PRIMARY KEY (`id`))")
	execSQL("create table if not exists teacher (id INT AUTO_INCREMENT NOT NULL, name VARCHAR(100) unique, password VARCHAR(100), PRIMARY KEY (`id`))")
	execSQL("create table if not exists attendance (id INT AUTO_INCREMENT NOT NULL, student_id INT NOT NULL, attendance_at timestamp, PRIMARY KEY (`id`))")
	execSQL("create table if not exists tag (id INT AUTO_INCREMENT NOT NULL, topic_description VARCHAR(200) unique, PRIMARY KEY (`id`))")
	execSQL("create table if not exists problem (id INT AUTO_INCREMENT NOT NULL, teacher_id INT, problem_description text, answer text, filename text, merit INT, effort INT, attempts INT, topic_id INT, tag INT, problem_uploaded_at timestamp, problem_ended_at timestamp, PRIMARY KEY (`id`))")
	execSQL("create table if not exists submission (id INT AUTO_INCREMENT NOT NULL, problem_id INT NOT NULL, student_id INT NOT NULL, student_code text, snapshot_id INT default 0, submission_category INT, code_submitted_at timestamp, completed timestamp, verdict text, attempt_number INT, answer text, PRIMARY KEY (`id`))")
	execSQL("create table if not exists score (id INT AUTO_INCREMENT NOT NULL, problem_id INT NOT NULL, student_id INT, teacher_id INT, score INT, graded_submission_number INT, score_given_at timestamp, unique(problem_id,student_id), PRIMARY KEY (`id`))")
	execSQL("create table if not exists feedback (id INT AUTO_INCREMENT NOT NULL, teacher_id INT, student_id INT, feedback text, feedback_given_at timestamp, submission_id INT, PRIMARY KEY (`id`))")
	execSQL("create table if not exists test_case (id INT AUTO_INCREMENT NOT NULL, problem_id INT, student_id INT, test_cases text, added_at timestamp, PRIMARY KEY (`id`))")
	execSQL("create table if not exists code_explanation (id INT AUTO_INCREMENT NOT NULL, problem_id INT, student_id INT, snapshot_id INT, trying_what text, need_help_with text, code_submitted_at timestamp, PRIMARY KEY (`id`))")
	execSQL("create table if not exists help_message (id INT AUTO_INCREMENT NOT NULL, code_explanation_id INT, student_id INT, message text, given_at timestamp, useful text, updated_at timestamp, PRIMARY KEY (`id`))")
	execSQL("create table if not exists code_snapshot (id INT AUTO_INCREMENT NOT NULL, student_id INT, problem_id INT, code text, last_updated_at timestamp, status int default 0, event VARCHAR(50), PRIMARY KEY (`id`))") // 0 = not submitted, 1 = submitted but not graded, 2 = submitted and incorrect, 3 = submitted and correct
	execSQL("create table if not exists snapshot_feedback (id INT AUTO_INCREMENT NOT NULL, snapshot_id INT, feedback text, author_id INT, author_role VARCHAR(50), given_at timestamp, PRIMARY KEY (`id`))")
	execSQL("create table if not exists snapshot_back_feedback (id INT AUTO_INCREMENT NOT NULL, snapshot_feedback_id INT, author_id INT, author_role VARCHAR(50), is_helpful VARCHAR(50), given_at timestamp, PRIMARY KEY (`id`))")
	execSQL("create table if not exists message (id INT AUTO_INCREMENT NOT NULL, snapshot_id INT, message text, author_id INT, author_role VARCHAR(50), given_at timestamp, type INT, PRIMARY KEY (`id`))")
	execSQL("create table if not exists message_feedback (id INT AUTO_INCREMENT NOT NULL, message_id INT, feedback text, author_id INT, author_role VARCHAR(50), given_at timestamp, PRIMARY KEY (`id`))")
	execSQL("create table if not exists message_back_feedback (id INT AUTO_INCREMENT NOT NULL, message_feedback_id INT, author_id INT, author_role VARCHAR(50), useful VARCHAR(50), given_at timestamp, PRIMARY KEY (`id`))")
	execSQL("create table if not exists help_eligible (id INT AUTO_INCREMENT NOT NULL, problem_id INT, student_id INT, became_eligible_at timestamp, PRIMARY KEY (`id`))")
	execSQL("create table if not exists user_event_log (id INT AUTO_INCREMENT NOT NULL, name VARCHAR(50), user_id INT, user_type VARCHAR(50), event_type VARCHAR(50), referral_info VARCHAR(50), event_time timestamp, PRIMARY KEY (`id`))")
	execSQL("create table if not exists student_status (id INT AUTO_INCREMENT NOT NULL, student_id INT, problem_id INT, coding_stat VARCHAR(50), help_stat VARCHAR(50), submission_stat VARCHAR(50), tutoring_stat VARCHAR(50), last_updated_at timestamp, PRIMARY KEY (`id`))")
	execSQL("create table if not exists problem_statistics (id INT AUTO_INCREMENT NOT NULL, problem_id INT not null, active INT default 0, submission INT default 0, help_request INT default 0, graded_correct INT default 0, graded_incorrect INT default 0, PRIMARY KEY (`id`))")
	// foreign key example: http://www.sqlitetutorial.net/sqlite-foreign-key/
}

//-----------------------------------------------------------------
func init_database(db_name string, username string, pass string, server string) {
	var err error
	prepare := func(s string) *sql.Stmt {
		stmt, err := Database.Prepare(s)
		if err != nil {
			log.Fatal(err)
		}
		return stmt
	}

	Database, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/", username, pass, server))
	if err != nil {
		log.Fatal(err)
	}
	_, err = Database.Exec("CREATE DATABASE IF NOT EXISTS " + db_name)
	if err != nil {
		log.Fatal(err)
	}
	// _, err = Database.Exec("USE " + db_name)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	Database, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", username, pass, server, db_name))
	if err != nil {
		log.Fatal(err)
	}
	Database.SetConnMaxLifetime(time.Minute * 3)
	create_tables()
	AddStudentSQL = prepare("insert into student (name, password) values (?, ?)")
	AddTeacherSQL = prepare("insert into teacher (name, password) values (?, ?)")
	AddProblemSQL = prepare("insert into problem (teacher_id, problem_description, answer, filename, merit, effort, attempts, topic_id, tag, problem_uploaded_at) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	AddSubmissionSQL = prepare("insert into submission (problem_id, student_id, student_code, submission_category, attempt_number, code_submitted_at, snapshot_id, answer) values (?, ?, ?, ?, ?, ?, ?, ?)")
	AddSubmissionCompleteSQL = prepare("insert into submission (problem_id, student_id, student_code, submission_category, attempt_number, code_submitted_at, completed, snapshot_id, answer) values (?, ?, ?, ?, ?, ?, ?, ?, ?)")
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
	AddCodeSnapshotSQL = prepare("insert into code_snapshot (student_id, problem_id, code, status, last_updated_at, event) values(?, ?, ?, ?, ?, ?)")
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
	AddMessageBackFeedbackSQL = prepare("insert into message_back_feedback (message_feedback_id, author_id, author_role, useful, given_at) values(?, ?, ?, ?, ?)")
	UpdateMessageBackFeedbackSQL = prepare("update message_back_feedback set useful=?, given_at=? where message_feedback_id=? and author_id=? and author_role=?")
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

func databaseTransaction(stmt *sql.Stmt, args ...any) (sql.Result, error) {
	tx, err := Database.Begin()
	if err != nil {
		log.Fatal(err)
	}
	result, err := tx.Stmt(stmt).Exec(args...)
	if err != nil {
		fmt.Println("doing rollback")
		tx.Rollback()
	} else {
		tx.Commit()
	}
	return result, err
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
		rows.Close()
		if codingStat != "" {
			_, err = UpdateStudentCodingStatSQL.Exec(codingStat, now, studentID, problemID)
			if err != nil {
				log.Fatal(err)
			}
		}
		if helpStat != "" {
			_, err = UpdateStudentHelpStatSQL.Exec(helpStat, now, studentID, problemID)
			if err != nil {
				log.Fatal(err)
			}
		}
		if submissionStat != "" {
			_, err = UpdateStudentSubmissionStatSQL.Exec(submissionStat, now, studentID, problemID)
			if err != nil {
				log.Fatal(err)
			}
		}
		if tutoringStat != "" {
			_, err = UpdateStudentTutoringStatSQL.Exec(tutoringStat, now, studentID, problemID)
			if err != nil {
				log.Fatal(err)
			}
		}
	} else {
		rows.Close()
		_, err = AddStudentStatusSQL.Exec(studentID, problemID, codingStat, helpStat, submissionStat, tutoringStat, time.Now())
		if err != nil {
			log.Fatal(err)
		}
	}
}

//-----------------------------------------------------------------
func init_teacher(id int, name string, password string) {
	Teacher[id] = password
	TeacherPass[name] = password
	TeacherNameToId[name] = id
	TeacherIdToName[id] = name
	SeenHelpSubmissions[id] = map[int]bool{}
}

//-----------------------------------------------------------------
// initialize once per session
//-----------------------------------------------------------------
func init_student(student_id int, name string, password string) {
	_, err := AddAttendanceSQL.Exec(student_id, time.Now())
	if err != nil {
		log.Fatal(err)
	}

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
	StudentSnapshot[student_id] = map[int]int{}
}

//-----------------------------------------------------------------
func load_and_authorize_student(student_id int, password string) bool {
	var pw, name string
	found := false
	rows, err := Database.Query("select name, password from student where id=?", student_id)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		rows.Scan(&name, &pw)
		found = true
		break
	}
	if !found || pw != password {
		return false
	}
	init_student(student_id, name, password)
	return true
}

//-----------------------------------------------------------------
func load_teachers() {
	rows, _ := Database.Query("select id,name, password from teacher")
	defer rows.Close()
	var password string
	var name string
	var id int
	for rows.Next() {
		rows.Scan(&id, &name, &password)
		Teacher[id] = password
		TeacherPass[name] = password
		TeacherNameToId[name] = id
		TeacherIdToName[id] = name
	}
	Passcode = RandStringRunes(20)
}
