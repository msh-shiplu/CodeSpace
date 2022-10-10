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
	submissionID, _ := strconv.Atoi(r.FormValue("submission_id"))
	message := r.FormValue("message")
	res, err := AddHelpMessageSQL.Exec(submissionID, uid, message, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	messageID, _ := res.LastInsertId()
	// student_id := 0
	// rows, _ := Database.Query("select student_id from code_explanation where id=?", submission_id)
	// for rows.Next() {
	// 	rows.Scan(&student_id)
	// 	break
	// }
	// rows.Close()
	helpSub := HelpSubmissions[submissionID]
	studentID := helpSub.Uid
	message = helpSub.Content + "\n\nFeedback: " + message
	b := &Board{
		Content:      message,
		Answer:       "",
		Attempts:     0,
		Filename:     "peer_feedback.txt",
		Pid:          int(messageID),
		StartingTime: time.Now(),
		Type:         "peer_feedback",
	}
	Students[studentID].Boards = append(Students[studentID].Boards, b)
	fmt.Fprint(w, "Dear "+who+", Your feedback has been sent.")

}
func sendThankYouHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	messageID, _ := strconv.Atoi(r.FormValue("message_id"))
	useful := r.FormValue("useful")
	_, err := UpdateHelpMessageSQL.Exec(useful, time.Now(), messageID)
	if err != nil {
		log.Fatal(err)
	}
	if useful == "yes" {
		studentID := 0
		rows, _ := Database.Query("select student_id from help_message where id=?", messageID)
		for rows.Next() {
			rows.Scan(&studentID)
			break
		}
		rows.Close()
		Students[studentID].ThankStatus = 1
	}

}

func studentSendBackFeedbackHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	backFeedback := r.FormValue("feedback")
	feedbackID, _ := strconv.Atoi(r.FormValue("feedback_id"))
	authorRole := r.FormValue("role")
	rows, err := Database.Query("select * from message_back_feedback where message_feedback_id = ? and author_id = ? and author_role = ?", feedbackID, uid, authorRole)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	if rows.Next() {
		rows.Close()
		_, err = UpdateMessageBackFeedbackSQL.Exec(backFeedback, time.Now(), feedbackID, uid, authorRole)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		rows.Close()
		_, err = AddMessageBackFeedbackSQL.Exec(feedbackID, uid, authorRole, backFeedback, time.Now())
		if err != nil {
			log.Fatal(err)
		}
	}
}
