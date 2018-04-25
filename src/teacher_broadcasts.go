//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//-----------------------------------------------------------------------------------
func extract_problems(content, answers, merits, efforts, attempts, tags, filenames, divider string) []*ProblemInfo {
	if divider == "" {
		merit, _ := strconv.Atoi(merits)
		effort, _ := strconv.Atoi(efforts)
		attempt, _ := strconv.Atoi(attempts)
		return []*ProblemInfo{&ProblemInfo{
			Description: content,
			Filename:    filenames,
			Answer:      answers,
			Merit:       merit,
			Effort:      effort,
			Attempts:    attempt,
			Tag:         tags,
		}}
	}
	c := strings.Split(content, divider)
	an := strings.Split(answers, "\n")
	m := strings.Split(merits, "\n")
	ef := strings.Split(efforts, "\n")
	at := strings.Split(attempts, "\n")
	tg := strings.Split(tags, "\n")
	fn := strings.Split(filenames, "\n")
	problems := make([]*ProblemInfo, 0)
	for i := 0; i < len(c); i++ {
		merit, _ := strconv.Atoi(m[i])
		effort, _ := strconv.Atoi(ef[i])
		attempt, _ := strconv.Atoi(at[i])
		p := &ProblemInfo{
			Description:     strings.TrimLeft(c[i], " \n"),
			Filename:        fn[i],
			Answer:          an[i],
			Merit:           merit,
			Effort:          effort,
			Attempts:        attempt,
			Tag:             tg[i],
			NextIfCorrect:   0,
			NextIfIncorrect: 0,
		}
		problems = append(problems, p)
		// fmt.Println(p)
	}
	return problems
}

//-----------------------------------------------------------------------------------
func assign_next_problem_pid(problems []*ProblemInfo, next_if_correct, next_if_incorrect string) {
	nic := strings.Split(next_if_correct, "\n")
	nii := strings.Split(next_if_incorrect, "\n")
	for i := 0; i < len(problems)-1; i++ {
		if nic[i] != "-1" {
			idx, err := strconv.Atoi(nic[i])
			if err != nil {
				log.Fatal(err)
			}
			problems[i].NextIfCorrect = problems[idx].Pid
		}
		if nii[i] != "-1" {
			idx, err := strconv.Atoi(nii[i])
			if err != nil {
				log.Fatal(err)
			}
			problems[i].NextIfIncorrect = problems[idx].Pid
		}
	}
}

//-----------------------------------------------------------------------------------
func insert_problems(uid int, problems []*ProblemInfo) {
	// Create new problems
	for i := 0; i < len(problems); i++ {
		pid := int64(0)
		if problems[i].Merit > 0 {
			// Find Tag id
			rows, _ := Database.Query("select id from tag where description=?", problems[i].Tag)
			tagID := int64(0)
			for rows.Next() {
				rows.Scan(&tagID)
				break
			}
			rows.Close()
			if tagID == 0 {
				result, err := AddTagSQL.Exec(problems[i].Tag)
				if err != nil {
					fmt.Println(err)
				} else {
					tagID, _ = result.LastInsertId()
				}
			}

			// Insert only real problems into database
			result, err := AddProblemSQL.Exec(
				uid,
				problems[i].Description,
				problems[i].Answer,
				problems[i].Filename,
				problems[i].Merit,
				problems[i].Effort,
				problems[i].Attempts,
				int(tagID),
				time.Now(),
			)
			if err != nil {
				log.Fatal(err)
			}
			pid, _ = result.LastInsertId()
			problems[i].Pid = int(pid)
			ActiveProblems[int(pid)] = &ActiveProblem{
				Info:     problems[i],
				Answers:  make([]string, 0),
				Active:   true,
				Attempts: make(map[int]int),
			}
		}
	}
}

//-----------------------------------------------------------------------------------
// Teacher starts one or more problems.
//-----------------------------------------------------------------------------------
func teacher_broadcastsHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	content := r.FormValue("content")
	answers := r.FormValue("answers")
	merits := r.FormValue("merits")
	efforts := r.FormValue("efforts")
	attempts := r.FormValue("attempts")
	tags := r.FormValue("tags")
	filenames := r.FormValue("filenames")
	divider := r.FormValue("divider")
	mode := r.FormValue("mode")
	nic, nii := r.FormValue("nic"), r.FormValue("nii")

	problems := make([]*ProblemInfo, 0)

	// Deactivate active problems and clear student boards
	for _, prob := range ActiveProblems {
		prob.Active = false
	}
	for stid, _ := range Students {
		Students[stid].Boards = make([]*Board, 0)
		Students[stid].SubmissionStatus = 0
	}

	// Extract info
	problems = extract_problems(content, answers, merits, efforts, attempts, tags, filenames, divider)

	insert_problems(uid, problems)

	if mode == "multicast_seq" {
		assign_next_problem_pid(problems, nic, nii)
	}

	BoardsSem.Lock()
	defer BoardsSem.Unlock()
	if mode == "unicast" || mode == "multicast_seq" || mode == "multicast_and" {
		end := 1
		if mode == "multicast_and" {
			end = len(problems)
		}
		for stid, _ := range Students {
			for i := 0; i < end; i++ {
				b := &Board{
					Content:      problems[i].Description,
					Answer:       problems[i].Answer,
					Attempts:     problems[i].Attempts,
					Filename:     problems[i].Filename,
					Pid:          problems[i].Pid,
					StartingTime: time.Now(),
				}
				Students[stid].Boards = append(Students[stid].Boards, b)
			}
		}
		fmt.Fprintf(w, "Content copied to white boards.")
	} else if mode == "multicast_or" {
		// Initialize random indices
		rand_idx := make([]int, len(Students))
		j := 0
		for i := 0; i < len(Students); i++ {
			rand_idx[i] = j
			j = (j + 1) % len(problems)
		}
		rand.Shuffle(len(rand_idx), func(i, j int) {
			rand_idx[i], rand_idx[j] = rand_idx[j], rand_idx[i]
		})
		// Insert into boards
		i := 0
		for stid, _ := range Students {
			b := &Board{
				Content:      problems[rand_idx[i]].Description,
				Answer:       problems[rand_idx[i]].Answer,
				Attempts:     problems[rand_idx[i]].Attempts,
				Filename:     problems[rand_idx[i]].Filename,
				Pid:          int(problems[rand_idx[i]].Pid),
				StartingTime: time.Now(),
			}
			Students[stid].Boards = append(Students[stid].Boards, b)
			i++
		}
		fmt.Fprintf(w, "Files saved randomly to white boards.")
	} else {
		fmt.Fprintf(w, "Unknown mode.")
	}
}

//-----------------------------------------------------------------------------------
