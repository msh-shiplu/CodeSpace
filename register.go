package main

import (
	"fmt"
	"net/http"
	"time"
)

//-----------------------------------------------------------------
func setup_new_teacherHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	mesg := ""
	rows, err := Database.Query("select name from teacher where name=?", name)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	i := 0
	for rows.Next() {
		i++
	}
	if i > 0 {
		mesg = fmt.Sprintf("%s already exists. Choose a different name.", name)
	} else {
		Teacher[name] = RandStringRunes(12)
		mesg = fmt.Sprintf("%s is temporarily added. This person needs to register under the same name", name)
	}
	fmt.Fprintf(w, mesg)
}

//-----------------------------------------------------------------
func register_teacherHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if password, ok := Teacher[name]; ok {
		result, err := AddTeacherSQL.Exec(name, password)
		if err != nil {
			fmt.Fprintf(w, "exist")
		} else {
			id, err2 := result.LastInsertId()
			if err2 != nil {
				panic(err2)
			}
			fmt.Fprintf(w, fmt.Sprintf("%d,%s", id, password))
		}
	} else {
		fmt.Fprintf(w, "notsetup")
	}
}

//-----------------------------------------------------------------
func register_studentHandler(w http.ResponseWriter, r *http.Request) {
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
		Student[name] = password

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
