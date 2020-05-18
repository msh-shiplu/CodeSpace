//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

//-----------------------------------------------------------------------------------
// func insert_problems(uid int, problems []*ProblemInfo) {
//-----------------------------------------------------------------------------------
func insert_problem(uid int, problem *ProblemInfo) {
	// Create new problem
	pid := int64(0)
	if problem.Merit > 0 {
		// Find Tag id
		rows, _ := Database.Query("select id from tag where topic_description=?", problem.Tag)
		tagID := int64(0)
		for rows.Next() {
			rows.Scan(&tagID)
			break
		}
		rows.Close()
		if tagID == 0 {
			result, err := AddTagSQL.Exec(problem.Tag)
			if err != nil {
				fmt.Println(err)
			} else {
				tagID, _ = result.LastInsertId()
			}
		}

		// Insert only real problems into database
		result, err := AddProblemSQL.Exec(
			uid,
			problem.Description,
			problem.Answer,
			problem.Filename,
			problem.Merit,
			problem.Effort,
			problem.Attempts,
			problem.Topic_id,
			int(tagID),
			time.Now(),
		)
		if err != nil {
			log.Fatal(err)
		}
		pid, _ = result.LastInsertId()
		problem.Pid = int(pid)
		ActiveProblems[problem.Filename] = &ActiveProblem{
			Info:     problem,
			Answers:  make([]string, 0),
			Active:   true,
			Attempts: make(map[int]int),
		}
	}
}

//-----------------------------------------------------------------------------------
// Teacher starts one or more problems.
//-----------------------------------------------------------------------------------
func teacher_broadcastsHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	content := r.FormValue("content")
	answer := r.FormValue("answer")
	merit, _ := strconv.Atoi(r.FormValue("merit"))
	effort, _ := strconv.Atoi(r.FormValue("effort"))
	attempts, _ := strconv.Atoi(r.FormValue("attempts"))
	topic_id, _ := strconv.Atoi(r.FormValue("topic_id"))
	tag := r.FormValue("tag")
	filename := r.FormValue("filename")
	exact_answer := r.FormValue("exact_answer")

	// fmt.Printf("%d,Answer:%s, Merit:%d, Effort:%d, Attempts:%d, Tag:%s, Filename:%s\n", len(content), answer, merit, effort, attempts, tag, filename)

	problem := &ProblemInfo{
		Description: content,
		Filename:    filename,
		Answer:      answer,
		Merit:       merit,
		Effort:      effort,
		Attempts:    attempts,
		Topic_id:	topic_id,
		Tag:         tag,
		ExactAnswer: exact_answer == "True",
	}
	// fmt.Println("answer:", problem.Answer, problem.ExactAnswer)

	insert_problem(uid, problem)
	BoardsSem.Lock()
	defer BoardsSem.Unlock()
	for student_id, _ := range Students {
		b := &Board{
			Content:      problem.Description,
			Answer:       problem.Answer,
			Attempts:     problem.Attempts,
			Filename:     problem.Filename,
			Pid:          problem.Pid,
			StartingTime: time.Now(),
			Type:         "new",
		}
		Students[student_id].Boards = append(Students[student_id].Boards, b)
	}
	fmt.Fprintf(w, "Content copied to white boards.")
}

//-----------------------------------------------------------------------------------
