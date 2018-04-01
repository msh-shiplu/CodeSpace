//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"database/sql"
	// "fmt"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

func create_tables() {
	execSQL := func(s string) {
		sql_stmt, err := Database.Prepare(s)
		if err != nil {
			panic(err)
		}
		sql_stmt.Exec()
	}
	execSQL("create table if not exists student (id integer primary key, name text unique, password text)")
	execSQL("create table if not exists teacher (id integer primary key, name text unique, password text)")
	execSQL("create table if not exists problem (id integer primary key, tid integer, content blob, answer text, filename text, merit integer, effort integer, attempts integer, at timestamp)")
	execSQL("create table if not exists submission (id integer primary key, pid integer, sid integer, content blob, priority integer, at timestamp, completed timestamp)")
	execSQL("create table if not exists score (id integer primary key, pid integer, stid integer, tid integer, points integer, attempts integer, at timestamp, unique(pid,stid))")
	execSQL("create table if not exists feedback (id integer primary key, tid integer, stid integer, content text, date timestamp)")
	execSQL("create table if not exists attendance (id integer primary key, stid integer, at timestamp)")

	// foreign key example: http://www.sqlitetutorial.net/sqlite-foreign-key/
}

//-----------------------------------------------------------------
func init_database(db_name string) {
	var err error
	prepare := func(s string) *sql.Stmt {
		stmt, err := Database.Prepare(s)
		if err != nil {
			panic(err)
		}
		return stmt
	}

	Database, err = sql.Open("sqlite3", db_name)
	if err != nil {
		panic(err)
	}
	create_tables()
	AddStudentSQL = prepare("insert into student (name, password) values (?, ?)")
	AddTeacherSQL = prepare("insert into teacher (name, password) values (?, ?)")
	AddProblemSQL = prepare("insert into problem (tid, content, answer, filename, merit, effort, attempts, at) values (?, ?, ?, ?, ?, ?, ?, ?)")
	AddSubmissionSQL = prepare("insert into submission (pid, sid, content, priority, at) values (?, ?, ?, ?, ?)")
	AddScoreSQL = prepare("insert into score (pid, stid, tid, points, attempts, at) values (?, ?, ?, ?, ?, ?)")
	AddFeedbackSQL = prepare("insert into feedback (tid, stid, content, date) values (?, ?, ?, ?)")
	UpdateScoreSQL = prepare("update score set tid=?, points=?, attempts=? where id=?")
	AddAttendanceSQL = prepare("insert into attendance (stid, at) values (?, ?)")

	// Initialize passcode for current session and default board
	Passcode = RandStringRunes(12)
	Boards[-1] = make([]*Board, 0) // content for a newly logged in student
}

//-----------------------------------------------------------------
func add_next_problem_to_board(pid, stid int) string {
	prob, ok := ActiveProblems[pid]
	if ok && prob.Next > 0 {
		new_content, new_answer, new_fn, new_merit, new_effort, new_attempts := "", "", "", 0, 0, 0
		rows, _ := Database.Query("select content, answer, filename, merit, effort, attempts from problem where id=?", prob.Next)
		for rows.Next() {
			rows.Scan(&new_content, &new_answer, &new_fn, &new_merit, &new_effort, &new_attempts)
			break
		}
		rows.Close()
		b := &Board{
			Content:      new_content,
			Answer:       new_answer,
			Attempts:     new_attempts,
			Filename:     new_fn,
			Pid:          int(prob.Next),
			StartingTime: time.Now(),
		}
		Boards[stid] = append(Boards[stid], b)
		return "\nNew problem added to white board."
	}
	return ""
}

//-----------------------------------------------------------------
// Add or update score based on a decision. If decision is "correct"
// a new problem, if there's one, is added to student's board.
//-----------------------------------------------------------------
func add_or_update_score(decision string, pid, stid, tid int) string {
	mesg := ""

	// Find score information for this student (stid) for this problem (pid)
	score_id, current_points, current_attempts, current_tid := 0, 0, 0, 0
	rows, _ := Database.Query("select id, points, attempts, tid from score where pid=? and stid=?", pid, stid)
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
	points, teacher := 0, tid
	if decision == "correct" {
		points = merit
		mesg = "Answer is correct."
		m := add_next_problem_to_board(pid, stid)
		mesg = mesg + m
	} else {
		points = effort
		// If the problem was previously graded correct, this submission
		// does not reduce it.  Grading is asynchronous.
		if points < current_points {
			points = current_points
			teacher = current_tid
		}
		mesg = "Answer is incorrect."
	}

	// Add a new score or update a current score for this student & problem
	if score_id == 0 {
		_, err := AddScoreSQL.Exec(pid, stid, tid, points, current_attempts+1, time.Now())
		if err != nil {
			panic(err)
		}
	} else {
		_, err := UpdateScoreSQL.Exec(teacher, points, current_attempts+1, score_id)
		if err != nil {
			panic(err)
		}
	}
	return mesg
}

//-----------------------------------------------------------------
// initialize once per session
//-----------------------------------------------------------------
func init_student(stid int, password string) {
	AddAttendanceSQL.Exec(stid, time.Now())

	BoardsSem.Lock()
	defer BoardsSem.Unlock()
	Student[stid] = password
	MessageBoards[stid] = ""
	Boards[stid] = make([]*Board, 0)
	for i := 0; i < len(Boards[-1]); i++ {
		b := &Board{
			Content:      Boards[-1][i].Content,
			Answer:       Boards[-1][i].Answer,
			Attempts:     Boards[-1][i].Attempts,
			Filename:     Boards[-1][i].Filename,
			Pid:          Boards[-1][i].Pid,
			StartingTime: time.Now(),
		}
		Boards[stid] = append(Boards[stid], b)
	}
}

//-----------------------------------------------------------------
func load_and_authorize_student(stid int, password string) bool {
	var pw string
	found := false
	rows, _ := Database.Query("select password from student where id=?", stid)
	for rows.Next() {
		rows.Scan(&pw)
		found = true
		break
	}
	rows.Close()
	if !found || pw != password {
		return false
	}
	init_student(stid, password)
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
