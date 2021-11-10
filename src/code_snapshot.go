package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func addCodeSnapshot(studentID int, problemID int, code string, status int, lastUpdate time.Time) int {
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
			return -1
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
			Status:      status,
			FirstUpdate: lastUpdate,
			LastUpdated: lastUpdate,
			LinesOfCode: getLinesOfCode(code),
			Code:        code,
		})
	} else {
		currentStatus := Snapshots[idx].Status
		if currentStatus > status {
			status = currentStatus
		}
		Snapshots[idx] = &Snapshot{
			ID:          int(snapshotID),
			StudentName: Snapshots[idx].StudentName,
			StudentID:   studentID,
			ProblemName: Snapshots[idx].ProblemName,
			ProblemID:   problemID,
			Status:      status,
			FirstUpdate: Snapshots[idx].FirstUpdate,
			LastUpdated: lastUpdate,
			LinesOfCode: getLinesOfCode(code),
			Code:        code,
			NumFeedback: Snapshots[idx].NumFeedback,
		}
	}
	fmt.Println("Code snapshot saved!")
	return int(snapshotID)
}

func codeSnapshotHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	code := r.FormValue("code")
	problemID, _ := strconv.Atoi(r.FormValue("problem_id"))
	studentID, _ := strconv.Atoi(r.FormValue("uid"))
	addCodeSnapshot(studentID, problemID, code, 0, time.Now())
}

func codeSnapshotFeedbackHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	snapshotID, _ := strconv.Atoi(r.FormValue("snapshot_id"))
	feedback := r.FormValue("feedback")
	authorID, _ := strconv.Atoi(r.FormValue("uid"))
	authorRole := r.FormValue("role")
	now := time.Now()
	rows, err := Database.Query("select student_id, problem_id, code, filename from code_snapshot cs, problem p where cs.problem_id=p.id and cs.id=?", snapshotID)
	if err != nil {
		log.Fatal(err)
	}
	// defer rows.Close()
	studentID := -1
	code := ""
	filename := ""
	problemID := -1
	for rows.Next() {
		rows.Scan(&studentID, &problemID, &code, &filename)
	}
	rows.Close()
	result, err := AddSnapShotFeedbackSQL.Exec(snapshotID, feedback, authorID, authorRole, now)
	if err != nil {
		log.Fatal(err)
	}
	idx := StudentSnapshot[studentID][problemID]
	Snapshots[idx].NumFeedback++
	feedbackID, _ := result.LastInsertId()
	Students[studentID].SnapShotFeedbackQueue = append(Students[studentID].SnapShotFeedbackQueue, &SnapShotFeedback{
		FeedbackID:  int(feedbackID),
		Snapshot:    code,
		Feedback:    feedback,
		ProblemName: filename,
	})
	fmt.Println("Feedback on code snapshot saved!")
	http.Redirect(w, r, "/get_codespace?uid="+strconv.Itoa(authorID)+"&role="+authorRole+"&password="+r.FormValue("password"), http.StatusSeeOther)
}

func getSnapshotFeedbackHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	feedfback := Students[uid].SnapShotFeedbackQueue[0]
	Students[uid].SnapShotFeedbackQueue = Students[uid].SnapShotFeedbackQueue[1:]
	js, err := json.Marshal(feedfback)
	if err != nil {
		log.Fatal(err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}
