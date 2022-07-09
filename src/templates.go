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
	<meta http-equiv="refresh" content="120" >
	<script src="https://kit.fontawesome.com/923539b4ee.js" crossorigin="anonymous"></script>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bulma/0.9.3/css/bulma.min.css" integrity="sha512-IgmDkwzs96t4SrChW29No3NXBIBv8baW490zk5aXvhCD8vuZM3yUSkbyTBcXohkySecyzIrUwiF/qV0cuPcL3Q==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	</head>
	<body>
	<div class="container">
		<h3 class="title is-3">CodeSpace: List of Code Snapshots</h3>
		<table class="table is-striped is-fullwidth is-hoverable is-narrow">
			<thead>
				<tr>
					<th>Student</th>
					<th>Last Snapshot</th>
					<th>Time Spent</th>
					<th>Lines of Code</th>
					<th>Number of Feedback Messages</th>
					<th>Status</th>
					<th></th>
				</tr>
			</thead>
			<tbody>
			{{ range .Snapshots }}
			<tr>
				<td>{{ .StudentName }}</td>
				<td>{{ formatTimeSince .LastUpdated }} ago</td>
				<td>{{ formatTimeSince .FirstUpdate }}</td>
				<td>{{ .LinesOfCode }}</td>
				<td>{{ .NumFeedback }}</td>
				<td>{{ if eq .Status 0 }} Not Submitted {{else if eq .Status 1}} Submitted {{else if eq .Status 2}} <span style="font-size: 1.5em; color: red;"> <i class="far fa-times-circle"></i> </span> {{else if eq .Status 3}} <span style="font-size: 1.5em; color: green;"> <i class="far fa-check-circle"></i> </span> {{end}}</td>
				<td><a href="/get_snapshot?student_id={{ .StudentID }}&problem_id={{ .ProblemID }}&uid={{$.UserID}}&role={{$.UserRole}}&password={{$.Password}}">View</a></td>
			</tr>
			{{ end }}
			</tbody>
		</table>
	</div>

	</body>
	</html>
`
var CODE_SNAPSHOT_TEMPLATE = `
<!DOCTYPE html>
	<html>
	<head>
	<title>Latest Code Snapshot from {{.Snapshot.StudentName}}</title>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/codemirror.min.js" integrity="sha512-hGVnilhYD74EGnPbzyvje74/Urjrg5LSNGx0ARG1Ucqyiaz+lFvtsXk/1jCwT9/giXP0qoXSlVDjxNxjLvmqAw==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/mode/python/python.min.js" integrity="sha512-/mavDpedrvPG/0Grj2Ughxte/fsm42ZmZWWpHz1jCbzd5ECv8CB7PomGtw0NAnhHmE/lkDFkRMupjoohbKNA1Q==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/mode/clike/clike.min.js" integrity="sha512-GAled7oA9WlRkBaUQlUEgxm37hf43V2KEMaEiWlvBO/ueP2BLvBLKN5tIJu4VZOTwo6Z4XvrojYngoN9dJw2ug==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/codemirror.min.css" integrity="sha512-6sALqOPMrNSc+1p5xOhPwGIzs6kIlST+9oGWlI4Wwcbj1saaX9J3uzO3Vub016dmHV7hM+bMi/rfXLiF5DNIZg==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/theme/monokai.min.css" integrity="sha512-R6PH4vSzF2Yxjdvb2p2FA06yWul+U0PDDav4b/od/oXf9Iw37zl10plvwOXelrjV2Ai7Eo3vyHeyFUjhXdBCVQ==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bulma/0.9.3/css/bulma.min.css" integrity="sha512-IgmDkwzs96t4SrChW29No3NXBIBv8baW490zk5aXvhCD8vuZM3yUSkbyTBcXohkySecyzIrUwiF/qV0cuPcL3Q==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<script src="https://kit.fontawesome.com/923539b4ee.js" crossorigin="anonymous"></script>
	<script src="https://code.jquery.com/jquery-3.6.0.min.js" integrity="sha256-/xUj+3OJU5yExlq6GSYGSHk7tPXikynS7ogEvDej/m4=" crossorigin="anonymous"></script>
	<script src="https://code.jquery.com/ui/1.12.1/jquery-ui.min.js" integrity="sha256-VazP97ZCwtekAsvgPBSUwPFKdrwD3unUfSGVYrahUqU=" crossorigin="anonymous"></script>
	<link rel="stylesheet" href="https://code.jquery.com/ui/1.12.1/themes/base/jquery-ui.css" />
	</head>
	<body>
		<div class="container">
			<section class="section">
				<h3 class="title is-3">Latest Code Snapshot from {{.Snapshot.StudentName}}</h3>
				<h4 class="title is-4">{{.Snapshot.StudentName}} ({{.Snapshot.ProblemName}} @ {{.Snapshot.LastUpdated.Format "Jan 02, 2006 3:04:05 PM"}})</h4>
				<h5 class="title is-5">If you think that this student needs help, feel free to offer a brief comment.</h5>
				{{$l := (len .HelpRequestIDs)}}
				{{if ne $l 0}}
					<b>Help requests: </b>
				{{end}}
				{{range $i, $v := .HelpRequestIDs}}
					<a href="/view_help_request?request_id={{.}}&uid={{$.UserID}}&role={{$.UserRole}}&password={{$.Password}}">Request {{add $i 1}}</a>{{if lt (add $i 1) $l}} | {{end}}
				{{end}}
				<textarea id="editor">{{ .Snapshot.Code }}</textarea>
			</section>
			{{if lt .Snapshot.Status 3}}
			<section class="section" style="margin-top: 0px !important;">
				<form action="/save_snapshot_feedback" method="POST">
					<textarea class="textarea" placeholder="Write your feedback!" name="feedback"></textarea>
					<input class="button" type="submit" value="Send Feedback">
					
					<input type="hidden" name="snapshot_id" value="{{.Snapshot.ID}}">
					<input type="hidden" name="uid" value="{{.UserID}}">
					<input type="hidden" name="role" value="{{.UserRole}}">
					<input type="hidden" name="password" value="{{.Password}}">
				</form>
			</section>
			{{end}}
			<section class="section">
				{{range .Feedbacks}}
					<article class="message">
						<div class="message-header">
						<p>{{.GivenBy}} ({{.FeedbackTime.Format "Jan 02, 2006 3:4:5 PM"}})</p>
						</div>
						<div class="message-body">
							<div class="columns">
								<div class="column is-three-quarters">{{.Feedback}}</div>
								<div class="column">
									<a onclick="autoFeedbackSubmit('yes', {{.FeedbackID}})">
										<span style="font-size: 1.5em; {{if eq .CurrentUserVote "yes"}} color: green; {{end}}">
											<i class="fas fa-thumbs-up"></i>
										</span>
									</a>
									<span>
											{{.Upvote}}
									</span>
								</div>
								<div class="column">
									<a onclick="autoFeedbackSubmit('no', {{.FeedbackID}})">
										<span style="font-size: 1.5em; {{if eq .CurrentUserVote "no"}} color: red; {{end}}">
											<i class="fas fa-thumbs-down"></i>
										</span>
									</a>
									<span>
										{{.Downvote}}
									</span>
								</div>
							</div>
							<div class="codesnapshots">
								<h3>Code Snapshot</h3>
								<div>
									<textarea class="editors">{{ .Code }}</textarea>
								</div>
							</div>
						</div>
					</article>
				{{end}}
			</section>
		</div>
		<script>
			var editor = document.getElementById("editor");
			var myCodeMirror = CodeMirror.fromTextArea(editor, {lineNumbers: true, mode: "{{getEditorMode .Snapshot.ProblemName}}", theme: "monokai", matchBrackets: true, indentUnit: 4, indentWithTabs: true, readOnly: "nocursor"});
			myCodeMirror.setSize("100%", 400)
			var snapshotEditors = document.getElementsByClassName("editors");
			for (i = 0;i<snapshotEditors.length; i++) {
				CodeMirror.fromTextArea(snapshotEditors[i], {lineNumbers: true, mode: "{{getEditorMode .Snapshot.ProblemName}}", theme: "monokai", matchBrackets: true, indentUnit: 4, indentWithTabs: true, readOnly: "nocursor"});
			}
			$( function() {
				$( ".codesnapshots" ).accordion({
					collapsible: true,
					active: false
				});
			} );
			function autoFeedbackSubmit(backFeedback, fID) {
				$.ajax({
					url: "/save_snapshot_back_feedback",
					type: "POST",
					data:  {
						feedback: backFeedback,
						feedback_id: fID,
						uid: {{.UserID}},
						role: "{{.UserRole}}",
						password: "{{.Password}}",
					},
					success: function(data){
						console.log("Success!")
					}
				});
				
				location.reload();
			}
		</script>
	</body>
	</html>
`
var STUDENT_VIEWS_FEEDBACK_TEMPLATE = `
<!DOCTYPE html>
	<html>
	<head>
	<title>Review Feedback</title>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/codemirror.min.js" integrity="sha512-hGVnilhYD74EGnPbzyvje74/Urjrg5LSNGx0ARG1Ucqyiaz+lFvtsXk/1jCwT9/giXP0qoXSlVDjxNxjLvmqAw==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/mode/python/python.min.js" integrity="sha512-/mavDpedrvPG/0Grj2Ughxte/fsm42ZmZWWpHz1jCbzd5ECv8CB7PomGtw0NAnhHmE/lkDFkRMupjoohbKNA1Q==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/mode/clike/clike.min.js" integrity="sha512-GAled7oA9WlRkBaUQlUEgxm37hf43V2KEMaEiWlvBO/ueP2BLvBLKN5tIJu4VZOTwo6Z4XvrojYngoN9dJw2ug==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/codemirror.min.css" integrity="sha512-6sALqOPMrNSc+1p5xOhPwGIzs6kIlST+9oGWlI4Wwcbj1saaX9J3uzO3Vub016dmHV7hM+bMi/rfXLiF5DNIZg==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/theme/monokai.min.css" integrity="sha512-R6PH4vSzF2Yxjdvb2p2FA06yWul+U0PDDav4b/od/oXf9Iw37zl10plvwOXelrjV2Ai7Eo3vyHeyFUjhXdBCVQ==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bulma/0.9.3/css/bulma.min.css" integrity="sha512-IgmDkwzs96t4SrChW29No3NXBIBv8baW490zk5aXvhCD8vuZM3yUSkbyTBcXohkySecyzIrUwiF/qV0cuPcL3Q==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<script src="https://kit.fontawesome.com/923539b4ee.js" crossorigin="anonymous"></script>
	<script src="https://code.jquery.com/jquery-3.6.0.min.js" integrity="sha256-/xUj+3OJU5yExlq6GSYGSHk7tPXikynS7ogEvDej/m4=" crossorigin="anonymous"></script>
	<script src="https://code.jquery.com/ui/1.12.1/jquery-ui.min.js" integrity="sha256-VazP97ZCwtekAsvgPBSUwPFKdrwD3unUfSGVYrahUqU=" crossorigin="anonymous"></script>
	<link rel="stylesheet" href="https://code.jquery.com/ui/1.12.1/themes/base/jquery-ui.css" />
	</head>
	<body>
		<div class="container">
		<h1 class="title">Review Feedback for Problem: {{.Filename}}</h1>
			<div class="tabs is-centered is-boxed is-medium">
				<ul>
					<li {{if eq .ViewType "forme"}}class="is-active"{{end}}>
						<a href="student_views_feedback?pid={{.CurrentPid}}&viewtype=forme&role={{.UserRole}}&uid={{.UserID}}&password={{.Password}}">
						<span class="icon is-small"><i class="fas fa-address-book" aria-hidden="true"></i></span>
						<span>For me</span>
						</a>
					</li>
					<li {{if eq .ViewType "all"}}class="is-active"{{end}}>
						<a href="student_views_feedback?pid={{.CurrentPid}}&viewtype=all&role={{.UserRole}}&uid={{.UserID}}&password={{.Password}}">
						<span class="icon is-small"><i class="fas fa-list-ul" aria-hidden="true"></i></span>
						<span>All</span>
						</a>
					</li>
				</ul>
			</div>
			<section class="section">
				{{range .Feedbacks}}
					<article class="message">
						<div class="message-header">
						<p>{{.GivenBy}} gave feedback on {{$.Filename}} at ({{.FeedbackTime.Format "Jan 02, 2006 3:04:05 PM"}})</p>
						</div>
						<div class="message-body">
							<div class="columns">
								<div class="column is-three-quarters">{{.Feedback}}</div>
								<div class="column">
									<a onclick="autoFeedbackSubmit('yes', {{.FeedbackID}})">
										<span style="font-size: 1.5em; {{if eq .CurrentUserVote "yes"}} color: green; {{end}}">
											<i class="fas fa-thumbs-up"></i>
										</span>
									</a>
									<span>
											{{.Upvote}}
									</span>
								</div>
								<div class="column">
									<a onclick="autoFeedbackSubmit('no', {{.FeedbackID}})">
										<span style="font-size: 1.5em; {{if eq .CurrentUserVote "no"}} color: red; {{end}}">
											<i class="fas fa-thumbs-down"></i>
										</span>
									</a>
									<span>
										{{.Downvote}}
									</span>
								</div>
							</div>
							<div class="codesnapshots">
								<h3>Code Snapshot</h3>
								<div>
									<textarea class="editors">{{ .Code }}</textarea>
								</div>
							</div>
						</div>
					</article>
				{{end}}
			</section>
			<nav class="pagination is-rounded" role="navigation" aria-label="pagination">
			{{if not (eq .NextPid -1)}}
				<a class="pagination-next" href="student_views_feedback?pid={{.NextPid}}&viewtype={{.ViewType}}&role={{.UserRole}}&uid={{.UserID}}&password={{.Password}}">Next</a>
			{{end}}
				<ul class="pagination-list">
				</ul>
			</nav>
		</div>
		<script>
			var snapshotEditors = document.getElementsByClassName("editors");
			
			for (let i = 0; i<snapshotEditors.length; i++){
				CodeMirror.fromTextArea(snapshotEditors[i], {lineNumbers: true, mode: "{{getEditorMode $.Filename}}", theme: "monokai", matchBrackets: true, indentUnit: 4, indentWithTabs: true, readOnly: "nocursor"});
			}
			$( function() {
				$( ".codesnapshots" ).accordion({
					collapsible: true,
					active: false
				});
			} );
			function autoFeedbackSubmit(backFeedback, fID) {
				$.ajax({
					url: "/save_snapshot_back_feedback",
					type: "POST",
					data:  {
						feedback: backFeedback,
						feedback_id: fID,
						uid: {{.UserID}},
						role: "{{.UserRole}}",
						password: "{{.Password}}",
					},
					success: function(data){
						console.log("Success!")
					}
				});
				
				location.reload();
			}
		</script>
	</body>
	</html>
`
var TEACHER_VIEWS_FEEDBACK_TEMPLATE = `
<!DOCTYPE html>
	<html>
	<head>
	<title>Review Feedback</title>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/codemirror.min.js" integrity="sha512-hGVnilhYD74EGnPbzyvje74/Urjrg5LSNGx0ARG1Ucqyiaz+lFvtsXk/1jCwT9/giXP0qoXSlVDjxNxjLvmqAw==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/mode/python/python.min.js" integrity="sha512-/mavDpedrvPG/0Grj2Ughxte/fsm42ZmZWWpHz1jCbzd5ECv8CB7PomGtw0NAnhHmE/lkDFkRMupjoohbKNA1Q==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/mode/clike/clike.min.js" integrity="sha512-GAled7oA9WlRkBaUQlUEgxm37hf43V2KEMaEiWlvBO/ueP2BLvBLKN5tIJu4VZOTwo6Z4XvrojYngoN9dJw2ug==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/codemirror.min.css" integrity="sha512-6sALqOPMrNSc+1p5xOhPwGIzs6kIlST+9oGWlI4Wwcbj1saaX9J3uzO3Vub016dmHV7hM+bMi/rfXLiF5DNIZg==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/theme/monokai.min.css" integrity="sha512-R6PH4vSzF2Yxjdvb2p2FA06yWul+U0PDDav4b/od/oXf9Iw37zl10plvwOXelrjV2Ai7Eo3vyHeyFUjhXdBCVQ==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bulma/0.9.3/css/bulma.min.css" integrity="sha512-IgmDkwzs96t4SrChW29No3NXBIBv8baW490zk5aXvhCD8vuZM3yUSkbyTBcXohkySecyzIrUwiF/qV0cuPcL3Q==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<script src="https://kit.fontawesome.com/923539b4ee.js" crossorigin="anonymous"></script>
	<script src="https://code.jquery.com/jquery-3.6.0.min.js" integrity="sha256-/xUj+3OJU5yExlq6GSYGSHk7tPXikynS7ogEvDej/m4=" crossorigin="anonymous"></script>
	<script src="https://code.jquery.com/ui/1.12.1/jquery-ui.min.js" integrity="sha256-VazP97ZCwtekAsvgPBSUwPFKdrwD3unUfSGVYrahUqU=" crossorigin="anonymous"></script>
	<link rel="stylesheet" href="https://code.jquery.com/ui/1.12.1/themes/base/jquery-ui.css" />
	</head>
	<body>
		<div class="container">
		<<h1 class="title">Review Feedback for Problem: {{.Filename}}</h1>
			<section class="section">
				{{range .Feedbacks}}
					<article class="message">
						<div class="message-header">
						<p>{{.GivenBy}} gave feedback on {{$.Filename}} at ({{.FeedbackTime.Format "Jan 02, 2006 3:04:05 PM"}})</p>
						</div>
						<div class="message-body">
							<div class="columns">
								<div class="column is-three-quarters">{{.Feedback}}</div>
								<div class="column">
									<a onclick="autoFeedbackSubmit('yes', {{.FeedbackID}})">
										<span style="font-size: 1.5em; {{if eq .CurrentUserVote "yes"}} color: green; {{end}}">
											<i class="fas fa-thumbs-up"></i>
										</span>
									</a>
									<span>
											{{.Upvote}}
									</span>
								</div>
								<div class="column">
									<a onclick="autoFeedbackSubmit('no', {{.FeedbackID}})">
										<span style="font-size: 1.5em; {{if eq .CurrentUserVote "no"}} color: red; {{end}}">
											<i class="fas fa-thumbs-down"></i>
										</span>
									</a>
									<span>
										{{.Downvote}}
									</span>
								</div>
							</div>
							<div class="codesnapshots">
								<h3>Code Snapshot</h3>
								<div>
									<textarea class="editors">{{ .Code }}</textarea>
								</div>
							</div>
						</div>
					</article>
				{{end}}
			</section>
			<nav class="pagination is-rounded" role="navigation" aria-label="pagination">
			{{if not (eq .NextPid -1)}}
				<a class="pagination-next" href="teacher_views_feedback?pid={{.NextPid}}&role={{.UserRole}}&uid={{.UserID}}&password={{.Password}}">Next</a>
			{{end}}
				<ul class="pagination-list">
				</ul>
			</nav>
		</div>
		<script>
			var snapshotEditors = document.getElementsByClassName("editors");
			
			for (let i = 0; i<snapshotEditors.length; i++){
				CodeMirror.fromTextArea(snapshotEditors[i], {lineNumbers: true, mode: "{{getEditorMode .Filename}}", theme: "monokai", matchBrackets: true, indentUnit: 4, indentWithTabs: true, readOnly: "nocursor"});
			}
			$( function() {
				$( ".codesnapshots" ).accordion({
					collapsible: true,
					active: false
				});
			} );
			function autoFeedbackSubmit(backFeedback, fID) {
				$.ajax({
					url: "/save_snapshot_back_feedback",
					type: "POST",
					data:  {
						feedback: backFeedback,
						feedback_id: fID,
						uid: {{.UserID}},
						role: "{{.UserRole}}",
						password: "{{.Password}}",
					},
					success: function(data){
						console.log("Success!")
					}
				});
				
				location.reload();
			}
		</script>
	</body>
	</html>
`
var HELP_REQUEST_LIST_TEMPLATE = `
	<!DOCTYPE html>
	<html>
	<head>
	<title>Help Hotline</title>
	<meta http-equiv="refresh" content="120" >
	<script src="https://kit.fontawesome.com/923539b4ee.js" crossorigin="anonymous"></script>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bulma/0.9.3/css/bulma.min.css" integrity="sha512-IgmDkwzs96t4SrChW29No3NXBIBv8baW490zk5aXvhCD8vuZM3yUSkbyTBcXohkySecyzIrUwiF/qV0cuPcL3Q==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	</head>
	<body>
	<div class="container">
	<h3 class="title is-3">Help Hotline/<h3>
	<h5 class="title is-5">Currently, there are {{.NumHelpNeeded}} students who need help.</h5>
		<table class="table is-striped is-fullwidth is-hoverable is-narrow">
			<thead>
				<tr>
					<th>Student</th>
					<th>Given At</th>
					<th># of Reply</th>
					<th></th>
				</tr>
			</thead>
			<tbody>
			{{ range .HelpRequests }}
			<tr>
				<td>{{ .StudentName }}</td>
				<td>{{ formatTimeSince .GivenAt }} ago</td>
				<td>{{ .NumReply}}</td>
				<td><a href="/view_help_request?request_id={{.ID}}&uid={{$.UserID}}&role={{$.UserRole}}&password={{$.Password}}">View</a></td>
			</tr>
			{{ end }}
			</tbody>
		</table>
	</div>

	</body>
	</html>
`

var HELP_REQUEST_VIEW_TEMPLATE = `
<!DOCTYPE html>
	<html>
	<head>
	<title>Help Request</title>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/codemirror.min.js" integrity="sha512-hGVnilhYD74EGnPbzyvje74/Urjrg5LSNGx0ARG1Ucqyiaz+lFvtsXk/1jCwT9/giXP0qoXSlVDjxNxjLvmqAw==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/mode/python/python.min.js" integrity="sha512-/mavDpedrvPG/0Grj2Ughxte/fsm42ZmZWWpHz1jCbzd5ECv8CB7PomGtw0NAnhHmE/lkDFkRMupjoohbKNA1Q==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/mode/clike/clike.min.js" integrity="sha512-GAled7oA9WlRkBaUQlUEgxm37hf43V2KEMaEiWlvBO/ueP2BLvBLKN5tIJu4VZOTwo6Z4XvrojYngoN9dJw2ug==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/codemirror.min.css" integrity="sha512-6sALqOPMrNSc+1p5xOhPwGIzs6kIlST+9oGWlI4Wwcbj1saaX9J3uzO3Vub016dmHV7hM+bMi/rfXLiF5DNIZg==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/theme/monokai.min.css" integrity="sha512-R6PH4vSzF2Yxjdvb2p2FA06yWul+U0PDDav4b/od/oXf9Iw37zl10plvwOXelrjV2Ai7Eo3vyHeyFUjhXdBCVQ==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bulma/0.9.3/css/bulma.min.css" integrity="sha512-IgmDkwzs96t4SrChW29No3NXBIBv8baW490zk5aXvhCD8vuZM3yUSkbyTBcXohkySecyzIrUwiF/qV0cuPcL3Q==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<script src="https://kit.fontawesome.com/923539b4ee.js" crossorigin="anonymous"></script>
	<script src="https://code.jquery.com/jquery-3.6.0.min.js" integrity="sha256-/xUj+3OJU5yExlq6GSYGSHk7tPXikynS7ogEvDej/m4=" crossorigin="anonymous"></script>
	<script src="https://code.jquery.com/ui/1.12.1/jquery-ui.min.js" integrity="sha256-VazP97ZCwtekAsvgPBSUwPFKdrwD3unUfSGVYrahUqU=" crossorigin="anonymous"></script>
	<link rel="stylesheet" href="https://code.jquery.com/ui/1.12.1/themes/base/jquery-ui.css" />
	</head>
	<body>
		<div class="container">
			<h3 class="title is-3">Help Request</h3>
			<section class="section">
				<article class="message">
					<div class="message-header">
					<p>Help Request from {{.StudentName}} At ({{.GivenAt.Format "Jan 02, 2006 3:04:05 PM"}})</p>
					</div>
					<div class="message-body">
						{{.Explanation}}
					</div>
				</article>
				<textarea id="editor">{{ .Snapshot }}</textarea>
				<form action="/save_snapshot_feedback" method="POST">
					<label class="label">Feedback</label>
					<div class="field">
						<textarea class="textarea" placeholder="Write your feedback!" name="feedback"></textarea>
					</div>
					<div class="control">
						<input class="button" type="submit" value="Send Feedback">
					</div>
						
						
						<input type="hidden" name="snapshot_id" value="{{.SnapshotID}}">
						<input type="hidden" name="uid" value="{{.UserID}}">
						<input type="hidden" name="role" value="{{.UserRole}}">
						<input type="hidden" name="password" value="{{.Password}}">
					
				</form>
				<footer class="footer">
					<div class="content has-text-centered">
					<a href="/get_snapshot?snapshot_id={{.SnapshotID}}&uid={{$.UserID}}&role={{$.UserRole}}&password={{$.Password}}">View Snapshot</a>
					</div>
				</footer>
			</section>
		</div>
	</body>
	<script>
			var editor = document.getElementById("editor");
			var myCodeMirror = CodeMirror.fromTextArea(editor, {lineNumbers: true, mode: "{{getEditorMode .ProblemName}}", theme: "monokai", matchBrackets: true, indentUnit: 4, indentWithTabs: true, readOnly: "nocursor"});
			myCodeMirror.setSize("100%", 400)
			
		</script>
	</html>
`
var FEEDBACK_PROVISION_TEMPLATE = `
	<!DOCTYPE html>
	<html>
	<head>
	<title>Student Dashboard</title>
	<meta http-equiv="refresh" content="120" >
	<script src="https://kit.fontawesome.com/923539b4ee.js" crossorigin="anonymous"></script>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bulma/0.9.3/css/bulma.min.css" integrity="sha512-IgmDkwzs96t4SrChW29No3NXBIBv8baW490zk5aXvhCD8vuZM3yUSkbyTBcXohkySecyzIrUwiF/qV0cuPcL3Q==" crossorigin="anonymous" referrerpolicy="no-referrer" />

	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/codemirror.min.js" integrity="sha512-hGVnilhYD74EGnPbzyvje74/Urjrg5LSNGx0ARG1Ucqyiaz+lFvtsXk/1jCwT9/giXP0qoXSlVDjxNxjLvmqAw==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/mode/python/python.min.js" integrity="sha512-/mavDpedrvPG/0Grj2Ughxte/fsm42ZmZWWpHz1jCbzd5ECv8CB7PomGtw0NAnhHmE/lkDFkRMupjoohbKNA1Q==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/mode/clike/clike.min.js" integrity="sha512-GAled7oA9WlRkBaUQlUEgxm37hf43V2KEMaEiWlvBO/ueP2BLvBLKN5tIJu4VZOTwo6Z4XvrojYngoN9dJw2ug==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/codemirror.min.css" integrity="sha512-6sALqOPMrNSc+1p5xOhPwGIzs6kIlST+9oGWlI4Wwcbj1saaX9J3uzO3Vub016dmHV7hM+bMi/rfXLiF5DNIZg==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/theme/monokai.min.css" integrity="sha512-R6PH4vSzF2Yxjdvb2p2FA06yWul+U0PDDav4b/od/oXf9Iw37zl10plvwOXelrjV2Ai7Eo3vyHeyFUjhXdBCVQ==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<script src="https://code.jquery.com/jquery-3.6.0.min.js" integrity="sha256-/xUj+3OJU5yExlq6GSYGSHk7tPXikynS7ogEvDej/m4=" crossorigin="anonymous"></script>
	<script src="https://code.jquery.com/ui/1.12.1/jquery-ui.min.js" integrity="sha256-VazP97ZCwtekAsvgPBSUwPFKdrwD3unUfSGVYrahUqU=" crossorigin="anonymous"></script>
	<link rel="stylesheet" href="https://code.jquery.com/ui/1.12.1/themes/base/jquery-ui.css" />
	<script src="https://cdn.jsdelivr.net/npm/@creativebulma/bulma-collapsible"></script>
	</head>
	<body>
	<div class="container">
	<nav class="breadcrumb" aria-label="breadcrumbs">
	<ul>
	  <li>
		<a id="view-exercise-link" href="#">
		  <span class="icon is-small">
			<i class="fas fa-home" aria-hidden="true"></i>
		  </span>
		  <span>Exercises</span>
		</a>
	  </li>
	  <li>
		<a id="problem-dashboard-link" href="#">
		<span class="icon is-small">
			<i class="fas fa-book" aria-hidden="true"></i>
		  </span>
			<span>Problem Dashboard ({{.ProblemName}})</span>
		</a>
	   </li>
	  <li class="is-active">
		<a href="#">
			<span class="icon is-small">
				<i class="fas fa-puzzle-piece" aria-hidden="true"></i>
			</span>
		  <span>Student Dashboard</span>
		</a>
	  </li>
	</ul>
	</nav>
		<h2 class="title is-2">{{.StudentName}}'s Dashboard for {{.ProblemName}}</h2>
		<div class="tabs">
			<ul>
			<li><a href="/student_dashboard_code_snapshot?student_id={{.StudentID}}&problem_id={{.ProblemID}}&uid={{.UserID}}&role={{.UserRole}}{{if ne .Password ""}}&password={{.Password}}{{end}}">CodeSpace</a></li>
				<li class="is-active"><a>Feedback</a></li>
			</ul>
		</div>
		<div>
			<section class="section">
				{{range .Messages}}
					<article class="message">
						<div class="message-header">
						<p>{{if eq .Type 0}}{{.Name}} asked for help{{else if eq .Event "at_submission"}} Submission Snapshot taken {{else}} Regular Snapshot taken {{end}} at ({{.GivenAt.Format "Jan 02, 2006 3:04:05 PM"}})</p>
						</div>
						<div class="message-body">
							{{.Message}}
						</div>
						<div style="margin-left:20px;">
							
						<!--
							<div class="accordions">
								<h3>Code Snapshot</h3>
								<div>
									<textarea class="editor">{{ .Code }}</textarea>
								</div>
							</div>
						-->
							
							{{range $index, $el := .Feedbacks}}
								<article class="message" style="margin-left: 25px;">
									<div class="message-header">
									<p>Reply from {{.Name}} given at {{.GivenAt.Format "Jan 02, 2006 3:04:05 PM"}} </p>
									</div>
									<div class="message-body">
										<div class="columns">
											<div class="column is-three-quarters">
												<textarea class="editor">{{ .Feedback }}</textarea>
											</div>
											<div class="column">
												<a onclick="autoFeedbackSubmit('yes', {{.FeedbackID}})">
													<span style="font-size: 1.5em; {{if eq .CurrentUserVote "yes"}} color: green; {{end}}">
														<i class="fas fa-thumbs-up"></i>
													</span>
												</a>
												<span>
														{{.Upvote}}
												</span>
											</div>
											
										</div>

										<div id="feedback-block-{{ $index }}"></div>
										
									</div>
								</article>
							{{end}}

						</div>
					</article>
					
				{{end}}
			</section>
		</div>
		<script>
			$(document).ready(function(){
				$('#view-exercise-link').attr("href", "/view_exercises"+window.location.search);
				$('#problem-dashboard-link').attr("href", "/problem_dashboard"+window.location.search+"&problem_id={{.ProblemID}}");
		  	});
			var snapshotEditors = document.getElementsByClassName("editor");
				
			for (let i = 0; i<snapshotEditors.length; i++){
				var code = CodeMirror.fromTextArea(snapshotEditors[i], {lineNumbers: true, mode: "{{getEditorMode .ProblemName}}", theme: "monokai", matchBrackets: true, indentUnit: 4, indentWithTabs: true, readOnly: "nocursor"});
				code.setSize("100%", 500);
			}

			// TODO Get inline feedbacks and their labels.


			function autoFeedbackSubmit(backFeedback, fID) {
				$.ajax({
					url: "/save_snapshot_back_feedback",
					type: "POST",
					data:  {
						feedback: backFeedback,
						feedback_id: fID,
						uid: {{.UserID}},
						role: "{{.UserRole}}",
						{{if ne .Password ""}}password: "{{.Password}}",{{end}}
					},
					success: function(data){
						console.log("Success!")
					}
				});
				
				location.reload();
			}
		</script>

	</body>
	</html>
`
var PROBLEM_DASHBOARD_TEMPLATE = `
<!DOCTYPE html>
<html>
<head>
<title>Problem Dashboard</title>
<meta http-equiv="refresh" content="120" >
<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/codemirror.min.js" integrity="sha512-hGVnilhYD74EGnPbzyvje74/Urjrg5LSNGx0ARG1Ucqyiaz+lFvtsXk/1jCwT9/giXP0qoXSlVDjxNxjLvmqAw==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/mode/python/python.min.js" integrity="sha512-/mavDpedrvPG/0Grj2Ughxte/fsm42ZmZWWpHz1jCbzd5ECv8CB7PomGtw0NAnhHmE/lkDFkRMupjoohbKNA1Q==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/mode/clike/clike.min.js" integrity="sha512-GAled7oA9WlRkBaUQlUEgxm37hf43V2KEMaEiWlvBO/ueP2BLvBLKN5tIJu4VZOTwo6Z4XvrojYngoN9dJw2ug==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/codemirror.min.css" integrity="sha512-6sALqOPMrNSc+1p5xOhPwGIzs6kIlST+9oGWlI4Wwcbj1saaX9J3uzO3Vub016dmHV7hM+bMi/rfXLiF5DNIZg==" crossorigin="anonymous" referrerpolicy="no-referrer" />
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/theme/monokai.min.css" integrity="sha512-R6PH4vSzF2Yxjdvb2p2FA06yWul+U0PDDav4b/od/oXf9Iw37zl10plvwOXelrjV2Ai7Eo3vyHeyFUjhXdBCVQ==" crossorigin="anonymous" referrerpolicy="no-referrer" />
<script src="https://kit.fontawesome.com/923539b4ee.js" crossorigin="anonymous"></script>
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bulma/0.9.3/css/bulma.min.css" integrity="sha512-IgmDkwzs96t4SrChW29No3NXBIBv8baW490zk5aXvhCD8vuZM3yUSkbyTBcXohkySecyzIrUwiF/qV0cuPcL3Q==" crossorigin="anonymous" referrerpolicy="no-referrer" />
<script src="https://code.jquery.com/jquery-3.6.0.min.js" integrity="sha256-/xUj+3OJU5yExlq6GSYGSHk7tPXikynS7ogEvDej/m4=" crossorigin="anonymous"></script>
<script src="https://code.jquery.com/ui/1.12.1/jquery-ui.min.js" integrity="sha256-VazP97ZCwtekAsvgPBSUwPFKdrwD3unUfSGVYrahUqU=" crossorigin="anonymous"></script>
<link rel="stylesheet" href="https://code.jquery.com/ui/1.12.1/themes/base/jquery-ui.css" />
</head>
<body>
<div class="container">
<nav class="breadcrumb" aria-label="breadcrumbs">
<ul>
  <li>
	<a id="view-exercise-link" href="#">
	  <span class="icon is-small">
		<i class="fas fa-home" aria-hidden="true"></i>
	  </span>
	  <span>Exercises</span>
	</a>
  </li>
  <li class="is-active">
	<a href="#">
	  <span class="icon is-small">
		<i class="fas fa-book" aria-hidden="true"></i>
	  </span>
	  <span>Problem Dashboard</span>
	</a>
  </li>
</ul>
</nav>
	<h2 class="title is-2">Dashboard for {{.ProblemName}}</h2> 
	{{if eq .UserRole "teacher"}} {{if eq .IsActive true}}<button id="deactivate-button" class="button is-danger">Deactivate!</button>{{end}}{{end}}
	<div class="accordions">
		<h3>{{.ProblemName}}</h3>
		<div>
			<textarea id="editor">{{ .Code }}</textarea>
		</div>
	</div>
	<table class="table">
			<thead>
				<tr>
					<th>Active Students</th>
					<th>Help Requests</th>
					<th>Not Graded</th>
					<th>Correct</th>
					<th>Incorrect</th>
				</tr>
			</thead>
			<tbody>
				<tr>
					<td>{{.NumActive}}</td>
					<td>{{.NumHelpRequest}}</td>
					<td>{{.NumNotGraded}}</td>
					<td>{{.NumGradedCorrect}}</td>
					<td>{{.NumGradedIncorrect}}</td>
				</tr>
			</tbody>
	</table>

	<table class="table">
			<thead>
				<tr>
					<th>Student</th>
					<th>Active</th>
					<th>Coding Status</th>
					<th>Help Status</th>
					<th>Submission Status</th>
					<th>Tutoring Status</th>
				</tr>
			</thead>
			<tbody>
				{{range .StudentInfo}}
				<tr>
					<td>{{if ne .CodingStat "Idle"}}<a href="/student_dashboard_code_snapshot?student_id={{.StudentID}}&problem_id={{$.ProblemID}}&uid={{$.UserID}}&role={{$.UserRole}}{{if ne $.Password ""}}&password={{$.Password}}{{end}}">{{.StudentName}}</a>{{else}}{{.StudentName}}{{end}}</td>
					<td>{{if and (eq $.IsActive true) (ne .CodingStat "Idle") (ne .LastUpdatedAt.IsZero true) }}<a href="/student_dashboard_code_snapshot?student_id={{.StudentID}}&problem_id={{$.ProblemID}}&uid={{$.UserID}}&role={{$.UserRole}}{{if ne $.Password ""}}&password={{$.Password}}{{end}}">{{ formatTimeSince .LastUpdatedAt }} ago</a>{{end}}</td>
					<td>{{.CodingStat}}</td>
					<td>{{if ne .HelpStat ""}}<a href="/student_dashboard_feedback_provision?student_id={{.StudentID}}&problem_id={{$.ProblemID}}&uid={{$.UserID}}&role={{$.UserRole}}{{if ne $.Password ""}}&password={{$.Password}}{{end}}">{{.HelpStat}}</a>{{end}}</td>
					<td>{{if ne .SubmissionStat ""}}<a href="/student_dashboard_submissions?student_id={{.StudentID}}&problem_id={{$.ProblemID}}&uid={{$.UserID}}&role={{$.UserRole}}{{if ne $.Password ""}}&password={{$.Password}}{{end}}">{{.SubmissionStat}}</a>{{end}}</td>
					<td>{{.TutoringStat}}</td>
				</tr>
				{{end}}
			</tbody>
	</table>
	<script>
		var editor = document.getElementById("editor");
		var myCodeMirror = CodeMirror.fromTextArea(editor, {lineNumbers: true, mode: get_editor_mode({{.ProblemName}}), theme: "monokai", matchBrackets: true, indentUnit: 4, indentWithTabs: true, readOnly: "nocursor"});
		myCodeMirror.setSize("100%", 400)
		function get_editor_mode(filename) {
			filename = filename.toLowerCase();
			if (filename.endsWith('.py')) {
				return "python";
			}
			if (filename.endsWith('.java')) {
				return "text/x-java";
			}
			if (filename.endsWith('.cpp') || filename.endsWith('.c++') || filename.endsWith('.c')) {
				return "text/x-c++src";
			}
			return "text";
		  }
		  $(document).ready(function(){
			$('#view-exercise-link').attr("href", "/view_exercises"+window.location.search);
			$('#deactivate-button').click(function(){
				if (confirm("Deactivate the problem?") == true) {
					$.post("/teacher_deactivates_problems", {filename: {{.ProblemName}}, uid: {{.UserID}}, role: {{.UserRole}}{{if ne .Password ""}}, password: {{.Password}}{{end}} })
					.done(function(data){
						if (data == "-1"){
							alert("Couldn't deactivate the problem! Please try again!");
						} else {
							alert("Problem deactiavated!");
							window.location.reload();
						}
					})
					.fail(function(){
						alert("Couldn't deactivate the problem! Please try again!");
					});
				}
			});
		  });
		  $(".accordions").accordion({ header: "h3", active: false, collapsible: true });
		  $(".accordions").show();
	</script>
</body>
</html>
`
var PROBLEM_LIST_TEMPLATE = `
<!DOCTYPE html>
<html>
<head>
<title>Exercises</title>
<meta http-equiv="refresh" content="120" >
<style>
.switch {
  position: relative;
  display: inline-block;
  width: 60px;
  height: 34px;
}

.switch input { 
  opacity: 0;
  width: 0;
  height: 0;
}

.slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: #ccc;
  -webkit-transition: .4s;
  transition: .4s;
}

.slider:before {
  position: absolute;
  content: "";
  height: 26px;
  width: 26px;
  left: 4px;
  bottom: 4px;
  background-color: white;
  -webkit-transition: .4s;
  transition: .4s;
}

input:checked + .slider {
  background-color: #2196F3;
}

input:focus + .slider {
  box-shadow: 0 0 1px #2196F3;
}

input:checked + .slider:before {
  -webkit-transform: translateX(26px);
  -ms-transform: translateX(26px);
  transform: translateX(26px);
}

/* Rounded sliders */
.slider.round {
  border-radius: 34px;
}

.slider.round:before {
  border-radius: 50%;
}
</style>
<script src="https://kit.fontawesome.com/923539b4ee.js" crossorigin="anonymous"></script>
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bulma/0.9.3/css/bulma.min.css" integrity="sha512-IgmDkwzs96t4SrChW29No3NXBIBv8baW490zk5aXvhCD8vuZM3yUSkbyTBcXohkySecyzIrUwiF/qV0cuPcL3Q==" crossorigin="anonymous" referrerpolicy="no-referrer" />
<script src="https://code.jquery.com/jquery-3.6.0.min.js" integrity="sha256-/xUj+3OJU5yExlq6GSYGSHk7tPXikynS7ogEvDej/m4=" crossorigin="anonymous"></script>
</head>
<body>
<div class="container">
	<nav class="breadcrumb" aria-label="breadcrumbs">
	<ul>
		<li class="is-active">
		<a href="#">
			<span class="icon is-small">
			<i class="fas fa-home" aria-hidden="true"></i>
			</span>
			<span>Exercises</span>
		</a>
		</li>
	</ul>
	</nav>
	<h2 class="title is-2">Exercises</h2>
	<a id="new-problem" class="button is-success" href="">
		<span class="icon is-small">
		<i class="fa-solid fa-plus"></i>
		</span>
		<span>Broadcast New Exercise</span>
	</a>
	Peer Tutoring: 
	<label class="switch">
		<input id="peer_tutoring_button" type="checkbox">
		<span class="slider round"></span>
	</label>
	<table class="table">
			<thead>
				<tr>
					<th>Filename</th>
					<th>Posted At</th>
					<th>Attendance</th>
					<th>Active Students</th>
					<th>Help Requests</th>
					<th>Correct</th>
					<th>Incorrect</th>
					<th>Not Graded</th>
				</tr>
			</thead>
			<tbody>
				{{range .Problems}}
				<tr {{if eq .IsActive true}}class="is-selected"{{end}}>
					<td><a href="/problem_dashboard?problem_id={{.ID}}&uid={{$.UserID}}&role={{$.UserRole}}{{if ne $.Password ""}}&password={{$.Password}}{{end}}">{{.Filename}}</a></td>
					<td>{{ .UploadedAt.Format "Jan 02, 2006 3:04:05 PM" }}</td>
					<td>{{.Attendance}}</td>
					<td>{{.NumActive}}</td>
					<td>{{.NumHelpRequest}}</td>
					<td>{{.NumGradedCorrect}}</td>
					<td>{{.NumGradedIncorrect}}</td>
					<td>{{.NumNotGraded}}</td>
				</tr>
				{{end}}
			</tbody>
	</table>
<script>
$(document).ready(function(){
	{{if eq .PeerTutorAllowed true}}$('#peer_tutoring_button').prop('checked', true);{{end}}
	$('#new-problem').attr("href", "/teacher_web_broadcast"+window.location.search);
	$('#peer_tutoring_button').change(function(){
		var val = document.getElementById('peer_tutoring_button').checked;
		var valInt = 0;
		if (val == true)
			valInt = 1;
		$.post("/set_peer_tutor", {turn_on: valInt, uid: {{.UserID}}, role: {{.UserRole}}{{if ne .Password ""}}, password: {{.Password}}{{end}} }, function(data, status){
		});
	});
	
});
</script>
</body>
</html>
`

var SUBMISSION_VIEW_TEMPLATE = `
	<!DOCTYPE html>
	<html>
	<head>
	<title>Student Dashboard</title>
	<meta http-equiv="refresh" content="120" >
	<script src="https://kit.fontawesome.com/923539b4ee.js" crossorigin="anonymous"></script>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bulma/0.9.3/css/bulma.min.css" integrity="sha512-IgmDkwzs96t4SrChW29No3NXBIBv8baW490zk5aXvhCD8vuZM3yUSkbyTBcXohkySecyzIrUwiF/qV0cuPcL3Q==" crossorigin="anonymous" referrerpolicy="no-referrer" />

	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/codemirror.min.js" integrity="sha512-hGVnilhYD74EGnPbzyvje74/Urjrg5LSNGx0ARG1Ucqyiaz+lFvtsXk/1jCwT9/giXP0qoXSlVDjxNxjLvmqAw==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/mode/python/python.min.js" integrity="sha512-/mavDpedrvPG/0Grj2Ughxte/fsm42ZmZWWpHz1jCbzd5ECv8CB7PomGtw0NAnhHmE/lkDFkRMupjoohbKNA1Q==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/mode/clike/clike.min.js" integrity="sha512-GAled7oA9WlRkBaUQlUEgxm37hf43V2KEMaEiWlvBO/ueP2BLvBLKN5tIJu4VZOTwo6Z4XvrojYngoN9dJw2ug==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/codemirror.min.css" integrity="sha512-6sALqOPMrNSc+1p5xOhPwGIzs6kIlST+9oGWlI4Wwcbj1saaX9J3uzO3Vub016dmHV7hM+bMi/rfXLiF5DNIZg==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/theme/monokai.min.css" integrity="sha512-R6PH4vSzF2Yxjdvb2p2FA06yWul+U0PDDav4b/od/oXf9Iw37zl10plvwOXelrjV2Ai7Eo3vyHeyFUjhXdBCVQ==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<script src="https://code.jquery.com/jquery-3.6.0.min.js" integrity="sha256-/xUj+3OJU5yExlq6GSYGSHk7tPXikynS7ogEvDej/m4=" crossorigin="anonymous"></script>
	<script src="https://code.jquery.com/ui/1.12.1/jquery-ui.min.js" integrity="sha256-VazP97ZCwtekAsvgPBSUwPFKdrwD3unUfSGVYrahUqU=" crossorigin="anonymous"></script>
	<link rel="stylesheet" href="https://code.jquery.com/ui/1.12.1/themes/base/jquery-ui.css" />
	</head>
	<body>
	<div class="container">
	<nav class="breadcrumb" aria-label="breadcrumbs">
		<ul>
		<li>
			<a id="view-exercise-link" href="#">
			<span class="icon is-small">
				<i class="fas fa-home" aria-hidden="true"></i>
			</span>
			<span>Exercises</span>
			</a>
		</li>
		<li>
			<a id="problem-dashboard-link" href="#">
			<span class="icon is-small">
				<i class="fas fa-book" aria-hidden="true"></i>
			</span>
				<span>Problem Dashboard ({{.ProblemName}})</span>
			</a>
		</li>
		<li class="is-active">
			<a href="#">
				<span class="icon is-small">
					<i class="fas fa-puzzle-piece" aria-hidden="true"></i>
				</span>
			<span>Student Dashboard</span>
			</a>
		</li>
		</ul>
	</nav>
		<h2 class="title is-2">{{.StudentName}}'s Submissions for {{.ProblemName}}</h2>
		<div class="tabs">
			<ul>
				<li><a href="/student_dashboard_code_snapshot?student_id={{.StudentID}}&problem_id={{.ProblemID}}&uid={{.UserID}}&role={{.UserRole}}{{if ne .Password ""}}&password={{.Password}}{{end}}">Code Snapshot</a></li>
				<li><a href="/student_dashboard_feedback_provision?student_id={{.StudentID}}&problem_id={{.ProblemID}}&uid={{.UserID}}&role={{.UserRole}}{{if ne .Password ""}}&password={{.Password}}{{end}}" >Feedback</a></li>
				<li class="is-active"><a>Submissions</a></li>
			</ul>
		</div>
		<div>
			{{range .Submissions}}
			<div class="box">
			<h4 class="title is-4">Submitted at {{.SubmittedAt.Format "Jan 02, 2006 3:04:05 PM"}}</h4>
			{{if eq .Grade ""}} Not Graded {{else}} Graded {{if eq .Grade "correct"}} <span class="tag is-success">correct</span> {{else if eq .Grade "incorrect"}} <span class="tag is-danger">incorrect</span> {{else}} {{.Grade}} {{end}} {{end}}
			<div class="accordions">
				<h3>Code</h3>
				<div>
					<textarea class="editor" id="editor-{{.ID}}">{{ .Code }}</textarea>
				</div>
			</div>
			
			{{if eq .Grade ""}}
			<div class="columns">
				<div class="column is-three-quarters"><input  class="input is-info" id="{{.ID}}" type="text" placeholder="Provide your feedback!"></div>
				<div class="column"><button  class="button is-success" onclick="sendGrade({{.ID}}, {{.SnapshotID}}, 'correct')">Correct</button></div>
				<div class="column"><button  class="button is-danger" onclick="sendGrade({{.ID}}, {{.SnapshotID}}, 'incorrect')">Incorrect</button></div>
			</div>
			{{end}}
			</div>
			{{end}}
		</div>
		<script>
			$(document).ready(function(){
				$('#view-exercise-link').attr("href", "/view_exercises"+window.location.search);
				$('#problem-dashboard-link').attr("href", "/problem_dashboard"+window.location.search+"&problem_id={{.ProblemID}}");
			});
			var snapshotEditors = document.getElementsByClassName("editor");
				
			for (let i = 0; i<snapshotEditors.length; i++){
				var code = CodeMirror.fromTextArea(snapshotEditors[i], {lineNumbers: true, mode: "{{getEditorMode .ProblemName}}", theme: "monokai", matchBrackets: true, indentUnit: 4, indentWithTabs: true, readOnly: "nocursor"});
				code.setSize("100%", "auto");
			}
			

			function sendGrade(submission_id, snapshot_id, grade) {
				var feedback = $('#'+submission_id).val().trim();
				var code = $('#editor-'+submission_id).val();
				$.post("/teacher_grades", {content: code, changed: "", decision: grade, sid: submission_id, uid: {{.UserID}}, role: {{.UserRole}}{{if ne .Password ""}}, password: {{.Password}}{{end}}  }, function(data, status){
					if (status == "success"){
						if (feedback != "") {
							$.post("/save_snapshot_feedback", {snapshot_id: snapshot_id, feedback: feedback, uid: {{.UserID}}, role: {{.UserRole}}{{if ne .Password ""}}, password: {{.Password}}{{end}} }, function(data1, status1){
							});
						}
						alert("Graded successfully!");
						window.location.reload();
					} else {
						alert("Could not grade the submission. Please try again!");
					}
				});
			}
			$(".accordions").accordion({ header: "h3", active: false, collapsible: true });
		</script>

	</body>
	</html>
`
var TEACHER_LOGIN = `
<!DOCTYPE html>
<html>
   <head>
      <meta charset="utf-8">
      <meta name="viewport" content="width=device-width, initial-scale=1">
      <title>Teacher Login</title>
      <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/font-awesome/4.7.0/css/font-awesome.min.css">
      <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bulma/0.6.0/css/bulma.min.css">
	  <script src="https://code.jquery.com/jquery-3.6.0.min.js" integrity="sha256-/xUj+3OJU5yExlq6GSYGSHk7tPXikynS7ogEvDej/m4=" crossorigin="anonymous"></script>
   </head>
   <body>
   <div class="container">
      <section class="section">      
       <div class="columns">
       <div class="column is-4 is-offset-4">
		  <div class="field">
		  <p class="control has-icons-left has-icons-right">
		    <input id="name" class="input" type="email" placeholder="Name">
		    <span class="icon is-small is-left">
		      <i class="fa fa-user"></i>
		    </span>
		    <span class="icon is-small is-right">
		      <i class="fa fa-check"></i>
		    </span>
		  </p>
		</div>
		<div class="field">
		  <p class="control has-icons-left">
		    <input id="password" class="input" type="password" placeholder="Password">
		    <span class="icon is-small is-left">
		      <i class="fa fa-lock"></i>
		    </span>
		  </p>
		</div>
		<div class="field">
		  <p class="control">
		    <button id="login" class="button is-success">
		      Login
		    </button>
		  </p>
		</div>
      </div>         
       </div>
      </section>
	</div>
	  <script>
	  	$(document).ready(function(){
			$('#login').click(function(){
				var name = $('#name').val().trim();
				var pass = $('#password').val().trim();
				if(name == "" || pass == "") {
					alert("Please enter both name and password!");
				} else {
					$.post("/teacher_signin_complete", {username: name, password: pass}, function(data, status){
						if (status == "success"){
							window.location.replace("/view_exercises?role=teacher&uid="+data);
						} else {
							alert("Unauthorized access");
						}
					});
				}
			});
		});
	  </script>
   </body>
</html>
`
var PROBLEM_FILE_UPLOAD_VIEW = `
<!DOCTYPE html>
<html>
   <head>
      <meta charset="utf-8">
      <meta name="viewport" content="width=device-width, initial-scale=1">
      <title>Broadcast Problem</title>
      <script src="https://kit.fontawesome.com/923539b4ee.js" crossorigin="anonymous"></script>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bulma/0.9.3/css/bulma.min.css" integrity="sha512-IgmDkwzs96t4SrChW29No3NXBIBv8baW490zk5aXvhCD8vuZM3yUSkbyTBcXohkySecyzIrUwiF/qV0cuPcL3Q==" crossorigin="anonymous" referrerpolicy="no-referrer" />

	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/codemirror.min.js" integrity="sha512-hGVnilhYD74EGnPbzyvje74/Urjrg5LSNGx0ARG1Ucqyiaz+lFvtsXk/1jCwT9/giXP0qoXSlVDjxNxjLvmqAw==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/mode/python/python.min.js" integrity="sha512-/mavDpedrvPG/0Grj2Ughxte/fsm42ZmZWWpHz1jCbzd5ECv8CB7PomGtw0NAnhHmE/lkDFkRMupjoohbKNA1Q==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/mode/clike/clike.min.js" integrity="sha512-GAled7oA9WlRkBaUQlUEgxm37hf43V2KEMaEiWlvBO/ueP2BLvBLKN5tIJu4VZOTwo6Z4XvrojYngoN9dJw2ug==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/codemirror.min.css" integrity="sha512-6sALqOPMrNSc+1p5xOhPwGIzs6kIlST+9oGWlI4Wwcbj1saaX9J3uzO3Vub016dmHV7hM+bMi/rfXLiF5DNIZg==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/theme/monokai.min.css" integrity="sha512-R6PH4vSzF2Yxjdvb2p2FA06yWul+U0PDDav4b/od/oXf9Iw37zl10plvwOXelrjV2Ai7Eo3vyHeyFUjhXdBCVQ==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	  <script src="https://code.jquery.com/jquery-3.6.0.min.js" integrity="sha256-/xUj+3OJU5yExlq6GSYGSHk7tPXikynS7ogEvDej/m4=" crossorigin="anonymous"></script>
   </head>
   <body>
   <div class="container">
   <nav class="breadcrumb" aria-label="breadcrumbs">
	<ul>
		<li>
			<a id="view-exercise-link" href="#">
			<span class="icon is-small">
				<i class="fas fa-home" aria-hidden="true"></i>
			</span>
			<span>Exercises</span>
			</a>
		</li>
		<li class="is-active">
			<a href="#">
			<span class="icon is-small">
				<i class="fas fa-book" aria-hidden="true"></i>
			</span>
			<span>Problem Broadcast</span>
			</a>
		</li>
		</ul>
		</nav>
   <div id="problem" class="file is-centered is-boxed is-success has-name">
		<label class="file-label">
			<input class="file-input" type="file" name="resume">
			<span class="file-cta">
			<span class="file-icon">
				<i class="fas fa-upload"></i>
			</span>
			<span class="file-label">
				Select Exercise File
			</span>
			</span>
		</label>
	</div>
   <div style="visibility:hidden;" id="editor-area">
	<article class="message">
			<div class="message-header">
				<p><span id="filename"></span></p>
			</div>
			<div class="message-body">
				<div>
					<textarea id="editor"></textarea>
				</div>
			</div>
		</article>
   </div>
   <button style="visibility:hidden" id="submit" class="button is-success is-rounded">
		<span class="icon is-small">
			<i class="fas fa-check"></i>
		</span>
		<span>Broadcast</span>
	</button>
	<input type="hidden" id="points" value="">
	<input type="hidden" id="effort" value="">
	<input type="hidden" id="attempt" value="">
	<input type="hidden" id="tag" value="">
  </div>
	<script>
	$(document).ready(function(){
		$('#view-exercise-link').attr("href", "/view_exercises"+window.location.search);
	  });
	document.querySelector('#problem input[type=file]').onchange = function(){
		document.querySelector('#problem').style.visibility = "hidden";
		var file = this.files[0];
		document.querySelector('#filename').textContent = file.name;
		var reader = new FileReader();
		reader.onload = function(progressEvent){
	  
		  // By lines
		  var lines = this.result.split('\n');
		  var firstLine = lines[0];
		  lines.splice(0, 1);
		  var content = lines.join('\n');
		  if (firstLine.length == 0 || (firstLine[0]!='#' && !firstLine.startsWith('//') )){
			  alert("Invalid problem header!");
			return;
		  }
		  var prefix = '';
		  if (firstLine[0] == '#') {
			  prefix = '#';
			firstLine.replace("#", '');
		  } else {
			  prefix = '//';
			firstLine.replace("//", '');
		  }
		  params = get_problem_info(firstLine);
		//   alert(params[1]+" Points, "+params[2]+" for effort. Maximum attempts: " + params[3]);
		 $('#editor-area').css('visibility', 'visible');
		  document.querySelector('#editor').textContent = prefix + ' ' + params[1]+" Points, "+params[2]+" for effort. Maximum attempts: " + params[3] + "\n" + content;
		  var editor = document.getElementById("editor");
		  var myCodeMirror = CodeMirror.fromTextArea(editor, {lineNumbers: true, mode: get_editor_mode(file.name), theme: "monokai", matchBrackets: true, indentUnit: 4, indentWithTabs: true, readOnly: "nocursor"});
		  myCodeMirror.setSize("100%", 400)
		  $('#submit').css('visibility', 'visible');
		  $('#points').val(params[1]);
		  $('#effort').val(params[2]);
		  $('#attempt').val(params[3]);
		  $('#tag').val(params[4]);
		};
		reader.readAsText(file);
	  };
	  
	  function get_problem_info(content) {
		  let regexpNames =  /\s*(\d+)\s+(\d+)\s+(\d+)(?:\s+(\w.*))?/mg;
		let match = regexpNames.exec(content);
		return match;
	  }
	  function get_editor_mode(filename) {
		filename = filename.toLowerCase();
		if (filename.endsWith('.py')) {
			return "python";
		}
		if (filename.endsWith('.java')) {
			return "text/x-java";
		}
		if (filename.endsWith('.cpp') || filename.endsWith('.c++') || filename.endsWith('.c')) {
			return "text/x-c++src";
		}
		return "text";
	  }
	$(document).ready(function() {
		$.ajaxSetup({
			xhrFields: {
			  withCredentials: true
			}
		});
		$('#submit').click(function() {
			var editor = document.querySelector('.CodeMirror').CodeMirror;
			var uid = new URLSearchParams(window.location.search).get('uid');
			var points = $('#points').val();
			var effort = $('#effort').val();
			var attempt = $('#attempt').val();
			var tag = $('#tag').val();
			var filename = $('#filename').text();
			$.post("/teacher_broadcasts", {role: "teacher", uid: uid, content: editor.getValue(), answer: "", merit: points, effort: effort, attempts: attempt, tag: tag, filename: filename, exact_answer: true}, function(data, status){
				if (status == "success"){
					alert("Exercise broadcasted successfully!");
					window.location.replace("/view_exercises?role=teacher&uid="+uid);
				} else {
					alert("Failed to broadcast. Try agian!");
				}
			});
		});
	});
	
	</script>
   </body>
</html>
`
var CODE_SNAPSHOT_TAB_TEMPLATE = `
	<!DOCTYPE html>
	<html>
	<head>
	<title>Student Dashboard</title>
	<meta http-equiv="refresh" content="120" >
	<script src="https://kit.fontawesome.com/923539b4ee.js" crossorigin="anonymous"></script>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bulma/0.9.3/css/bulma.min.css" integrity="sha512-IgmDkwzs96t4SrChW29No3NXBIBv8baW490zk5aXvhCD8vuZM3yUSkbyTBcXohkySecyzIrUwiF/qV0cuPcL3Q==" crossorigin="anonymous" referrerpolicy="no-referrer" />

	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/codemirror.min.js" integrity="sha512-hGVnilhYD74EGnPbzyvje74/Urjrg5LSNGx0ARG1Ucqyiaz+lFvtsXk/1jCwT9/giXP0qoXSlVDjxNxjLvmqAw==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/mode/python/python.min.js" integrity="sha512-/mavDpedrvPG/0Grj2Ughxte/fsm42ZmZWWpHz1jCbzd5ECv8CB7PomGtw0NAnhHmE/lkDFkRMupjoohbKNA1Q==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/mode/clike/clike.min.js" integrity="sha512-GAled7oA9WlRkBaUQlUEgxm37hf43V2KEMaEiWlvBO/ueP2BLvBLKN5tIJu4VZOTwo6Z4XvrojYngoN9dJw2ug==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/codemirror.min.css" integrity="sha512-6sALqOPMrNSc+1p5xOhPwGIzs6kIlST+9oGWlI4Wwcbj1saaX9J3uzO3Vub016dmHV7hM+bMi/rfXLiF5DNIZg==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/theme/monokai.min.css" integrity="sha512-R6PH4vSzF2Yxjdvb2p2FA06yWul+U0PDDav4b/od/oXf9Iw37zl10plvwOXelrjV2Ai7Eo3vyHeyFUjhXdBCVQ==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<script src="https://code.jquery.com/jquery-3.6.0.min.js" integrity="sha256-/xUj+3OJU5yExlq6GSYGSHk7tPXikynS7ogEvDej/m4=" crossorigin="anonymous"></script>
	<script src="https://code.jquery.com/ui/1.12.1/jquery-ui.min.js" integrity="sha256-VazP97ZCwtekAsvgPBSUwPFKdrwD3unUfSGVYrahUqU=" crossorigin="anonymous"></script>
	<link rel="stylesheet" href="https://code.jquery.com/ui/1.12.1/themes/base/jquery-ui.css" />
	<script src="https://cdn.jsdelivr.net/npm/@creativebulma/bulma-collapsible"></script>
	<style>
		#code-snapshot {
			background: darkseagreen;
			padding: 20px;
			margin: 30px;
		}
		#ask-for-help {
			background: cornflowerblue;
			margin: 30px;
		}
		#submission {
			background: greenyellow;
			padding: 20px;
			margin: 30px;
		}
		.wrapper {
			display: grid;
			grid-template-columns: repeat(2, 1fr);
			gap: 20px;
		}
	</style>
	</head>
	<body>
	<div class="container">
	<nav class="breadcrumb" aria-label="breadcrumbs">
	<ul>
	  <li>
		<a id="view-exercise-link" href="#">
		  <span class="icon is-small">
			<i class="fas fa-home" aria-hidden="true"></i>
		  </span>
		  <span>Exercises</span>
		</a>
	  </li>
	  <li>
		<a id="problem-dashboard-link" href="#">
		<span class="icon is-small">
			<i class="fas fa-book" aria-hidden="true"></i>
		  </span>
			<span>Problem Dashboard ({{ .Feedback.ProblemName}})</span>
		</a>
	   </li>
	  <li class="is-active">
		<a href="#">
			<span class="icon is-small">
				<i class="fas fa-puzzle-piece" aria-hidden="true"></i>
			</span>
		  <span>Student Dashboard</span>
		</a>
	  </li>
	</ul>
	</nav>
		<h2 class="title is-2">{{ .Feedback.StudentName}}'s Dashboard for {{ .Feedback.ProblemName}}</h2>
		<div class="tabs">
			<ul>
				<li class="is-active"><a>CodeSpace</a></li>
				<li><a href="/student_dashboard_feedback_provision?student_id={{.Feedback.StudentID}}&problem_id={{.Feedback.ProblemID}}&uid={{.Feedback.UserID}}&role={{.Feedback.UserRole}}{{if ne .Feedback.Password ""}}&password={{.Feedback.Password}}{{end}}">Feedback</a></li>
			</ul>
		</div>
		<div id="code-snapshot">
			<div class="message-header">
				<p>Latest Code Snapshot at {{.Feedback.LastSnapshot.LastUpdated.Format "Jan 02, 2006 3:04:05 PM"}}</p>
				<button class="button is-info" onclick="codeSnapshotFeedback({{ .Feedback.LastSnapshot.Code }}, {{ .Feedback.UserID }})" style="margin-right:10px;" >Send Feedback</button>
			</div>	
			<textarea id="snapshot-editor"> {{ .Feedback.LastSnapshot.Code }} </textarea>
			<div id="code-snapshot-feedback-block"></div> 
		</div>

		<div id="ask-for-help">
			<section class="section" style="padding: 20px;">
				{{range $index, $el := .Feedback.Messages}}
					<article class="message">
						<div class="message-header">
							<p>{{if eq .Type 0}}{{.Name}} asked for help{{else if eq .Event "at_submission"}} Submission Snapshot taken {{else}} Regular Snapshot taken {{end}} at ({{.GivenAt.Format "Jan 02, 2006 3:04:05 PM"}})</p>
							<button class="button is-info" onclick="messageFeedback( {{ $index }} ,{{ .Code }} , {{ .ID }})" style="margin-right:10px;" >Send Feedback</button>
						</div>
						<div class="message-body">
							<h3>Student says: </h3> {{.Message}}
						</div>
						<div style="background: cornflowerblue;">
							
							<div class="accordions">
								
								<div>
									<textarea class="feedback-editor">{{ .Code }}</textarea>
								</div>
							</div>
							<div id="feedback-block-{{ $index }}"></div>

						</div>
					</article>
					
				{{end}}
			</section>
		</div>

		<div id="submission">
			{{range $index, $el := .Submission.Submissions}}
				<div class="box" style="padding: 0px">
					<div class="message-header">
						<p>Submitted at {{.SubmittedAt.Format "Jan 02, 2006 3:04:05 PM"}}</p>
					</div>
					<div class="message-body">
					{{if eq .Grade ""}} Not Graded {{else}} Graded {{if eq .Grade "correct"}} <span class="tag is-success">correct</span> {{else if eq .Grade "incorrect"}} <span class="tag is-danger">incorrect</span> {{else}} {{.Grade}} {{end}} {{end}}
					<h3>Code</h3>
					</div>
					
					<div>
						<textarea class="submission-editor" id="editor-{{.ID}}">{{ .Code }}</textarea>
					</div>
				
				{{if eq .Grade ""}}
					<div class="columns">
						<div class="column is-three-quarters">Grade this Submissions.</div>
						<div class="column"><button  class="button is-success" onclick="sendGrade( {{ $index }}, {{.ID}}, {{.SnapshotID}},{{ .Code }}, 'correct')">Correct</button></div>
						<div class="column"><button  class="button is-danger" onclick="sendGrade( {{ $index }}, {{.ID}}, {{.SnapshotID}}, {{ .Code }}, 'incorrect')">Incorrect</button></div>
					</div>
				{{end}}
				</div>
			{{end}}
		</div>


		<script>
			$(document).ready(function(){
				$('#view-exercise-link').attr("href", "/view_exercises"+window.location.search);
				$('#problem-dashboard-link').attr("href", "/problem_dashboard"+window.location.search+"&problem_id={{.Feedback.ProblemID}}");
		  	});

			var snapshotCodeChanged = "";
			var snapshotcode = CodeMirror.fromTextArea(document.getElementById("snapshot-editor"), {
				lineNumbers: true, 
				mode: "{{getEditorMode .Feedback.ProblemName}}", 
				theme: "monokai", 
				matchBrackets: true, 
				indentUnit: 4, 
				indentWithTabs: true, 
				readOnly: false
			});
			snapshotcode.setSize("100%", 500);
			snapshotcode.on('change', (snapshotcode) => {
				snapshotCodeChanged = snapshotcode.doc.getValue()
			});
			var snapshotCounter = 0;
			var feedbackCounter = 0;
			
			var feedbackEditors = document.getElementsByClassName("feedback-editor");
			var feedbackChangedCode = []
			for (let i = 0; i<feedbackEditors.length; i++){
				var code = CodeMirror.fromTextArea(feedbackEditors[i], {lineNumbers: true, mode: "{{getEditorMode .Feedback.ProblemName}}", theme: "monokai", matchBrackets: true, indentUnit: 4, indentWithTabs: true, readOnly: false});
				code.setSize("100%", 500);
				code.on('change', (code) => {
					feedbackChangedCode[i] = code.doc.getValue()
				});
			}

			var submissionEditors = document.getElementsByClassName("submission-editor");
			var submissionChangedCode = []
			for (let i = 0; i<submissionEditors.length; i++){
				var code = CodeMirror.fromTextArea(submissionEditors[i], {lineNumbers: true, mode: "{{getEditorMode .Feedback.ProblemName}}", theme: "monokai", matchBrackets: true, indentUnit: 4, indentWithTabs: true, readOnly: false});
				code.setSize("100%", 500);
				code.on('change', (code) => {
					submissionChangedCode[i] = code.doc.getValue()
				});
			}

			function runNLP(student_code, ta_code, user_id, write, counter) {
				console.log("Insert here: ", write);
				return $.ajax({
					url: "http://127.0.0.1:5000/feedback_classify",
					type: "POST",
					data: JSON.stringify({
							student_code: student_code,
							ta_code: ta_code,
							user_id: user_id
					}),
					"headers": {
						"Content-Type": "application/json"
					},
					success: function(data){
						var pre = ''
						var str = '';
						data.results.forEach(function(item){
							pre += '<div class="card" style="background-color: lightgray; margin-top: 15px"><div class="card-content"><ul style="margin:5px">'
							for (const key in item) {
								if (key == "output") {
									str = '<li><ul style="list-style-type:square">'
									item.output.forEach(function(suggestion){
										str += '<li>' + suggestion + '</li>'
									});
									str += '</ul></li>'
									pre += str
								} else {
									pre += '<li>' + key + ' : ' + item[key] + '</li>'
								}
							}
							pre += '</ul></div></div>'
						})
						$(write).html("");
						$('<div class="wrapper">' + pre +  '</div>').appendTo( write )

						if (counter == 0 ){
							alert("Feedback is not sent to student yet. Click on Send Feedback again.");
						}
					},
					error: function(err) {
						alert(JSON.stringify(err));
					}
				});
			}

			function codeSnapshotFeedback(code, user_id) {
				runNLP(code, snapshotCodeChanged, user_id, '#code-snapshot-feedback-block',snapshotCounter);

				if (snapshotCounter !== 0) {
					$.post("/save_snapshot_feedback", {feedback: snapshotCodeChanged, snapshot_id: {{.Feedback.LastSnapshot.ID}}, uid: {{ .Feedback.UserID}}, role: {{ .Feedback.UserRole}}{{if ne .Feedback.Password ""}}, password: {{ .Feedback.Password}}{{end}}  }, function(data, status){
						if (status == "success"){
							alert("Feedback posted successfully!");
							window.location.replace("/student_dashboard_feedback_provision?student_id={{ .Feedback.StudentID}}&problem_id={{ .Feedback.ProblemID}}&uid={{ .Feedback.UserID}}&role={{.Feedback.UserRole}}{{if ne .Feedback.Password ""}}&password={{.Feedback.Password}}{{end}}");
						} else {
							alert("Could not post the feedback. Please try again!");
						}
					});
				}
				snapshotCounter ++;
			}

			function messageFeedback(i,code,message_id){
				runNLP(code, feedbackChangedCode[i], {{ .Feedback.UserID }}, '#feedback-block-'+i, feedbackCounter );
				if (feedbackCounter != 0 ) {
					$.post("/save_message_feedback", {feedback: feedbackChangedCode[i], message_id: message_id, uid: {{ .Feedback.UserID}}, role: {{ .Feedback.UserRole}}{{if ne .Feedback.Password ""}}, password: {{ .Feedback.Password}}{{end}}  }, function(data, status){
						if (status == "success"){
							alert("Feedback posted successfully!");
							window.location.reload();
						} else {
							alert("Could not post the feedback. Please try again!");
						}
					});
				}
				feedbackCounter ++
			}

			// Run NLP for submission also ?
			
			function sendGrade(i, submission_id, snapshot_id, code, grade) {
				// var feedback = $('#'+submission_id).val().trim();
				// var code = $('#editor-'+submission_id).val();
				var feedbackCode = code;
				var code = submissionChangedCode[i]
				$.post("/teacher_grades", {content: feedbackCode, changed: "", decision: grade, sid: submission_id, uid: {{ .Submission.UserID}}, role: {{ .Submission.UserRole}}{{if ne .Submission.Password ""}}, password: {{ .Submission.Password}}{{end}}  }, function(data, status){
					if (status == "success"){
						// Save and send feedback if the code is changed
						if (submissionChangedCode[i] != "") {
							feedbackCode = submissionChangedCode[i]
							$.post("/save_snapshot_feedback", {snapshot_id: snapshot_id, feedback: feedbackCode, uid: {{ .Submission.UserID}}, role: {{ .Submission.UserRole}}{{if ne .Submission.Password ""}}, password: {{ .Submission.Password}}{{end}} }, function(data1, status1){
							});
						}
						alert("Graded successfully!");
						window.location.reload();
					} else {
						alert("Could not grade the submission. Please try again!");
					}
				});
			}
			
		</script>

	</body>
	</html>
`
