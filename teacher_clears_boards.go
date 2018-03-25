//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"fmt"
	"net/http"
)

//-----------------------------------------------------------------------------------
func teacher_clears_boardsHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	for stid, _ := range Boards {
		Boards[stid] = make([]*Board, 0)
	}
	for stid, _ := range MessageBoards {
		MessageBoards[stid] = "White board is empty."
	}
	fmt.Fprintf(w, "Whiteboards cleared.")
}
