//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"net/http"
	"time"
)

//-----------------------------------------------------------------------------------
type ScoreEntry struct {
	Name     string
	Points   int
	Attempts int
	Count    int
}

type TagsViewData struct {
	Tags            map[int]string
	SubmissionCount map[string]int
	Scores          map[int]*ScoreEntry
	PC              string
}

//-----------------------------------------------------------------------------------
func reportHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("pc") != Passcode {
		fmt.Fprintf(w, "Unauthorized")
		return
	}
	rows, _ := Database.Query("select id, topic_description from tag")
	record := &TagsViewData{
		Tags:            make(map[int]string),
		SubmissionCount: make(map[string]int),
		Scores:          make(map[int]*ScoreEntry),
		PC:              Passcode,
	}

	var id int
	var des string
	for rows.Next() {
		rows.Scan(&id, &des)
		record.Tags[id] = des
	}
	rows.Close()

	rows, _ = Database.Query("select code_submitted_at from submission")
	var at time.Time
	for rows.Next() {
		rows.Scan(&at)
		date := fmt.Sprintf("%d.%d.%d", at.Month(), at.Day(), at.Year())
		record.SubmissionCount[date]++
	}
	rows.Close()

	rows, err := Database.Query("select score.score, score.graded_submission_number, score.student_id, student.name from score join student on score.student_id=student.id")
	if err != nil {
		fmt.Println(err)
		return
	}
	var points, attempts, student_id int
	var stname string
	for rows.Next() {
		rows.Scan(&points, &attempts, &student_id, &stname)
		if _, ok := record.Scores[student_id]; !ok {
			record.Scores[student_id] = &ScoreEntry{Name: stname}
		}
		record.Scores[student_id].Points += points
		record.Scores[student_id].Attempts += attempts
		record.Scores[student_id].Count += 1
	}
	rows.Close()

	w.Header().Set("Content-Type", "text/html")
	t, _ := template.New("").Parse(TAGS_VIEW_TEMPLATE)
	err = t.Execute(w, record)
	if err != nil {
		fmt.Println(err)
	}
}

var TAGS_VIEW_TEMPLATE = `
<html>
  <head>
    <!--Load the AJAX API-->
    <script type="text/javascript" src="https://www.gstatic.com/charts/loader.js"></script>
    <script type="text/javascript">
      google.charts.load('current', {'packages':['corechart', 'table']});
      google.charts.setOnLoadCallback(drawActivity);
      google.charts.setOnLoadCallback(drawScores);

      function drawActivity() {
	      var data = google.visualization.arrayToDataTable([
	        ['Date', 'Total submissions'],
			{{ range $day, $val := .SubmissionCount }}
				[  {{$day}}, {{$val}} ],
			{{ end }}
	      ]);
	      var options = {
	        title: '',
        	height: 350,
            hAxis: { title: 'Total submissions' },
	        legend: { position: 'none' },
	      };
        var chart = new google.visualization.SteppedAreaChart(document.getElementById('chart_div'));
        chart.draw(data, options);
      }

      function drawScores() {
        var data = new google.visualization.DataTable();
        data.addColumn('string', 'Name');
        data.addColumn('number', 'Points');
        data.addColumn('number', 'Attempts');
        data.addRows([
			{{ range $student_id, $entry := .Scores }}
				[{{$entry.Name}},{{$entry.Points}},{{$entry.Attempts}}],
			{{ end }}
        ]);
        var table = new google.visualization.Table(document.getElementById('scores_div'));
        table.draw(data, {showRowNumber: true, width: '400px'});
      }
    </script>
    <style>
    body { margin:auto; width:90%; font-size:16pt;}
    #chart_div,#scores_div{ margin:auto; }
    .spacer{ width:100%; height:30px; }
    </style>
  </head>
  <body>
  	<div class="spacer"></div>
  	<h4>Activities</h4>
  	<div id="chart_div"></div>
  	<div class="spacer"></div>
  	<h4>Learning objectives</h4>
  	<ul>
	{{$pc := .PC}}
	{{ range $tag_id, $tag_des := .Tags }}
	<li><a href="report_tag?pc={{$pc}}&tag_id={{$tag_id}}" target="_blank">{{$tag_des}}</a></li>
	{{ end }}
	</ul>
	<h4>Points</h4>
	<div id="scores_div"></div>
  	<div class="spacer"></div>
  </body>
</html>
`

//-----------------------------------------------------------------------------------
type ProblemPerformance struct {
	Pid       int
	Timestamp int64
	Correct   int
	Incorrect int
	Activity  float32
	Success   float32
	PC        string
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
	row, _ := Database.Query("select topic_description from tag where id=? limit 1", tag_id)
	tag_description := ""
	for row.Next() {
		row.Scan(&tag_description)
	}
	row.Close()

	query := "select problem.id, problem.merit, problem.at, score.points, score.student_id from problem join score on problem.id=score.problem_id join student where problem.tag=?"
	rows, err := Database.Query(query, tag_id)
	if err != nil {
		fmt.Println(err)
		return
	}
	var pid, merit, points, student_id int
	var at time.Time
	record := make(map[int]*ProblemPerformance)
	for rows.Next() {
		rows.Scan(&pid, &merit, &at, &points, &student_id)
		if _, ok := record[pid]; !ok {
			record[pid] = &ProblemPerformance{
				Pid:       pid,
				Timestamp: at.UnixNano(),
				Correct:   0,
				Incorrect: 0,
				Activity:  0,
				PC:        Passcode,
			}
		}
		if merit == points {
			record[pid].Correct++
		} else {
			record[pid].Incorrect++
		}
		record[pid].Activity += 1.0
	}
	rows.Close()

	var student_count float32
	rows, err = Database.Query("select count(*) from student")
	for rows.Next() {
		rows.Scan(&student_count)
	}
	rows.Close()

	for pid, _ := range record {
		record[pid].Success = float32(record[pid].Correct) / float32(record[pid].Correct+record[pid].Incorrect)
		record[pid].Activity = record[pid].Activity / student_count
	}

	w.Header().Set("Content-Type", "text/html")
	t, _ := template.New("").Parse(TAG_REPORT_TEMPLATE)
	err = t.Execute(w, &TagData{Description: tag_description, Performance: record})
	if err != nil {
		fmt.Println(err)
	}
}

//-----------------------------------------------------------------------------------

var TAG_REPORT_TEMPLATE = `
<html>
  <head>
    <!--Load the AJAX API-->
    <script type="text/javascript" src="https://www.gstatic.com/charts/loader.js"></script>
    <script type="text/javascript">
      google.charts.load('current', {'packages':['bar','scatter']});
      google.charts.setOnLoadCallback(draw_success);
      google.charts.setOnLoadCallback(draw_participation);

      function draw_success() {
	      var data = google.visualization.arrayToDataTable([
	        ['Time', 'Success'],
			{{ range $pid, $rec := .Performance }}
				[ new Date({{$rec.Timestamp}} / 1000000), {{$rec.Success}} ],
			{{ end }}
	      ]);
	      var options = {
	        title: 'Success', legend: {position: 'none'},
	        vAxis: {
	        	textStyle: { fontSize: 18},
	            viewWindowMode:'explicit',
    	        viewWindow: { min:0, max:1.05 }
	        },
            hAxis: { title: '', textStyle: {fontSize: 18} },
        	fontSize: 24,
	      };
        var chart = new google.charts.Scatter(document.getElementById('success'));
        chart.draw(data, google.charts.Scatter.convertOptions(options));
      }

      function draw_participation() {
	      var data = google.visualization.arrayToDataTable([
	        ['Time', 'Participation'],
			{{ range $pid, $rec := .Performance }}
				[ new Date({{$rec.Timestamp}} / 1000000), {{$rec.Activity}} ],
			{{ end }}
	      ]);
	      var options = {
	        title: 'Participation', legend: {position: 'none'},
	        vAxis: {
	        	textStyle: {fontSize: 18},
	            viewWindowMode:'explicit',
    	        viewWindow: { min:0, max:1.05 }
	        },
            hAxis: { title: '', textStyle: {fontSize: 18} },
        	fontSize: 24,
	      };
        var chart = new google.charts.Scatter(document.getElementById('participation'));
        chart.draw(data, google.charts.Scatter.convertOptions(options));
      }
      </script>
    <style>
    body{ margin: auto; width:75%; }
    .row{
	    display:flex;
	    flex-direction:row;
	    justify-content: space-around;
    }
    #success,#participation{
	    width:450px; height:400px;
	    display:flex;
	    flex-direction:column;
    }
    .spacer{ width:100%; height:40px; }
    #problem_ids{
    	margin:auto;
    	height:100px;
    	padding-top:20px;
		overflow-x: scroll;
	    white-space: nowrap;
    	text-align: center;
		vertical-align: middle;
    }
    .problem_id{
		padding: 15px 20px 15px 20px;
    	text-align: center;
    	font-size: 110%;
    	border:2px solid #dedede;
		display: inline-block;
    }
    .problem_id a{
    	text-decoration: none;
    }
    </style>
  </head>

  <body>
  	<div class="spacer"><h2>{{.Description}}</h2></div>
	<div class="row">
	    <div id="success"></div>
	    <div id="participation"></div>
    </div>
	<div class="spacer"></div>
    <div id="problem_ids">
	{{ range $pid, $rec := .Performance }}
		<div class="problem_id"><a href="analyze_submissions?pid={{$pid}}&pc={{$rec.PC}}" target="_blank">{{$pid}}</a></div>
	{{ end }}
	</div>
  </body>
</html>
`
