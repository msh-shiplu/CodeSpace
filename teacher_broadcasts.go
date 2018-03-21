//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//-----------------------------------------------------------------------------------
type ProblemFormat struct {
	Header      string
	Description string
	Answer      string
	Merit       int
	Effort      int
	Attempts    int
	Pid         int64
}

//-----------------------------------------------------------------------------------
func extract_problem_info(content, ext, answer_tag string) *ProblemFormat {
	var err error
	problem := &ProblemFormat{}
	merit, effort, attempts, answer := 0, 0, 0, ""
	prefix := "//"
	if ext != "java" && ext != "c++" && ext != "c" && ext != ".go" {
		prefix = "#"
	}
	content = strings.Trim(content, "\n ")
	if strings.HasPrefix(content, prefix) {
		items := strings.SplitN(content, "\n", 2)
		header := strings.Trim(items[0], "\n "+prefix)
		description := items[1]
		items = strings.SplitN(header, " ", 2)
		triple := items[0]
		if strings.Count(triple, ",") == 2 {
			items = strings.Split(triple, ",")
			merit, err = strconv.Atoi(items[0])
			if err != nil {
				return problem
			}
			effort, err = strconv.Atoi(items[1])
			if err != nil {
				return problem
			}
			attempts, err = strconv.Atoi(items[2])
			if err != nil {
				return problem
			}
			items := strings.SplitN(description, answer_tag, 2)
			if len(items) == 2 {
				answer = strings.Trim(items[1], "\n ")
				description = items[0] + "\n" + answer_tag + " "
			}
			problem = &ProblemFormat{
				Header:      prefix + " " + header,
				Description: description,
				Answer:      answer,
				Merit:       merit,
				Effort:      effort,
				Attempts:    attempts,
			}
		}
	}
	return problem
}

//-----------------------------------------------------------------------------------
func extract_problems(body, ext, answer_tag, divider_tag string) []*ProblemFormat {
	problems := make([]*ProblemFormat, 0)
	p := strings.Split(body, divider_tag)
	for i := 0; i < len(p); i++ {
		problems = append(problems, extract_problem_info(p[i], ext, answer_tag))
	}
	return problems
}

//-----------------------------------------------------------------------------------
func teacher_broadcastsHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	content, ext := r.FormValue("content"), r.FormValue("ext")
	divider_tag, answer_tag := r.FormValue("divider_tag"), r.FormValue("answer_tag")
	mode := r.FormValue("mode")
	problems := make([]*ProblemFormat, 0)
	// Extract info
	if mode == "unicast" {
		problems = append(problems, extract_problem_info(content, ext, answer_tag))
	} else {
		problems = extract_problems(content, ext, answer_tag, divider_tag)
	}

	// Create new problems
	for i := 0; i < len(problems); i++ {
		pid := int64(0)
		if problems[i].Merit > 0 {
			// insert only real problems into database
			content := fmt.Sprintf("%s\n%s %s\n",
				problems[i].Header,
				problems[i].Description,
				problems[i].Answer,
			)
			result, err := AddProblemSQL.Exec(
				uid,
				content,
				problems[i].Merit,
				problems[i].Effort,
				problems[i].Attempts,
				time.Now(),
			)
			if err != nil {
				panic(err)
			}
			pid, _ = result.LastInsertId()
			problems[i].Pid = pid
		}
	}

	BoardsSem.Lock()
	defer BoardsSem.Unlock()
	if mode == "unicast" || mode == "multicast_seq" {
		for stid, _ := range Boards {
			b := &Board{
				Content:      problems[0].Description,
				Answer:       problems[0].Answer,
				Attempts:     problems[0].Attempts,
				Ext:          ext,
				Pid:          int(problems[0].Pid),
				StartingTime: time.Now(),
			}
			Boards[stid] = append(Boards[stid], b)
		}
		if mode == "unicast" {
			fmt.Fprintf(w, "Content copied to white boards.")
		} else if mode == "multicast_seq" {
			for i := 0; i < len(problems)-1; i++ {
				NextProblem[problems[i].Pid] = problems[i+1].Pid
			}
			fmt.Fprintf(w, "First file copied to white boards.")
		}
	} else if mode == "multicast_or" {
		// Initialize random indices
		rand_idx := make([]int, len(Boards))
		j := 0
		for i := 0; i < len(Boards); i++ {
			rand_idx[i] = j
			j = (j + 1) % len(problems)
		}
		rand.Shuffle(len(rand_idx), func(i, j int) {
			rand_idx[i], rand_idx[j] = rand_idx[j], rand_idx[i]
		})
		// Insert into boards
		i := 0
		for stid, _ := range Boards {
			b := &Board{
				Content:      problems[rand_idx[i]].Description,
				Answer:       problems[rand_idx[i]].Answer,
				Attempts:     problems[rand_idx[i]].Attempts,
				Ext:          ext,
				Pid:          int(problems[rand_idx[i]].Pid),
				StartingTime: time.Now(),
			}
			Boards[stid] = append(Boards[stid], b)
			i++
		}
		fmt.Fprintf(w, "Content files randomly to white boards.")
	} else {
		fmt.Fprintf(w, "Unknown mode.")
	}
}

//-----------------------------------------------------------------------------------
