package main

import (
	"fmt"
	"net/http"
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
func ta_registersHandler(w http.ResponseWriter, r *http.Request) {
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
		init_student(int(id), password)
		// Send password back to student
		fmt.Fprintf(w, fmt.Sprintf("%d,%s", id, password))
	}
}

//-----------------------------------------------------------------
