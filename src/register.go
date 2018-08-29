package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

//-----------------------------------------------------------------
func add_multiple(filename, role string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		name := strings.TrimSpace(scanner.Text())
		if name != "" {
			add_user(name, role)
		}
	}
}

//-----------------------------------------------------------------
func add_user(name, role string) {
	var err error
	var rows *sql.Rows
	var result sql.Result
	var id int64

	if role == "teacher" {
		rows, err = Database.Query("select name from teacher where name=?", name)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		rows, err = Database.Query("select name from student where name=?", name)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer rows.Close()
	for rows.Next() {
		fmt.Printf("%s already exists. Choose a different name.\n", name)
		return
	}
	password := RandStringRunes(12)
	if role == "teacher" {
		result, err = AddTeacherSQL.Exec(name, password)
	} else {
		result, err = AddStudentSQL.Exec(name, password)
	}
	if err != nil {
		log.Fatal(err)
	}
	id, err = result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	if role == "teacher" {
		init_teacher(int(id), password)
	} else {
		init_student(int(id), password)
	}
	fmt.Printf("|%s| is added. Must complete registeration.\n", name)
}

//-----------------------------------------------------------------
func complete_registrationHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	role := r.FormValue("role")
	course_id := r.FormValue("course_id")
	if course_id != Config.CourseId {
		fmt.Fprintf(w, "Failed")
		return
	}
	var err error
	var rows *sql.Rows
	if role == "teacher" {
		rows, err = Database.Query("select id, password from teacher where name=?", name)
		defer rows.Close()
	} else if role == "student" {
		rows, err = Database.Query("select id, password from student where name=?", name)
		defer rows.Close()
	} else {
		fmt.Fprintf(w, "Failed")
		return
	}
	if err != nil {
		fmt.Fprintf(w, "Failed")
		// log.Fatal(err)
	}
	var password string
	var id int
	for rows.Next() {
		rows.Scan(&id, &password)
		msg := fmt.Sprintf("%d,%s", id, password)
		// msg := ""
		// if Config.NameServer != "" {
		// 	msg = fmt.Sprintf("%d,%s,%s,%s", id, password, Config.CourseId, Config.NameServer)
		// } else {
		// 	msg = fmt.Sprintf("%d,%s,%s", id, password, Config.CourseId)
		// }
		fmt.Fprintf(w, msg)
		return
	}
	fmt.Fprintf(w, "Failed")
}

//-----------------------------------------------------------------
