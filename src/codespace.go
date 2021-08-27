package main

import (
	"bufio"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Snapshot contains information related to a code snapshot.
type Snapshot struct {
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

// CodeSpace contains a list of code snapshot information.
type CodeSpace struct {
	Snapshots []*Snapshot
}

var snapshotStatus = []string{"Not submitted", "Submitted: not graded", "Submitted: incorrect", "Submitted: correct"}

func getStudentInfo() map[int]string {
	var students = make(map[int]string)
	rows, err := Database.Query("select id, name from student")
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	id, name := -1, ""
	for rows.Next() {
		rows.Scan(&id, &name)
		students[id] = name
	}
	return students
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

func getCodeSpace() []*Snapshot {
	snapshots := make([]*Snapshot, 0)
	students := getStudentInfo()
	for _, problem := range ActiveProblems {
		if problem.Active == false {
			continue
		}
		rows, err := Database.Query("select student_id, status, code, last_updated_at, starting_time from code_snapshot where problem_id=?", problem.Info.Pid)
		if err != nil {
			log.Fatal(err)
		}
		studentID, status, code, lastUpdatedAt, startingTime := 0, 0, "", time.Now(), time.Now()
		for rows.Next() {
			rows.Scan(&studentID, &status, &code, &lastUpdatedAt, &startingTime)
			timeSpent := lastUpdatedAt.Sub(startingTime)
			loc := getLinesOfCode(code)
			s := &Snapshot{
				StudentName: students[studentID],
				StudentID:   studentID,
				ProblemName: problem.Info.Filename,
				ProblemID:   problem.Info.Pid,
				Status:      snapshotStatus[status],
				TimeSpent:   timeSpent,
				LastUpdated: lastUpdatedAt,
				LinesOfCode: loc,
			}
			snapshots = append(snapshots, s)
		}
		rows.Close()
	}
	return snapshots
}

func codespaceHandler(w http.ResponseWriter, r *http.Request) {
	temp := template.New("")
	t, err := temp.Parse(CODESPACE_TEMPLATE)
	if err != nil {
		log.Fatal(err)
	}
	data := getCodeSpace()
	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, data)
}

func getCodeSnapshot(studentID int, problemID int) *Snapshot {
	rows, err := Database.Query("select name from student where id=?", studentID)
	if err != nil {
		log.Fatal(err)
	}
	name := ""
	for rows.Next() {
		rows.Scan(&name)
	}
	rows.Close()
	problemName := ""
	for _, problem := range ActiveProblems {
		if problem.Active == true && problem.Info.Pid == problemID {
			problemName = problem.Info.Filename
			break
		}
	}
	var s *Snapshot
	if problemName == "" {
		s = &Snapshot{}
	} else {
		rows, err = Database.Query("select status, code, last_updated_at, starting_time from code_snapshot where problem_id=? and student_id=?", problemID, studentID)
		defer rows.Close()
		if err != nil {
			log.Fatal(err)
		}

		status, code, lastUpdatedAt, startingTime := 0, "", time.Now(), time.Now()
		for rows.Next() {
			rows.Scan(&status, &code, &lastUpdatedAt, &startingTime)
			timeSpent := lastUpdatedAt.Sub(startingTime)
			s = &Snapshot{
				StudentName: name,
				ProblemName: problemName,
				Status:      snapshotStatus[status],
				TimeSpent:   timeSpent,
				LastUpdated: lastUpdatedAt,
				Code:        code,
			}
		}
	}
	return s
}

func getCodeSnapshotHandler(w http.ResponseWriter, r *http.Request) {
	studentID, _ := strconv.Atoi(r.FormValue("student_id"))
	problemID, _ := strconv.Atoi(r.FormValue("problem_id"))
	temp := template.New("")
	t, err := temp.Parse(CODE_SNAPSHOT_TEMPLATE)
	if err != nil {
		log.Fatal(err)
	}
	data := getCodeSnapshot(studentID, problemID)
	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, data)
}
