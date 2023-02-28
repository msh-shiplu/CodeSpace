//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	// "encoding/json"
	"fmt"
	"log"
	"time"

	// "log"
	"net/http"
)

//-----------------------------------------------------------------------------------
// When problems are deactivated, boards cleared, no new submissions are possibile.
//-----------------------------------------------------------------------------------
func teacher_deactivates_problemsHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	CodeSnapshotSem.Lock()
	defer CodeSnapshotSem.Unlock()
	filename := r.FormValue("filename")
	if prob, ok := ActiveProblems[filename]; ok {
		prob.Active = false
		PeerTutorAllowed = false
		if len(prob.Answers) > 0 {
			fmt.Fprintf(w, "1")
		} else {
			fmt.Fprintf(w, "0")
		}
		tempSnapshots := Snapshots
		for studentID, _ := range StudentSnapshot {
			StudentSnapshot[studentID] = map[int]int{}
		}
		Snapshots = make([]*Snapshot, 0)
		idx := 0
		for _, s := range tempSnapshots {
			if s.ProblemID != prob.Info.Pid {
				Snapshots = append(Snapshots, s)
				for studentID, _ := range StudentSnapshot {
					StudentSnapshot[studentID][prob.Info.Pid] = idx
				}
				idx++
			}
		}
		_, err := UpdateProblemEndTimeSQL.Exec(time.Now(), prob.Info.Pid)
		if err != nil {
			log.Fatal(err)
		}
		for studentID, _ := range Students {
			for i, b := range Students[studentID].Boards {
				if b.Pid == prob.Info.Pid {
					Students[studentID].Boards = append(Students[studentID].Boards[:i], Students[studentID].Boards[i+1:]...)
					break
				}
			}
		}

	} else {
		fmt.Fprintf(w, "-1")
	}

	// filenames := make([]string, 0)
	// for fname, prob := range ActiveProblems {
	// 	if prob.Active {
	// 		prob.Active = false
	// 		if len(prob.Answers) > 0 {
	// 			filenames = append(filenames, fname)
	// 		}
	// 	}
	// }
	// for stid, _ := range Students {
	// 	Students[stid].Boards = make([]*Board, 0)
	// 	Students[stid].SubmissionStatus = 0
	// }
	// js, err := json.Marshal(filenames)
	// if err == nil {
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.Write(js)
	// } else {
	// 	log.Fatal(err)
	// }
}

//-----------------------------------------------------------------------------------
// Clear submissions, boards, statuses, and set all problems inactive.
//-----------------------------------------------------------------------------------
func teacher_clears_submissionsHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	WorkingSubs = make([]*Submission, 0)
	fmt.Fprintf(w, "Done.")
}
