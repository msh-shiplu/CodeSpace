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
					$("#p1u").html(data["P1Ungraded"]);
					$("#p1g").html(data["P1Graded"]);
					$("#p2u").html(data["P2Unanswered"]);
					$("#p2a").html(data["P2Answered"]);
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
		.bottom_left {
			position: fixed;
			bottom: 0;
			text-align: left,
			width: 50%;
		}
		.bottom_right {
			position: fixed;
			bottom: 0;
			right: 50px;
			text-align: right,
			width: 50%;
		}
		.label{ display: inline; }
		#p1, #p2, #p1g, #p1u, #p2a, #p2u, #ans, #ap, #bu, #at {
			padding: 0.75em;
			display: inline;
		}
		#p1g, #p2a { color: green; }
		#p1u, #p2u { color: red; }
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

	<div class="bottom_left">
		<div id="p2">{{.P2}}</div> <div class="label">Help Requests:</div>
		<div id="p2u">{{.P2Unanswered}}</div> <div class="label">Pending,</div>
		<div id="p2a">{{.P2Answered}}</div> <div class="label">Answered</div>
	</div>
	<div class="bottom_right">
		<div id="p1">{{.P1}}</div> <div class="label">Submissions:</div>
		<div id="p1u">{{.P1Ungraded}}</div> <div class="label">Pending,</div>
		<div id="p1g">{{.P1Graded}}</div> <div class="label">Graded</div>
	</div>
	</body>
</html>
`
var CODESPACE_TEMPLATE = `
	<!DOCTYPE html>
	<html>
	<head>
	<title>CodeSpace</title>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bulma/0.9.3/css/bulma.min.css" integrity="sha512-IgmDkwzs96t4SrChW29No3NXBIBv8baW490zk5aXvhCD8vuZM3yUSkbyTBcXohkySecyzIrUwiF/qV0cuPcL3Q==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	</head>
	<body>
	<div class="container">
		{{if .Authenticated }}
			<table class="table is-striped is-fullwidth is-hoverable">
				<thead>
					<tr>
						<th>Student</th>
						<th>Problem</th>
						<th>Last Snapshot Since</th>
						<th>Time Spent</th>
						<th>Number of Lines</th>
						<th>Status</th>
						<th></th>
					</tr>
				</thead>
				<tbody>
				{{ range .Snapshots }}
				<tr>
					<td>{{ .StudentName }}</td>
					<td>{{ .ProblemName }}</td>
					<td>{{ formatTimeSince .LastUpdated }}</td>
					<td>{{ formatTimeSince .FirstUpdate }}</td>
					<td>{{ .LinesOfCode }}</td>
					<td>{{ .Status }}</td>
					<td><a href="/get_snapshot?student_id={{ .StudentID }}&problem_id={{ .ProblemID }}&uid={{$.UserID}}&role={{$.UserRole}}&pc={{$.Passcode}}">View</a></td>
				</tr>
				{{ end }}
				</tbody>
			</table>
		{{else}}
			<h1 class="title">Unauthorized access!!!</h1>
		{{end}}
	</div>

	</body>
	</html>
`
var CODE_SNAPSHOT_TEMPLATE = `
<!DOCTYPE html>
	<html>
	<head>
	<title>CodeSpace</title>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/codemirror.min.js" integrity="sha512-hGVnilhYD74EGnPbzyvje74/Urjrg5LSNGx0ARG1Ucqyiaz+lFvtsXk/1jCwT9/giXP0qoXSlVDjxNxjLvmqAw==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/mode/python/python.min.js" integrity="sha512-/mavDpedrvPG/0Grj2Ughxte/fsm42ZmZWWpHz1jCbzd5ECv8CB7PomGtw0NAnhHmE/lkDFkRMupjoohbKNA1Q==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/mode/clike/clike.min.js" integrity="sha512-GAled7oA9WlRkBaUQlUEgxm37hf43V2KEMaEiWlvBO/ueP2BLvBLKN5tIJu4VZOTwo6Z4XvrojYngoN9dJw2ug==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/codemirror.min.css" integrity="sha512-6sALqOPMrNSc+1p5xOhPwGIzs6kIlST+9oGWlI4Wwcbj1saaX9J3uzO3Vub016dmHV7hM+bMi/rfXLiF5DNIZg==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/theme/monokai.min.css" integrity="sha512-R6PH4vSzF2Yxjdvb2p2FA06yWul+U0PDDav4b/od/oXf9Iw37zl10plvwOXelrjV2Ai7Eo3vyHeyFUjhXdBCVQ==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bulma/0.9.3/css/bulma.min.css" integrity="sha512-IgmDkwzs96t4SrChW29No3NXBIBv8baW490zk5aXvhCD8vuZM3yUSkbyTBcXohkySecyzIrUwiF/qV0cuPcL3Q==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	</head>
	<body>
		<div class="container">
			<section class="section">
				<h2 class="title is-2">Code Snapshot at {{.Snapshot.LastUpdated.Format "Jan 02, 2006 3:4:5 PM"}}</h2>
				<h3 class="title is-3">Student: {{.Snapshot.StudentName}}, Problem: {{.Snapshot.ProblemName}}</h3>
				<h3>
				<textarea id="editor">{{ .Snapshot.Code }}</textarea>
				<form action="/save_snapshot_feedback" method="POST">
					<textarea class="textarea" placeholder="Write your feedback!" name="feedback"></textarea>
					<input class="button" type="submit" value="Send Feedback">
					
					<input type="hidden" name="snapshot_id" value="{{.Snapshot.ID}}">
					<input type="hidden" name="uid" value="{{.UserID}}">
					<input type="hidden" name="role" value="{{.UserRole}}">
				</form>
			</section>
			<section class="section">
				{{range .Feedbacks}}
					<article class="message">
						<div class="message-header">
						<p>Feedback given at {{.FeedbackTime}}</p>
						</div>
						<div class="message-body">
							{{.Feedback}}
						</div>
					</article>
				{{end}}
			</section>
		</div>
		<script>
			var editor = document.getElementById("editor");
			var myCodeMirror = CodeMirror.fromTextArea(editor, {lineNumbers: true, mode: "{{getEditorMode .Snapshot.ProblemName}}", theme: "monokai", matchBrackets: true, indentUnit: 4, indentWithTabs: true, readOnly: "nocursor"});
			// myCodeMirror1.setSize("80%", 900)
		</script>
	</body>
	</html>
`
