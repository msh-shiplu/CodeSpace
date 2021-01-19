//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

//-----------------------------------------------------------------------------------
func extract_partial_credits(content string) int {
	re := regexp.MustCompile(`(\d)+ for effort`)
	result := re.FindSubmatch([]byte(content))
	if len(result) >= 2 {
		points, _ := strconv.Atoi(string(result[1]))
		return points
	} else {
		return -1
	}
}

//-----------------------------------------------------------------------------------
func teacher_gradesHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	content, decision := r.FormValue("content"), r.FormValue("decision")
	sid, _ := strconv.Atoi(r.FormValue("sid"))
	changed := r.FormValue("changed")
	// student_id, _ := strconv.Atoi(r.FormValue("student_id"))
	// pid, _ := strconv.Atoi(r.FormValue("pid"))
	mesg := ""
	sub, ok := Submissions[sid]
	if !ok {
		fmt.Fprintf(w, "Unknown submission cannot be graded.")
		return
	}
	student_id := sub.Uid
	if changed == "True" {
		// If the original file is changed, there's feedback.  Copy it to whiteboard.
		if prob, ok := ActiveProblems[sub.Filename]; ok {
			AddFeedbackSQL.Exec(uid, student_id, content, time.Now(), sub.Sid)
			mesg = "Feedback saved to student's board."
			BoardsSem.Lock()
			defer BoardsSem.Unlock()
			b := &Board{
				Content:      content,
				Answer:       prob.Info.Answer,
				Attempts:     0, // This tells the client this is an existing problem
				Filename:     sub.Filename,
				Pid:          sub.Pid,
				StartingTime: time.Now(),
				Type:         "feedback",
			}
			Students[student_id].Boards = append(Students[student_id].Boards, b)
		}
	}

	// If submission is dismissed, do not take that attempt away from the student.
	if decision == "dismissed" {
		Students[student_id].SubmissionStatus = 2
		ActiveProblems[sub.Filename].Attempts[student_id] += 1
		fmt.Fprintf(w, "Submission dismissed.")
	} else if decision == "ungraded" {
		Students[student_id].SubmissionStatus = 5
		fmt.Fprintf(w, mesg)
	} else {
		// Update score based on the grading decision
		partial_credits := -1
		if decision != "correct" {
			partial_credits = extract_partial_credits(content)
		}
		scoring_mesg := add_or_update_score(decision, sub.Pid, sub.Uid, uid, partial_credits)
		mesg = scoring_mesg + "\n" + mesg
		if decision == "correct" {
			Students[sub.Uid].SubmissionStatus = 4
			ActiveProblems[sub.Filename].Attempts[student_id] = 0 // This prevents further submission.
			pid := sub.Pid
			if _, ok := HelpEligibleStudents[pid][sub.Uid]; !ok {
				HelpEligibleStudents[pid][sub.Uid] = true
				SeenHelpSubmissions[sub.Uid] = map[int]bool{}
				mesg = mesg + "\nYou are now elligible to help you friends. To help please click on 'Help Friends' button."
			}
		} else {
			Students[student_id].SubmissionStatus = 3
		}

		// Update submission complete time
		_, err := CompleteSubmissionSQL.Exec(time.Now(), decision, sid)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, mesg)
	}
}

//-----------------------------------------------------------------------------------
