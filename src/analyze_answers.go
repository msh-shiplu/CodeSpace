//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"net/http"
	// "strconv"
)

type AnswersBoardMessage struct {
	Counts  map[string]int
	Content string
	Total   int
}

//-----------------------------------------------------------------------------------
func view_answersHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.FormValue("filename")
	passcode := r.FormValue("pc")
	if prob, ok := ActiveProblems[filename]; ok && passcode == Passcode {
		t, err := template.New("").Parse(VIEW_ANSWERS_TEMPLATE)
		if err == nil {
			answers := prob.Answers
			counts := make(map[string]int)
			total := 0
			for i := 0; i < len(answers); i++ {
				counts[answers[i]]++
				total++
			}
			content := prob.Info.Description
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

var VIEW_ANSWERS_TEMPLATE = `
<html>
  <head>
    <!--Load the AJAX API-->
    <script type="text/javascript" src="https://www.gstatic.com/charts/loader.js"></script>
    <script type="text/javascript">
      google.charts.load('current', {'packages':['bar']});
      google.charts.setOnLoadCallback(drawChart);
      function drawChart() {
        var data = google.visualization.arrayToDataTable([
          ["Answer", "Count", {role: 'annotation'}],
          	{{$total := .Total}}
			{{ range $key, $value := .Counts }}
				[{{ $key }}, {{ $value }}, Math.round(100 * {{$value}} / {{ $total }})  + '%'],
			{{ end }}
		]);
        var options = {
        	'title':'Total votes: {{ .Total }}',
        	'height':300,
        	'legend': {position: "none"},
        	'fontSize': 24,
            vAxis: { textStyle: { fontSize: 20} },
            hAxis: { title: "", textStyle: { fontSize: 20} },
        };
        var chart = new google.charts.Bar(document.getElementById('chart_div'));
        chart.draw(data, google.charts.Bar.convertOptions(options));
      }
    </script>
    <style>
    #chart_div{ margin: auto; width:70%; }
    pre{ margin: auto; width:60%}
    .spacer{ width:100%; height:40px; }
    </style>
  </head>

  <body>
    <div id="chart_div"></div>
    <div class="spacer"></div>
    <pre id="content">{{ .Content }}</pre>
    <div class="spacer"></div>
  </body>
</html>
`
