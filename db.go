//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"database/sql"
	// "fmt"
	_ "github.com/mattn/go-sqlite3"
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
	execSQL("create table if not exists attendance (id integer primary key, student integer, at timestamp)")
	execSQL("create table if not exists problem (id integer primary key, tid integer, content blob, merit integer, effort integer, attempts integer, at timestamp)")
	execSQL("create table if not exists submission (id integer primary key, pid integer, sid integer, content blob, priority integer, at timestamp, completed timestamp)")

	// execSQL("create table if not exists score (id integer primary key, problem integer, student integer, merit_points integer, effort_points integer, attempts integer)")
	// execSQL("create table if not exists feedback (id integer primary key, teacher integer, student integer, content text, date timestamp)")

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
	AddAttendanceSQL = prepare("insert into attendance (student, at) values (?, ?)")
	AddProblemSQL = prepare("insert into problem (tid, content, merit, effort, attempts, at) values (?, ?, ?, ?, ?, ?)")
	AddSubmissionSQL = prepare("insert into submission (pid, sid, content, priority, at) values (?, ?, ?, ?, ?)")

	// Passcode for current session
	Passcode = RandStringRunes(12)
}

//-----------------------------------------------------------------
func load_students() {
	rows, _ := Database.Query("select id, password from student")
	defer rows.Close()
	var stid int
	var password string

	BoardsSem.Lock()
	defer BoardsSem.Unlock()

	for rows.Next() {
		rows.Scan(&stid, &password)
		Boards[stid] = make([]*Board, 0)
		Student[stid] = password
	}
	Boards[-1] = make([]*Board, 0) // content for a newly registered student
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
