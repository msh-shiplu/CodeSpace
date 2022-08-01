package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func teacherSigninCompleteHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("username")
	password := r.FormValue("password")
	expectedPass, ok := TeacherPass[name]
	if !ok || expectedPass != password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(3 * time.Hour)

	sessions[sessionToken] = session{
		username: name,
		expiry:   expiresAt,
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	})
	fmt.Fprintf(w, "%d", TeacherNameToId[name])
	// http.Redirect(w, r, "view_exercises?role=teacher", http.StatusFound)
}

func teacherSigninHandler(w http.ResponseWriter, r *http.Request) {
	temp := template.New("")
	t, err := temp.Parse(TEACHER_LOGIN)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
}
