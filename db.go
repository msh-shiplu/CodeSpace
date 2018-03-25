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
	execSQL("create table if not exists problem (id integer primary key, tid integer, content blob, answer text, ext text, merit integer, effort integer, attempts integer, at timestamp)")
	execSQL("create table if not exists submission (id integer primary key, pid integer, sid integer, content blob, priority integer, at timestamp, completed timestamp)")
	execSQL("create table if not exists score (id integer primary key, pid integer, stid integer, points integer, attempts integer, unique(pid,stid))")
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
	AddProblemSQL = prepare("insert into problem (tid, content, answer, ext, merit, effort, attempts, at) values (?, ?, ?, ?, ?, ?, ?, ?)")
	AddSubmissionSQL = prepare("insert into submission (pid, sid, content, priority, at) values (?, ?, ?, ?, ?)")
	AddScoreSQL = prepare("insert into score (pid, stid, points, attempts) values (?, ?, ?, ?)")
	AddFeedbackSQL = prepare("insert into feedback (tid, stid, content, date) values (?, ?, ?, ?)")
	UpdateScoreSQL = prepare("update score set points=?, attempts=? where id=?")
	AddAttendanceSQL = prepare("insert into attendance (stid, at) values (?, ?)")

	// Initialize passcode for current session and default board
	Passcode = RandStringRunes(12)
	Boards[-1] = make([]*Board, 0) // content for a newly logged in student
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
			Ext:          Boards[-1][i].Ext,
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
