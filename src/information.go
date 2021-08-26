package main

import (
	"encoding/json"
	"net/http"
)

type ActiveProblemInfo struct {
	ProblemID int
	Filename  string
}
type globalInfo struct {
	ActiveProblems []*ActiveProblemInfo
}

func globalInfoHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {

	activeProblems := make([]*ActiveProblemInfo, 0)
	for _, problem := range ActiveProblems {
		if problem.Active == true {
			activeProblems = append(activeProblems, &ActiveProblemInfo{problem.Info.Pid, problem.Info.Filename})
		}
	}
	g := &globalInfo{
		ActiveProblems: activeProblems,
	}
	js, _ := json.Marshal(g)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
