//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//-----------------------------------------------------------------------------------
func student_getsHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	var js []byte
	var err error

	BoardsSem.Lock()
	defer BoardsSem.Unlock()

	if _, ok := Students[uid]; ok {
		js, err = json.Marshal(Students[uid].Boards)
		Students[uid].Boards = []*Board{}
		if err == nil {
			// fmt.Println(string(js))
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			return
		}
		for _, b := range Students[uid].Boards {
			addOrUpdateStudentStatus(uid, b.Pid, "Working", "", "", "")
		}
	}
	fmt.Println(err.Error())
	js, err = json.Marshal([]*Board{})
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

//-----------------------------------------------------------------------------------
