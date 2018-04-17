//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"database/sql"
	"math/rand"
	"sync"
	"time"
)

const VERSION = "0.1"

//---------------------------------------------------------
type Configuration struct {
	Id         string
	NameServer string
	IP         string
	Port       int
	Database   string
	Address    string
}

var Config *Configuration

//---------------------------------------------------------
// Database
//---------------------------------------------------------

var Database *sql.DB
var AddStudentSQL *sql.Stmt
var AddTeacherSQL *sql.Stmt
var AddAttendanceSQL *sql.Stmt
var AddProblemSQL *sql.Stmt
var AddSubmissionSQL *sql.Stmt
var AddSubmissionCompleteSQL *sql.Stmt
var CompleteSubmissionSQL *sql.Stmt
var AddFeedbackSQL *sql.Stmt
var AddScoreSQL *sql.Stmt
var UpdateScoreSQL *sql.Stmt

//---------------------------------------------------------
// Authentication
//---------------------------------------------------------

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
	Filename     string
	Pid          int // problem id
	StartingTime time.Time
}

type StudenInfo struct {
	Password         string
	Boards           []*Board
	SubmissionStatus int
	/*
		1 submission being looked at.
		2 teacher did not grade your submission (dismissed).
		3 your submission was not correct.
		4 your submission was correct.
	*/
}

var Students = make(map[int]*StudenInfo)

//---------------------------------------------------------

var BulletinBoard = make([]string, 0)

//---------------------------------------------------------
type Submission struct {
	Sid      int // submission id
	Uid      int // student id
	Pid      int // problem id
	Content  string
	Filename string
	Priority int
	At       time.Time
}

var WorkingSubs = make([]*Submission, 0)

//---------------------------------------------------------
type ProblemInfo struct {
	Description     string
	Filename        string
	Answer          string
	Merit           int
	Effort          int
	Attempts        int
	Pid             int
	NextIfCorrect   int
	NextIfIncorrect int
}

type ActiveProblem struct {
	Info     *ProblemInfo
	Answers  []string
	Active   bool
	Attempts map[int]int
}

var ActiveProblems = make(map[int]*ActiveProblem)

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
