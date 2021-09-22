//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"bufio"
	"database/sql"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

//---------------------------------------------------------
type Configuration struct {
	CourseId   string
	NameServer string
	IP         string
	Port       int
	Database   string
	Address    string
	LogFile    string
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
var AddTagSQL *sql.Stmt
var AddTestCaseSQL *sql.Stmt
var UpdateTestCaseSQL *sql.Stmt
var AddHelpSubmissionSQL *sql.Stmt
var AddHelpMessageSQL *sql.Stmt
var UpdateHelpMessageSQL *sql.Stmt
var AddCodeSnapshotSQL *sql.Stmt
var AddSnapShotFeedbackSQL *sql.Stmt
var UpdateSnapshotFeedbackSQL *sql.Stmt

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
var HelpSubSem sync.Mutex
var CodeSnapshotSem sync.Mutex

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
	Type         string
}
type StudentSubmissionStatus struct {
	Filename      string
	AttemptNumber int
	Status        int
	/*
		1 submission being looked at.
		2 teacher did not grade your submission (dismissed).
		3 your submission was not correct.
		4 your submission was correct.
	*/
}
type SnapShotFeedback struct {
	FeedbackID  int
	Snapshot    string
	Feedback    string
	ProblemName string
}
type StudenInfo struct {
	Password              string
	Boards                []*Board
	SubmissionStatus      []*StudentSubmissionStatus
	SnapShotFeedbackQueue []*SnapShotFeedback
	ThankStatus           int
	/*
		0 Nothing
		1 Got a new thanks for feedback
	*/
}

var Students = make(map[int]*StudenInfo)

//---------------------------------------------------------

var BulletinBoard = make([]string, 0)

//---------------------------------------------------------
type Submission struct {
	Sid           int // submission id
	Uid           int // student id
	Pid           int // problem id
	Content       string
	Filename      string
	Priority      int
	AttemptNumber int
	At            time.Time
	Name          string
}

var WorkingSubs = make([]*Submission, 0)
var Submissions = make(map[int]*Submission)

//---------------------------------------------------------

type HelpSubmission struct {
	Sid      int // submission id
	Uid      int // student id
	Pid      int // problem id
	Status   int // 0=ok, 1=queue empty, 2=not elligible
	Content  string
	Filename string
	At       time.Time
}

var WorkingHelpSubs = make([]*HelpSubmission, 0)
var HelpSubmissions = make(map[int]*HelpSubmission)

//---------------------------------------------------------
type ProblemInfo struct {
	Description string
	Filename    string
	Answer      string
	Merit       int
	Effort      int
	Attempts    int
	Topic_id    int
	Tag         string
	Pid         int
	ExactAnswer bool
}

type ActiveProblem struct {
	Info     *ProblemInfo
	Answers  []string
	Active   bool
	Attempts map[int]int
}

var ActiveProblems = make(map[string]*ActiveProblem)

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

//-----------------------------------------------------------------------------
func writeLog(filename, message string) {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.Println(time.Now(), " ", message)
}

//---------------------------------------------------------

var HelpEligibleStudents = map[int]map[int]bool{}
var SeenHelpSubmissions = map[int]map[int]bool{}

//---------------------------------------------------------

// Snapshot contains information related to a code snapshot.
type Snapshot struct {
	ID          int
	StudentName string
	StudentID   int
	ProblemName string
	ProblemID   int
	Status      string
	TimeSpent   time.Duration
	LastUpdated time.Time
	LinesOfCode int
	Code        string
}

// Snapshots contains all the current snapshots from students.
var Snapshots = make([]*Snapshot, 0)

// StudentSnapshot is the mapping from (student id, problem id) -> snashpt index in `Snapshots` list.
var StudentSnapshot = map[int]map[int]int{}

// SnapshotStatus maps from integer status to string named status.
var SnapshotStatus = []string{"Not submitted", "Submitted: not graded", "Submitted: incorrect", "Submitted: correct"}

// SnapshotStatusMapping maps from string snapshot status to integer.
var SnapshotStatusMapping = map[string]int{
	"Not submitted":         0,
	"Submitted: not graded": 1,
	"Submitted: incorrect":  2,
	"Submitted: correct":    3,
}

func getLinesOfCode(code string) int {
	scanner := bufio.NewScanner(strings.NewReader(code))
	scanner.Split(bufio.ScanLines)
	count := 0
	for scanner.Scan() {
		count++
	}
	return count
}
