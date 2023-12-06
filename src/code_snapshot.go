package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func addCodeSnapshot(studentID int, problemID int, code string, status int, lastUpdate time.Time, event string) int {
	result, err := AddCodeSnapshotSQL.Exec(studentID, problemID, code, status, lastUpdate, event)
	if err != nil {
		log.Fatal("Could not save the snapshot for error: ", err)
		return -1
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
	return int(snapshotID)
}

func codeSnapshotHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	code := r.FormValue("code")
	problemID, _ := strconv.Atoi(r.FormValue("problem_id"))
	studentID, _ := strconv.Atoi(r.FormValue("uid"))
	snapshotEvent := r.FormValue("event")
	if snapshotEvent == "" {
		snapshotEvent = "at_regular_interval"
	}
	addCodeSnapshot(studentID, problemID, code, 0, time.Now(), snapshotEvent)
}

func codeSnapshotFeedbackHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	snapshotID, _ := strconv.Atoi(r.FormValue("snapshot_id"))
	feedback := r.FormValue("feedback")
	authorID, _ := strconv.Atoi(r.FormValue("uid"))
	authorRole := r.FormValue("role")
	now := time.Now()

	result, err := AddMessageSQL.Exec(snapshotID, "", authorID, authorRole, now, 1)
	if err != nil {
		log.Fatal("Could not save feedback for error: ", err)
		// fmt.Fprintf(w, "Could not save feedback")
		return
	}
	rows, err := Database.Query("select student_id, problem_id, code, filename from code_snapshot cs, problem p where cs.problem_id=p.id and cs.id=?", snapshotID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	messageID, _ := result.LastInsertId()
	studentID := -1
	code := ""
	filename := ""
	problemID := -1
	for rows.Next() {
		rows.Scan(&studentID, &problemID, &code, &filename)
	}
	rows.Close()
	if authorRole == "student" && studentID == authorID {
		fmt.Fprintf(w, "You can not give feedback to your own code.")
		return
	}
	idx := StudentSnapshot[studentID][problemID]
	Snapshots[idx].NumFeedback++
	result, err = AddMessageFeedbackSQL.Exec(messageID, feedback, authorID, authorRole, now)
	if err != nil {
		log.Fatal("Could not save feedback for error: ", err)
		// fmt.Fprintf(w, "Could not save feedback")
		return
	}
	feedbackID, _ := result.LastInsertId()
	Students[studentID].SnapShotFeedbackQueue = append(Students[studentID].SnapShotFeedbackQueue, &SnapShotFeedback{
		FeedbackID:  int(feedbackID),
		Snapshot:    code,
		Feedback:    feedback,
		ProblemName: filename,
		Provider:    getName(uid, authorRole),
	})
	addOrUpdateStudentStatus(studentID, problemID, "", "Been helped", "", "")
	if authorRole == "student" {
		addOrUpdateStudentStatus(authorID, problemID, "", "", "", "Tutoring")
	}
	fmt.Println("Feedback on code snapshot saved!")
	// http.Redirect(w, r, "/get_codespace?uid="+strconv.Itoa(authorID)+"&role="+authorRole+"&password="+r.FormValue("password"), http.StatusSeeOther)
}

func messageFeedbackHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	messageID, _ := strconv.Atoi(r.FormValue("message_id"))
	feedback := r.FormValue("feedback")
	authorID, _ := strconv.Atoi(r.FormValue("uid"))
	authorRole := r.FormValue("role")
	now := time.Now()

	result, err := AddMessageFeedbackSQL.Exec(messageID, feedback, authorID, authorRole, now)
	if err != nil {
		log.Fatal("Could not save the feedback for error: ", err)
		// fmt.Fprintf(w, "Could not save the feedback.")
		return
	}
	rows, err := Database.Query("select student_id, problem_id, code, filename, m.type from code_snapshot cs, problem p, message m where cs.problem_id=p.id and m.snapshot_id = cs.id and m.id=?", messageID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	studentID := -1
	code := ""
	filename := ""
	problemID := -1
	messageType := -1
	for rows.Next() {
		rows.Scan(&studentID, &problemID, &code, &filename, &messageType)
	}
	rows.Close()
	if studentID == authorID {
		fmt.Fprintf(w, "You can not give feedback to your own code.")
		return
	}
	if messageType == 0 {
		addOrUpdateStudentStatus(studentID, problemID, "", "Been helped", "", "")
		if authorRole == "student" {
			addOrUpdateStudentStatus(authorID, problemID, "", "", "", "Tutoring")
		}
	}
	idx := StudentSnapshot[studentID][problemID]
	Snapshots[idx].NumFeedback++
	feedbackID, _ := result.LastInsertId()
	Students[studentID].SnapShotFeedbackQueue = append(Students[studentID].SnapShotFeedbackQueue, &SnapShotFeedback{
		FeedbackID:  int(feedbackID),
		Snapshot:    code,
		Feedback:    feedback,
		ProblemName: filename,
		Provider:    getName(uid, authorRole),
	})
	fmt.Println("Feedback on message saved!")
	// http.Redirect(w, r, "/get_codespace?uid="+strconv.Itoa(authorID)+"&role="+authorRole+"&password="+r.FormValue("password"), http.StatusSeeOther)
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
