//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"fmt"
	"net/http"
)

//-----------------------------------------------------------------------------------
func teacher_gets_passcodeHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
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
			fmt.Printf("Attempts: %d, Ext: %s, Pid: %d, Answer: %s, len of content: %d\n",
				b.Attempts, b.Ext, b.Pid, b.Answer, len(b.Content))
		}
	}

	fmt.Println("NextProblem")
	for pid, p := range NextProblem {
		fmt.Println(pid, p)
	}

	fmt.Println("WorkingSubs")
	for i := 0; i < len(WorkingSubs); i++ {
		fmt.Println(WorkingSubs[i].Sid, WorkingSubs[i].Uid, WorkingSubs[i].Pid, WorkingSubs[i].Priority)
		fmt.Println(WorkingSubs[i].Content)
	}

	fmt.Println("MessageBoards")
	for uid, mesg := range MessageBoards {
		fmt.Println(uid, mesg)
	}

	fmt.Fprintf(w, Passcode)
}

//-----------------------------------------------------------------------------------
