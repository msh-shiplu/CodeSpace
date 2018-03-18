package main

import (
	"database/sql"
	"math/rand"
	"sync"
	"time"
)

const VERSION = "0.1"

//---------------------------------------------------------
// Database
//---------------------------------------------------------

var Database *sql.DB
var AddStudentSQL *sql.Stmt
var AddTeacherSQL *sql.Stmt
var AddAttendanceSQL *sql.Stmt
var AddProblemSQL *sql.Stmt

//---------------------------------------------------------
// Authentication
//---------------------------------------------------------

var Student = make(map[int]string)
var Teacher = make(map[int]string)
var Passcode string

//---------------------------------------------------------
// Semaphores
//---------------------------------------------------------

var BoardsSem sync.Mutex

//---------------------------------------------------------
// Virtual boards for students
//---------------------------------------------------------
type Board struct {
	Content      string
	Answer       string
	Attempts     int
	Ext          string
	Pid          int // problem id
	StartingTime time.Time
}

var Boards = make(map[int][]*Board)

// type BroadcastData struct {
// 	Content  string `json:"content"`
// 	Answer   string `json:"answer"`
// 	Merit    int    `json:"merit"`
// 	Effort   int    `json:"effort"`
// 	Attempts int    `json:"attempts"`
// 	Ext      string `json:"ext"`
// }

//---------------------------------------------------------
// Utilities
//---------------------------------------------------------

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

//---------------------------------------------------------
