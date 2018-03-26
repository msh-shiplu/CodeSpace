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
var AddSubmissionSQL *sql.Stmt
var AddFeedbackSQL *sql.Stmt
var AddScoreSQL *sql.Stmt
var UpdateScoreSQL *sql.Stmt

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
var SubSem sync.Mutex
var BulletinSem sync.Mutex

//---------------------------------------------------------
// Virtual boards for students and student submissions
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

var MessageBoards = make(map[int]string)

var BulletinBoard = make([]string, 0)

type Submission struct {
	Sid      int // submission id
	Uid      int // student id
	Pid      int // problem id
	Content  string
	Ext      string
	Priority int
	At       time.Time
}

var WorkingSubs = make([]*Submission, 0)

var NextProblem = make(map[int64]int64)

var ActiveProblems = make(map[int]struct{})

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
