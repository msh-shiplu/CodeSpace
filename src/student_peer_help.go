//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

//-----------------------------------------------------------------------------------
func studentGetHelpCode(w http.ResponseWriter, r *http.Request, who string, uid int) {
	filename := r.FormValue("filename")
	pid := 0
	prob, ok := ActiveProblems[filename]
	HelpSubSem.Lock()
	defer HelpSubSem.Unlock()
	selected := &HelpSubmission{}
	selected.Status = 1
	if ok {
		if prob.Active {
			// fmt.Fprint(w, "This problem is not active.")
			pid = prob.Info.Pid
			if _, ok := HelpEligibleStudents[pid][uid]; ok {

				for idx, sub := range WorkingHelpSubs {
					if sub.Pid != pid || sub.Uid == uid {
						continue
					}
					if _, ok := SeenHelpSubmissions[uid][sub.Sid]; !ok {
						selected = sub
						WorkingHelpSubs = append(WorkingHelpSubs[:idx], WorkingHelpSubs[idx+1:]...)
						SeenHelpSubmissions[uid][sub.Sid] = true
						selected.Status = 0
						break
					}
				}

				// fmt.Fprintf(w, "You are elligible to help in this problem.")

			} else {
				selected.Status = 2
			}
		}
	}
	js, err := json.Marshal(selected)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}

}

//-----------------------------------------------------------------------------------

func student_return_without_feedbackHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	HelpSubSem.Lock()
	defer HelpSubSem.Unlock()
	tmp := r.FormValue("submission_id")
	submissionID, _ := strconv.Atoi(tmp)
	submission := HelpSubmissions[submissionID]
	WorkingHelpSubs = append(WorkingHelpSubs, submission)
	fmt.Fprint(w, "No feedback is given. This request is returned to the help queue.")
}

func student_send_help_messageHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	submission_id, _ := strconv.Atoi(r.FormValue("submission_id"))
	message := r.FormValue("message")
	res, err := AddHelpMessageSQL.Exec(submission_id, uid, message, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	message_id, _ := res.LastInsertId()
	// student_id := 0
	// rows, _ := Database.Query("select student_id from help_submission where id=?", submission_id)
	// for rows.Next() {
	// 	rows.Scan(&student_id)
	// 	break
	// }
	// rows.Close()
	helpSub := HelpSubmissions[submission_id]
	student_id := helpSub.Uid
	message = helpSub.Content + "\n\nFeedback: " + message
	b := &Board{
		Content:      message,
		Answer:       "",
		Attempts:     0,
		Filename:     "peer_feedback.txt",
		Pid:          int(message_id),
		StartingTime: time.Now(),
		Type:         "peer_feedback",
	}
	Students[student_id].Boards = append(Students[student_id].Boards, b)
	fmt.Fprint(w, "Dear "+who+", Your feedback has been sent.")

}
func sendThankYouHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	message_id, _ := strconv.Atoi(r.FormValue("message_id"))
	useful := r.FormValue("useful")
	_, err := UpdateHelpMessageSQL.Exec(useful, time.Now(), message_id)
	if err != nil {
		log.Fatal(err)
	}
	if useful == "yes" {
		studentID := 0
		rows, _ := Database.Query("select student_id from help_message where id=?", message_id)
		for rows.Next() {
			rows.Scan(&studentID)
			break
		}
		Students[studentID].ThankStatus = 1
	}

}
