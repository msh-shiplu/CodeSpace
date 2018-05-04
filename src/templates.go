//
// Author: Vinhthuy Phan, 2018
//
package main

var STUDENT_MESSAGING_TEMPLATE = `
<html>
	<head>
  		<title>Student messaging</title>
		<meta http-equiv="refresh" content="10" />
	</head>
	<style>
		.bottom {
			position: fixed;
			bottom: 0;
			font-size: 150%;
			color: red;
		}
	</style>
	<body>
	<div class="bottom">{{.Message}}</div>
	</body>
</html>
`

var TEACHER_MESSAGING_TEMPLATE = `
<html>
	<head>
  		<title>Teacher messaging</title>
		<script src="https://cdn.rawgit.com/google/code-prettify/master/loader/run_prettify.js?autoload=true&skin=sons-of-obsidian"></script>
  		<script src="http://code.jquery.com/jquery-3.1.1.min.js"></script>
	    <script type="text/javascript">
			var updateInterval = 5000;		// 5 sec update interval
			var maxUpdateTime =  1800000;   // no longer update after 30 min.
			var totalUpdateTime = 0;
			function getData() {
				var url = "http://{{.Address}}/bulletin_board_data";
				$.getJSON(url, function( data ) {
					console.log(data);
					$("#p1").html(data["P1"]);
					$("#p2").html(data["P2"]);
					$("#ap").html(data["ActiveProblems"]);
					$("#bu").html(data["BulletinItems"]);
					$("#at").html(data["Attendance"]);
				});
			}
			$(document).ready(function(){
				getData();
				handle = setInterval(getData, updateInterval);
			});
	    </script>
	</head>
	<style>
		.bottom {
			position: fixed;
			bottom: 0;
			text-align: center;
			width: 100%;
		}
		.label{ display: inline; }
		#p1, #p2, #ap, #bu, #at {
			padding: 0.75em;
			display: inline;
		}
		#p1 { color: green; }
		#p2 { color: red; }
		pre {
			font-family: monospace;
			font-size:120%;
			margin-top:50px;
			padding-left:2em;
			overflow-x:scroll;
			overflow-y:scroll;
			tab-size: 4;
			-moz-tab-size: 4;
		}
		.center {
		    text-align: center;
		}
		.pagination {
		    display: inline-block;
		    padding-bottom: 20px;
		}
		.pagination a {
		    color: black;
		    float: left;
		    padding: 8px 16px;
		    text-decoration: none;
		    transition: background-color .3s;
		    border: 1px solid #ddd;
		    margin: 0 4px;
		    border-radius: 5px;
		}
		.pagination a.active {
		    background-color: #4CAF50;
		    color: white;
		    border: 1px solid #4CAF50;
		    border-radius: 5px;
		}
		.pagination a:hover:not(.active) {background-color: #ddd;}
		.nav a { text-decoration: none; padding:3px;}
		.nav { display: inline-block; vertical-align: baseline;}
		#navWrap{position:absolute;top:20;right:10;}
	</style>
	<body>
	<div id="navWrap">
	{{ if .Authenticated }}
	<div class="nav"><a href="view_bulletin_board?i=0&pc={{.PC}}">First<a></div>
	<div class="nav"><a href="view_bulletin_board?i={{.PrevI}}&pc={{.PC}}">Prev<a></div>
	<div class="nav"><a href="view_bulletin_board?i={{.NextI}}&pc={{.PC}}">Next<a></div>
	<div class="nav"><a href="remove_bulletin_page?i={{.I}}&pc={{.PC}}">&#x2718;</a></div>
	{{ end }}
	</div>
	<pre class="prettyprint linenums">{{.Code}}</pre>

	<div class="bottom">
	<div class="label">&#128546;</div><div id="p2">{{.P2}}</div>
	<div class="label">&#128526;</div><div id="p1">{{.P1}}</div>
	<div class="label">Problems:</div><div id="ap">{{.ActiveProblems}}</div>
	<div class="label">Bulletin:</div><div id="bu">{{.BulletinItems}}</div>
	<div class="label">Attendance:</div><div id="at">{{.Attendance}}</div>
	</div>
	</body>
</html>
`

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

var TAG_REPORT_TEMPLATE = `
<html>
  <head>
    <!--Load the AJAX API-->
    <script type="text/javascript" src="https://www.gstatic.com/charts/loader.js"></script>
    <script type="text/javascript">
      google.charts.load('current', {'packages':['bar','scatter']});
      google.charts.setOnLoadCallback(draw_by_pids);
      google.charts.setOnLoadCallback(draw_success);
      google.charts.setOnLoadCallback(draw_participation);

      function draw_by_pids() {
	      var data = google.visualization.arrayToDataTable([
	        ['Time', 'Success', 'Participation'],
			{{ range $pid, $rec := .Performance }}
				[ '{{$pid}}', {{$rec.Success}}, {{$rec.Activity}} ],
			{{ end }}
	      ]);
	      var options = {
	        title: {{.Description}},
        	height: 300,
	        vAxis: { 
	        	textStyle: { fontSize: 20},
	            viewWindowMode:'explicit',
    	        viewWindow: { min:0, max:1.05 }
	        },
            hAxis: { title: 'Problem ID', textStyle: { fontSize: 20} },
        	fontSize: 24,
	      };
        var chart = new google.charts.Bar(document.getElementById('bypids'));
        chart.draw(data, google.charts.Bar.convertOptions(options));
      }

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
	        	textStyle: { fontSize: 20},
	            viewWindowMode:'explicit',
    	        viewWindow: { min:0, max:1.05 }
	        },
            hAxis: { title: '', textStyle: { fontSize: 20} },
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
	        	textStyle: { fontSize: 20},
	            viewWindowMode:'explicit',
    	        viewWindow: { min:0, max:1.05 }
	        },
            hAxis: { title: '', textStyle: { fontSize: 20} },
        	fontSize: 24,
	      };
        var chart = new google.charts.Scatter(document.getElementById('participation'));
        chart.draw(data, google.charts.Scatter.convertOptions(options));
      }
      </script>
    <style>
    #bypids{ margin: auto; width:75%; }
    #row{
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
    	width:70%;
    	height:40px;
		overflow-x: scroll;
	    white-space: nowrap;
    	border-bottom:10px solid #F0F0F0;
    	text-align: center;
    }
    .problem_id{
    	padding-left:10px;
    	padding-right:10px;
    	margin:auto;
    	text-align: center;
    	display: inline;
    	font-size: 110%;
    }
    .problem_id a{
    	text-decoration: none;
    }
    </style>
  </head>

  <body>
    <div id="bypids"></div>
	<div class="spacer"></div>
    <div id="problem_ids">
	{{ range $pid, $rec := .Performance }}
		<div class="problem_id"><a href="analyze_submissions?pid={{$pid}}&pc={{$rec.PC}}" target="_blank">{{$pid}}</a></div>
	{{ end }}
	</div>
	<div class="spacer"></div>
	<div id="row">
	    <div id="success"></div>
	    <div id="participation"></div>
    </div>
  </body>
</html>
`

var TAGS_VIEW_TEMPLATE = `
<html>
  <head>
    <style>
    body { font-size: 16pt;}
    </style>
  </head>
  <body>
  	<h1>Learning objectives</h1>
  	<ul>
	{{ range . }}
	<li><a href="report_tag?pc={{.PC}}&tag_id={{.Id}}">{{.Description}}</a></li>
	{{ end }}
	</ul>
  </body>
</html>
`

var ACTIVITY_VIEW_TEMPLATE = `
<html>
  <head>
    <!--Load the AJAX API-->
    <script type="text/javascript" src="https://www.gstatic.com/charts/loader.js"></script>
    <script type="text/javascript">
      google.charts.load('current', {'packages':['bar','line']});
      google.charts.setOnLoadCallback(drawStudentCount);
      google.charts.setOnLoadCallback(drawExerciseCount);

      function drawStudentCount() {
	      var data = google.visualization.arrayToDataTable([
	        ['Time', 'Student Count'],
			{{ range $day, $val := . }}
				[  new Date({{$day}} / 1000000), {{$val.SidCount}} ],
			{{ end }}
	      ]);
	      var options = {
	        title: 'Daily student participation',
        	height: 350,
	        vAxis: { textStyle: { fontSize: 20} },
            hAxis: { title: '', textStyle: { fontSize: 20} },
        	fontSize: 24,
	      };
        var chart = new google.charts.Bar(document.getElementById('chart1_div'));
        chart.draw(data, google.charts.Bar.convertOptions(options));
      }

      function drawExerciseCount() {
	      var data = google.visualization.arrayToDataTable([
	        ['Time', 'Excercise Count'],
			{{ range $day, $val := . }}
				[  new Date({{$day}} / 1000000), {{$val.PidCount}} ],
			{{ end }}
	      ]);
	      var options = {
	        title: 'Daily exercise',
        	height: 350,
	        vAxis: { textStyle: { fontSize: 20} },
            hAxis: { title: '', textStyle: { fontSize: 20} },
        	fontSize: 24,
	      };
        var chart = new google.charts.Bar(document.getElementById('chart2_div'));
        chart.draw(data, google.charts.Bar.convertOptions(options));
      }
    </script>
    <style>
    #chart1_div,#chart2_div{ margin: auto; width:75%; }
    .spacer{ width:100%; height:40px; }
    </style>
  </head>

  <body>
    <div id="chart1_div"></div>
    <div class="spacer"></div>
    <div id="chart2_div"></div>
  </body>
</html>
`

var ANALYZE_SUBMISSIONS_TEMPLATE = `
<html>
  <head>
    <!--Load the AJAX API-->
    <script type="text/javascript" src="https://www.gstatic.com/charts/loader.js"></script>
    <script type="text/javascript">
      google.charts.load('current', {'packages':['timeline','corechart','bar']});
      google.charts.setOnLoadCallback(drawSubmissions);
      google.charts.setOnLoadCallback(drawWaitingTime);
      google.charts.setOnLoadCallback(drawResponseTime);
      google.charts.setOnLoadCallback(drawAttempts);


      function drawResponseTime() {
        var data = google.visualization.arrayToDataTable([
			['Student', 'Duration'],
			{{ range $sid, $rec := . }}
				{{ range $rec }}
				[  String({{$sid}}), ({{.At}} - {{.Start}})/1e9],
				{{ end }}
			{{ end }}
        ]);

        var options = {
          title: 'Student response time to a problem',
          legend: { position: 'none' },
        };

        var chart = new google.visualization.Histogram(document.getElementById('response'));
        chart.draw(data, options);
      }

      //---------------------------------------------------
      function drawWaitingTime() {
        var data = google.visualization.arrayToDataTable([
			['Student', 'Waiting time'],
			{{ range $sid, $rec := . }}
				{{ range $rec }}
				[  String({{$sid}}), ({{.Completed}} - {{.At}})/1e9],
				{{ end }}
			{{ end }}
        ]);

        var options = {
          title: 'Waiting time for response from teacher',
          legend: { position: 'none' },
        };

        var chart = new google.visualization.Histogram(document.getElementById('waiting'));
        chart.draw(data, options);
      }

      //---------------------------------------------------
      function drawAttempts() {
        var data = google.visualization.arrayToDataTable([
			['Student', 'Attempts'],
			{{ range $sid, $rec := . }}
				{{ $length := len $rec }}
				[  String({{$sid}}), {{$length}} ],
			{{ end }}
        ]);

        var options = {
			title: 'Solution attempts',
			legend: { position: 'none' },
	        hAxis: { 
	            viewWindowMode:'explicit',
		        viewWindow: { min:1 }
	        },
        };

        var chart = new google.visualization.Histogram(document.getElementById('attempts'));
        chart.draw(data, options);
      }


      //---------------------------------------------------
      function drawSubmissions() {
        var container = document.getElementById('timeline');
        var chart = new google.visualization.Timeline(container);
        var dataTable = new google.visualization.DataTable();

        dataTable.addColumn({ type: 'string', id: 'Student' });
        dataTable.addColumn({ type: 'string', id: 'Flag' });
        dataTable.addColumn({ type: 'date', id: 'Start' });
        dataTable.addColumn({ type: 'date', id: 'End' });
        dataTable.addRows([
			{{ range $sid, $rec := . }}
				{{ range $rec }}
				[  String({{$sid}}), {{.Flag}}, new Date({{.At}}/1000000), new Date({{.Completed}}/1000000)],
				{{ end }}
			{{ end }}
		]);			
	    var options = {
	    	title: 'Submissions',
	    };
        chart.draw(dataTable, options);
      }
    </script>
    <style>
    #waiting,#attempts,#response{ margin: auto; width:75%; height:300px; }
    #timeline{ margin: auto; width:75%; height:500px; }
    .spacer{ width:100%; height:40px; }
    </style>
  </head>
  <body>
    <div id="attempts"></div>
    <div class="spacer"></div>
    <div id="response"></div>
    <div class="spacer"></div>
    <div id="waiting"></div>
    <div class="spacer"></div>
    <div id="timeline"></div>
  </body>
</html>
`
