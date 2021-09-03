package main

import (
	"log"
	"net/http"
	"strconv"
	"time"
)

func addCodeSnapshot(studentID int, problemID int, code string, status int, lastUpdate time.Time) {
	result, err := AddCodeSnapshotSQL.Exec(studentID, problemID, code, status, lastUpdate)
	if err != nil {
		log.Fatal(err)
	}
	snapshotID, _ := result.LastInsertId()
	idx, ok := StudentSnapshot[studentID][problemID]
	if !ok {
		idx = len(Snapshots)
		StudentSnapshot[studentID][problemID] = idx
		rows, err := Database.Query("select name from student where id=?", studentID)
		defer rows.Close()
		if err != nil {
			log.Fatal(err)
			return
		}
		name := ""
		if rows.Next() {
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
		Snapshots = append(Snapshots, &Snapshot{
			ID:          int(snapshotID),
			StudentName: name,
			StudentID:   studentID,
			ProblemName: problemName,
			ProblemID:   problemID,
			Status:      SnapshotStatus[status],
			TimeSpent:   time.Duration(0),
			LastUpdated: lastUpdate,
			LinesOfCode: getLinesOfCode(code),
			Code:        code,
		})
	} else {
		currentStatus := SnapshotStatusMapping[Snapshots[idx].Status]
		if currentStatus > status {
			status = currentStatus
		}
		Snapshots[idx] = &Snapshot{
			ID:          int(snapshotID),
			StudentName: Snapshots[idx].StudentName,
			StudentID:   studentID,
			ProblemName: Snapshots[idx].ProblemName,
			ProblemID:   problemID,
			Status:      SnapshotStatus[status],
			TimeSpent:   Snapshots[idx].TimeSpent + (lastUpdate.Sub(Snapshots[idx].LastUpdated)),
			LastUpdated: lastUpdate,
			LinesOfCode: getLinesOfCode(code),
			Code:        code,
		}
	}

}

func codeSnapshotHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	code := r.FormValue("code")
	problemID, _ := strconv.Atoi(r.FormValue("problem_id"))
	studentID, _ := strconv.Atoi(r.FormValue("uid"))
	addCodeSnapshot(studentID, problemID, code, 0, time.Now())
}

func codeSnapshotFeedbackHandler(w http.ResponseWriter, r *http.Request) {
	snapshotID, _ := strconv.Atoi(r.FormValue("snapshot_id"))
	feedback := r.FormValue("feedback")
	authorID, _ := strconv.Atoi(r.FormValue("uid"))
	authorRole := r.FormValue("role")
	now := time.Now()
	AddSnapShotFeedbackSQL.Exec(snapshotID, feedback, authorID, authorRole, now)

}
