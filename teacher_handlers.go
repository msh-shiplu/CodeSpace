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
	Description string
	Answer      string
}

//-----------------------------------------------------------------------------------
func extract_problems(body, problem_divider, answer_tag string) []*ProblemFormat {
	problems := make([]*ProblemFormat, 0)
	p := strings.Split(body, problem_divider)
	for i := 0; i < len(p); i++ {
		items := strings.SplitN(p[i], answer_tag, 2)
		answer := ""
		if len(items) == 2 {
			answer = items[1]
		}
		problems = append(problems, &ProblemFormat{Description: items[0], Answer: answer})
	}
	return problems
}

//-----------------------------------------------------------------------------------
func extract_problem_info(content, ext, problem_divider, answer_tag string) ([]*ProblemFormat, int, int, int) {
	var err error
	problems := []*ProblemFormat{&ProblemFormat{Description: content}}
	merit, effort, attempts := 0, 0, 0
	prefix := "//"
	if ext != "java" && ext != "c++" && ext != "c" && ext != ".go" {
		prefix = "#"
	}
	if strings.HasPrefix(content, prefix) {
		items := strings.SplitN(content, "\n", 2)
		header := strings.Trim(items[0], "\n "+prefix)
		body := items[1]
		items = strings.SplitN(header, " ", 2)
		triple := items[0]
		if strings.Count(triple, ",") == 2 {
			items = strings.Split(triple, ",")
			merit, err = strconv.Atoi(items[0])
			if err != nil {
				return problems, 0, 0, 0
			}
			effort, err = strconv.Atoi(items[1])
			if err != nil {
				return problems, 0, 0, 0
			}
			attempts, err = strconv.Atoi(items[2])
			if err != nil {
				return problems, 0, 0, 0
			}
			problems = extract_problems(body, problem_divider, answer_tag)
		}
	}
	return problems, merit, effort, attempts
}

//-----------------------------------------------------------------------------------
func teacher_broadcastHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	BoardsSem.Lock()
	defer BoardsSem.Unlock()

	content, ext := r.FormValue("content"), r.FormValue("ext")
	problem_divider, answer_tag := r.FormValue("problem_divider"), r.FormValue("answer_tag")

	// Extract info
	problems, merit, effort, attempts := extract_problem_info(content, ext, problem_divider, answer_tag)

	// Create new problem
	pid := int64(0)
	if merit > 0 {
		result, err := AddProblemSQL.Exec(uid, content, merit, effort, attempts, time.Now())
		if err != nil {
			panic(err)
		}
		pid, _ = result.LastInsertId()
	}

	// Initialize random indices
	rand_idx := make([]int, len(Boards))
	j := 0
	for i := 0; i < len(Boards); i++ {
		rand_idx[i] = j
		j = (j + 1) % len(problems)
	}
	if len(problems) > 1 {
		rand.Shuffle(len(rand_idx), func(i, j int) {
			rand_idx[i], rand_idx[j] = rand_idx[j], rand_idx[i]
		})
	}

	// Insert into boards
	i := 0
	for stid, _ := range Boards {
		b := &Board{
			Content:      problems[rand_idx[i]].Description,
			Answer:       problems[rand_idx[i]].Answer,
			Attempts:     attempts,
			Ext:          ext,
			Pid:          int(pid),
			StartingTime: time.Now(),
		}
		Boards[stid] = append(Boards[stid], b)
		i++
	}
	fmt.Fprintf(w, "Content copied to white boards.")
}

//-----------------------------------------------------------------------------------
func teacher_get_passcodeHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	fmt.Fprintf(w, Passcode)
}

//-----------------------------------------------------------------------------------
func testHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	// Show content of boards
	fmt.Println("Boards", len(Boards))
	for uid, board_pages := range Boards {
		fmt.Printf("Uid: %d has %d pages.\n", uid, len(board_pages))
		for i := 0; i < len(board_pages); i++ {
			b := board_pages[i]
			fmt.Printf("Attempts: %d, Ext: %s, Pid: %d, Answer: %s\n",
				b.Attempts, b.Ext, b.Pid, b.Answer)
			fmt.Println(b.Content)
			fmt.Println("********")
		}
		fmt.Println("--------------------------------")
	}
	fmt.Fprintf(w, Passcode)
}

//-----------------------------------------------------------------------------------
