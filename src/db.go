//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
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
	execSQL("create table if not exists problem (id integer primary key, teacher_id integer, problem_description blob, answer text, filename text, merit integer, effort integer, attempts integer, tag integer, problem_uploaded_at timestamp)")
	execSQL("create table if not exists submission (id integer primary key, problem_id integer, student_id integer, student_code blob, submission_category integer, code_submitted_at timestamp, completed timestamp)")
	execSQL("create table if not exists score (id integer primary key, problem_id integer, student_id integer, teacher_id integer, score integer, graded_submission_number integer, score_given_at timestamp, unique(problem_id,student_id))")
	execSQL("create table if not exists feedback (id integer primary key, teacher_id integer, student_id integer, feedback text, feedback_given_at timestamp, problem_id integer)")
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
	AddProblemSQL = prepare("insert into problem (teacher_id, problem_description, answer, filename, merit, effort, attempts, tag, problem_uploaded_at) values (?, ?, ?, ?, ?, ?, ?, ?, ?)")
	AddSubmissionSQL = prepare("insert into submission (problem_id, student_id, student_code, submission_category, code_submitted_at) values (?, ?, ?, ?, ?)")
	AddSubmissionCompleteSQL = prepare("insert into submission (problem_id, student_id, student_code, submission_category, code_submitted_at, completed) values (?, ?, ?, ?, ?, ?)")
	CompleteSubmissionSQL = prepare("update submission set completed=? where id=?")
	AddScoreSQL = prepare("insert into score (problem_id, student_id, teacher_id, score, graded_submission_number, score_given_at) values (?, ?, ?, ?, ?, ?)")
	AddFeedbackSQL = prepare("insert into feedback (teacher_id, student_id, feedback, feedback_given_at, problem_id) values (?, ?, ?, ?, ?)")
	UpdateScoreSQL = prepare("update score set teacher_id=?, score=?, graded_submission_number=? where id=?")
	AddAttendanceSQL = prepare("insert into attendance (student_id, attendance_at) values (?, ?)")
	AddTagSQL = prepare("insert into tag (topic_description) values (?)")
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
	rows, _ := Database.Query("select id, score, attempts, teacher_id from score where problem_id=? and student_id=?", pid, student_id)
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
			mesg = fmt.Sprintf("Unable to add score: %d %d %d", pid, student_id, teacher_idß)
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

//-----------------------------------------------------------------
func init_teacher(id int, password string) {
	Teacher[id] = password
}

//-----------------------------------------------------------------
// initialize once per session
//-----------------------------------------------------------------
func init_student(student_id int, password string) {
	AddAttendanceSQL.Exec(student_id, time.Now())

	BoardsSem.Lock()
	defer BoardsSem.Unlock()

	Students[student_id] = &StudenInfo{
		Password:         password,
		Boards:           make([]*Board, 0),
		SubmissionStatus: 0,
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
	var pw string
	found := false
	rows, _ := Database.Query("select password from student where id=?", student_id)
	for rows.Next() {
		rows.Scan(&pw)
		found = true
		break
	}
	rows.Close()
	if !found || pw != password {
		return false
	}
	init_student(student_id, password)
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

//-----------------------------------------------------------------
