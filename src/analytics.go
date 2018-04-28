//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	// "encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

//-----------------------------------------------------------------------------------
type TagsData struct {
	Id          int
	Description string
	PC          string
}

//-----------------------------------------------------------------------------------
func view_tagsHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("pc") != Passcode {
		fmt.Fprintf(w, "Unauthorized")
		return
	}
	rows, err := Database.Query("select id, description from tag")
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
	} else {
		tags := make([]*TagsData, 0)
		var id int
		var des string
		for rows.Next() {
			rows.Scan(&id, &des)
			tags = append(tags, &TagsData{Id: id, Description: des, PC: Passcode})
		}
		w.Header().Set("Content-Type", "text/html")
		t, _ := template.New("").Parse(TAGS_VIEW_TEMPLATE)
		err = t.Execute(w, tags)
		if err != nil {
			fmt.Println(err)
		}
	}
}

//-----------------------------------------------------------------------------------
type ProblemPerformance struct {
	Pid       int
	Timestamp int64
	Activity  int
	Correct   int
	Incorrect int
	Success   float32
}

type TagData struct {
	Description string
	Performance map[int]*ProblemPerformance
}

//-----------------------------------------------------------------------------------
func report_tagHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("pc") != Passcode {
		fmt.Fprintf(w, "Unauthorized")
		return
	}
	tag_id := r.FormValue("tag_id")
	row, _ := Database.Query("select description from tag where id=? limit 1", tag_id)
	tag_description := ""
	for row.Next() {
		row.Scan(&tag_description)
	}
	row.Close()

	query := "select problem.id, problem.merit, problem.at, score.points, score.stid from problem join score on problem.id == score.pid where problem.tag=?"
	rows, err := Database.Query(query, tag_id)
	if err != nil {
		fmt.Println(err)
		return
	}
	var pid, merit, points, stid int
	var at time.Time
	record := make(map[int]*ProblemPerformance)
	for rows.Next() {
		rows.Scan(&pid, &merit, &at, &points, &stid)
		if _, ok := record[pid]; !ok {
			record[pid] = &ProblemPerformance{
				Pid:       pid,
				Timestamp: at.UnixNano(),
				// Timestamp: at.Format("2016-01-02 15:04"),
				Correct:   0,
				Incorrect: 0,
				Activity:  0,
			}
		}
		if merit == points {
			record[pid].Correct++
		} else {
			record[pid].Incorrect++
		}
		record[pid].Activity++
	}
	rows.Close()
	for pid, _ := range record {
		record[pid].Success = float32(record[pid].Correct) / float32(record[pid].Correct+record[pid].Incorrect)
	}
	w.Header().Set("Content-Type", "text/html")
	t, _ := template.New("").Parse(TAG_REPORT_TEMPLATE)
	err = t.Execute(w, &TagData{Description: tag_description, Performance: record})
	if err != nil {
		fmt.Println(err)
	}
}

//-----------------------------------------------------------------------------------
type BulletinBoardMessage struct {
	Code           string
	I              int
	NextI          int
	PrevI          int
	PC             string
	P1             int
	P2             int
	ActiveProblems int
	BulletinItems  int
	Attendance     int
	Address        string
	Authenticated  bool
}

type AnswersBoardMessage struct {
	Counts  map[string]int
	Content string
	Total   int
}

//-----------------------------------------------------------------------------------
func view_answersHandler(w http.ResponseWriter, r *http.Request) {
	pid, err := strconv.Atoi(r.FormValue("pid"))
	passcode := r.FormValue("pc")
	if _, ok := ActiveProblems[pid]; err == nil && ok && passcode == Passcode {
		t, err := template.New("").Parse(VIEW_ANSWERS_TEMPLATE)
		if err == nil {
			answers := ActiveProblems[pid].Answers
			counts := make(map[string]int)
			total := 0
			for i := 0; i < len(answers); i++ {
				counts[answers[i]]++
				total++
			}
			content := ActiveProblems[pid].Info.Description
			w.Header().Set("Content-Type", "text/html")
			data := &AnswersBoardMessage{Counts: counts, Content: content, Total: total}
			err = t.Execute(w, data)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println(err)
		}
	}
}

//-----------------------------------------------------------------------------------
