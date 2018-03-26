//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type BulletinBoardMessage struct {
	Code           string
	Idx            []string
	Tbr            int // page to be removed
	PC             string
	P1             int
	P2             int
	ActiveProblems int
	Attendance     int
	Authenticated  bool
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
func view_bulletin_boardHandler(w http.ResponseWriter, r *http.Request) {
	i, err := strconv.Atoi(r.FormValue("i"))
	passcode := r.FormValue("pc")
	if err != nil {
		i = 0
	}
	if i >= len(BulletinBoard) {
		i = 0
	}
	funcMap := template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
	}
	temp := template.New("")
	t, err2 := temp.Funcs(funcMap).Parse(TEACHER_MESSAGING_TEMPLATE)
	if err2 != nil {
		panic(err2)
	}

	// Get code and build page links
	idx := make([]string, 0)
	code := ""
	if i >= 0 && i < len(BulletinBoard) {
		for j := 0; j < len(BulletinBoard); j++ {
			if i == j {
				idx = append(idx, "active")
			} else {
				idx = append(idx, "")
			}
		}
		code = BulletinBoard[i]
	}

	// Get priority counts
	priority := []int{0, 0, 0}
	for i := 0; i < len(WorkingSubs); i++ {
		priority[WorkingSubs[i].Priority]++
	}

	data := &BulletinBoardMessage{
		Code:           code,
		Idx:            idx,
		Tbr:            i,
		PC:             passcode,
		P1:             priority[1],
		P2:             priority[2],
		ActiveProblems: len(ActiveProblems),
		Attendance:     len(Student),
		Authenticated:  passcode == Passcode,
	}
	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, &data)
}

//-----------------------------------------------------------------------------------
func student_messagesHandler(w http.ResponseWriter, r *http.Request) {
	stid, err := strconv.Atoi(r.FormValue("stid"))
	if err != nil {
		fmt.Fprintf(w, "Error")
	}
	mesg, ok := MessageBoards[stid]
	if ok {
		t := template.New("")
		t, err := t.Parse(STUDENT_MESSAGING_TEMPLATE)
		if err == nil {
			data := struct{ Message string }{mesg}
			w.Header().Set("Content-Type", "text/html")
			t.Execute(w, data)
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Fprint(w, "Error")
	}
}
