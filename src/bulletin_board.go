//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

//-----------------------------------------------------------------------------------
type BulletinBoardMessage struct {
	Code           string
	I              int
	NextI          int
	PrevI          int
	PC             string
	P1             int
	P2             int
	ActiveProblems string
	// ActiveProblems int
	BulletinItems int
	AnswerCount   int
	Attendance    int
	Address       string
	Authenticated bool
}

//-----------------------------------------------------------------------------------
func teacher_adds_bulletin_pageHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	BulletinSem.Lock()
	defer BulletinSem.Unlock()
	BulletinBoard = append(BulletinBoard, r.FormValue("content"))
	fmt.Fprintf(w, "Content added to bulletin board")
}

//-----------------------------------------------------------------
func remove_bulletin_pageHandler(w http.ResponseWriter, r *http.Request) {
	BulletinSem.Lock()
	defer BulletinSem.Unlock()
	i, _ := strconv.Atoi(r.FormValue("i"))
	passcode := r.FormValue("pc")
	if passcode == Passcode && i >= 0 && i < len(BulletinBoard) {
		BulletinBoard = append(BulletinBoard[:i], BulletinBoard[i+1:]...)
		http.Redirect(w, r, "view_bulletin_board?i=0&pc="+passcode, http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "view_bulletin_board?i="+r.FormValue("i")+"&pc="+passcode, http.StatusSeeOther)
	}
}

//-----------------------------------------------------------------------------------
func get_bulletin_board_data(i int, passcode string) *BulletinBoardMessage {
	BulletinSem.Lock()
	defer BulletinSem.Unlock()

	if i >= len(BulletinBoard) {
		i = 0
	}
	// Get code and build page links
	code := ""
	if i >= 0 && i < len(BulletinBoard) {
		code = BulletinBoard[i]
	}

	// Get priority counts
	priority := []int{0, 0, 0}
	for j := 0; j < len(WorkingSubs); j++ {
		priority[WorkingSubs[j].Priority]++
	}
	next_i, prev_i := 0, 0
	if len(BulletinBoard) > 0 {
		next_i = (i + 1 + len(BulletinBoard)) % len(BulletinBoard)
		prev_i = (i - 1 + len(BulletinBoard)) % len(BulletinBoard)
	}
	answers := 0
	keys := make([]string, 0)
	for key := range ActiveProblems {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	submissions := make([]string, 0)
	for i, key := range keys {
		p := ActiveProblems[key]
		if p.Active {
			rows, err := Database.Query("select problem_uploaded_at from problem where id = ?", p.Info.Pid)
			if err != nil {
				fmt.Println("Error retrieving problem starting time", err)
				return &BulletinBoardMessage{}
			}
			var starting_time time.Time
			for rows.Next() {
				rows.Scan(&starting_time)
			}
			duration := time.Since(starting_time).Minutes()
			subs := len(p.Attempts)
			label := fmt.Sprintf("P%d: %d subs after %.0fm", i+1, subs, duration)
			submissions = append(submissions, label)
			answers += len(p.Answers)
		}
	}
	active_problems := strings.Join(submissions, ". ")
	fmt.Println(">", active_problems)

	// for pid, p := range ActiveProblems {
	// 	if p.Active {
	// 		answers += len(p.Answers)
	// 	}
	// 	fmt.Println(">", pid, p)
	// 	cur_submissions = append(cur_submissions, len(p.Attempts))
	// }

	data := &BulletinBoardMessage{
		Code:           code,
		I:              i,
		NextI:          next_i,
		PrevI:          prev_i,
		PC:             passcode,
		P1:             priority[1],
		P2:             priority[2],
		ActiveProblems: active_problems,
		// ActiveProblems: len(ActiveProblems),
		BulletinItems: len(BulletinBoard),
		AnswerCount:   answers,
		Attendance:    len(Students),
		Address:       Config.Address,
		Authenticated: passcode == Passcode,
	}
	return data
}

//-----------------------------------------------------------------------------------
func bulletin_board_dataHandler(w http.ResponseWriter, r *http.Request) {
	data := get_bulletin_board_data(0, "")
	js, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(js)
}

//-----------------------------------------------------------------------------------
func view_bulletin_boardHandler(w http.ResponseWriter, r *http.Request) {
	i, err := strconv.Atoi(r.FormValue("i"))
	passcode := r.FormValue("pc")
	if err != nil {
		i = 0
	}

	temp := template.New("")
	t, err2 := temp.Parse(TEACHER_MESSAGING_TEMPLATE)
	if err2 != nil {
		log.Fatal(err2)
	}
	data := get_bulletin_board_data(i, passcode)
	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, data)
}

//-----------------------------------------------------------------------------------
// func student_messagesHandler(w http.ResponseWriter, r *http.Request) {
// 	stid, err := strconv.Atoi(r.FormValue("stid"))
// 	if err != nil {
// 		fmt.Fprintf(w, "Error")
// 	}
// 	_, ok := Students[stid]
// 	if ok {
// 		t := template.New("")
// 		t, err := t.Parse(STUDENT_MESSAGING_TEMPLATE)
// 		if err == nil {
// 			data := struct{ Message string }{Students[stid].Status}
// 			w.Header().Set("Content-Type", "text/html")
// 			t.Execute(w, data)
// 		} else {
// 			fmt.Println(err)
// 		}
// 	} else {
// 		fmt.Fprint(w, "Error")
// 	}
// }
