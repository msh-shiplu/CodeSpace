package main

import (
	"fmt"
	"net/http"
	"time"
)

//-----------------------------------------------------------------
func teacher_adds_taHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	mesg := ""
	rows, err := Database.Query("select name from teacher where name=?", name)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		mesg = fmt.Sprintf("%s already exists. Choose a different name.", name)
		return
	}
	password := RandStringRunes(12)
	result, _ := AddTeacherSQL.Exec(name, password)
	id, _ := result.LastInsertId()
	Teacher[int(id)] = password
	mesg = fmt.Sprintf("%s is added. User must register under the same name", name)
	fmt.Fprintf(w, mesg)
}

//-----------------------------------------------------------------
func teacher_registersHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	rows, err := Database.Query("select id, password from teacher where name=?", name)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var password string
	var id int
	for rows.Next() {
		rows.Scan(&id, &password)
		fmt.Fprintf(w, fmt.Sprintf("%d,%s", id, password))
		return
	}
	fmt.Fprintf(w, "Failed")
}

//-----------------------------------------------------------------
func student_registersHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	rows, err := Database.Query("select name from student where name=?", name)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	i := 0
	for rows.Next() {
		i++
	}
	if i > 0 {
		fmt.Fprintf(w, "exist")
	} else {
		password := RandStringRunes(12)
		result, err2 := AddStudentSQL.Exec(name, password)
		if err2 != nil {
			panic(err2)
		}
		id, err3 := result.LastInsertId()
		if err3 != nil {
			panic(err3)
		}
		Student[int(id)] = password

		// Initialize student's board
		BoardsSem.Lock()
		defer BoardsSem.Unlock()
		Boards[int(id)] = make([]*Board, 0)
		for i := 0; i < len(Boards[-1]); i++ {
			b := &Board{
				Content:      Boards[-1][i].Content,
				Answer:       Boards[-1][i].Answer,
				Attempts:     Boards[-1][i].Attempts,
				Ext:          Boards[-1][i].Ext,
				Pid:          Boards[-1][i].Pid,
				StartingTime: time.Now(),
			}
			Boards[int(id)] = append(Boards[int(id)], b)
		}

		// Send password back to student
		fmt.Fprintf(w, fmt.Sprintf("%d,%s", id, password))
	}
}

//-----------------------------------------------------------------
