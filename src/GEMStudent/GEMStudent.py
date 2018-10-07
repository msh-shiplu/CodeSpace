# GEMStudent
# Author: Vinhthuy Phan, 2018
#
import sublime, sublime_plugin
import urllib.parse
import urllib.request
import os
import json
import time
import random
import shutil
import datetime
import webbrowser

gemsUpdateIntervalLong  = 20000		# Update interval
gemsUpdateIntervalShort = 10000		# When submission is being looked at
gemsFILE = os.path.join(os.path.dirname(os.path.realpath(__file__)), "info")
gemsFOLDER = ''
gemsTIMEOUT = 7
gemsAnswerTag = 'ANSWER:'
gemsTracking = False
gemsSERVER = ''
gemsSERVER_TIME = 0
gemsConnected = False
gemsUpdateMessage = {
	1 : "Your submission is being looked at.",
	2 : "Teacher did not grade your submission.",
	3 : "Good effort!!!  However, the teacher did not think your solution was correct.",
	4 : "Your solution was correct.",
}

# ------------------------------------------------------------------
class gemsAttendanceReport(sublime_plugin.ApplicationCommand):
	def run(self):
		response = gemsRequest('student_checks_in', {})
		json_obj = json.loads(response)
		dates = set()
		for d in json_obj:
			dates.add(datetime.datetime.fromtimestamp(d).strftime('%Y-%m-%d'))

		with open(gemsFILE, 'r') as f:
			info = json.loads(f.read())
		report_file = os.path.join(info['Folder'], 'Attendance.txt')
		with open(report_file, 'w', encoding='utf-8') as f:
			f.write('Your attendance was taken on these dates:\n')
			for d in sorted(dates, reverse=True):
				f.write('{}\n'.format(d))
		sublime.run_command('new_window')
		sublime.active_window().open_file(report_file)

# ------------------------------------------------------------------
class gemsPointsReport(sublime_plugin.ApplicationCommand):
	def run(self):
		response = gemsRequest('student_gets_report', {})
		if response != None:
			json_obj = json.loads(response)
			report = {}
			total_points = 0
			for i in json_obj:
				Date, Points, Filename = i['Date'], i['Points'], i['Filename']
				if Date not in report:
					report[Date] = []
				if '.' in Filename:
					prefix, ext = Filename.rsplit('.',1)
				else:
					prefix = Filename
				report[Date].append((Points,prefix))
				total_points += Points

			with open(gemsFILE, 'r') as f:
				info = json.loads(f.read())
			report_file = os.path.join(info['Folder'], 'Points.txt')
			with open(report_file, 'w', encoding='utf-8') as f:
				f.write('Total points: {}\n'.format(total_points))
				for d,v in reversed(sorted(report.items())):
					date = datetime.datetime.fromtimestamp(d).strftime('%Y-%m-%d')
					for entry in v:
						f.write('{}\t{}\t{}\n'.format(date,entry[0],entry[1]))

			sublime.run_command('new_window')
			sublime.active_window().open_file(report_file)

# ------------------------------------------------------------------
def gems_periodic_update():
	global gemsTracking
	response = gemsRequest('student_periodic_update', {}, verbal=False)
	if response is None:
		print('Response is None. Stop tracking.')
		gemsTracking = False
		return
	try:
		submission_stat, board_stat = response.split(';')
		submission_stat = int(submission_stat)
		board_stat = int(board_stat)

		# Display messages if necessary
		mesg = ""
		if submission_stat > 0 and submission_stat in gemsUpdateMessage:
			mesg = gemsUpdateMessage[submission_stat]
		if board_stat == 1:
			mesg += "\nTeacher placed new material on your board."
		mesg = mesg.strip()
		if mesg != "":
			sublime.message_dialog(mesg)

		# Open board pages and feedback automatically
		# if board_stat == 1:
		# 	if sublime.active_window().id() == 0:
		# 		sublime.run_command('new_window')
		# 	sublime.active_window().run_command('gems_get_board_content')

		# Keep checking periodically
		update_timeout = gemsUpdateIntervalLong
		if submission_stat == 1:
			update_timeout = gemsUpdateIntervalShort
		# print('checking', submission_stat, board_stat, update_timeout)
		sublime.set_timeout_async(gems_periodic_update, update_timeout)
	except:
		gemsTracking = False

# ------------------------------------------------------------------
def gems_share(self, edit, priority):
	global gemsTracking
	fname = self.view.file_name()
	if fname is None:
		sublime.message_dialog('Cannot share unsaved content.')
		return
	content = self.view.substr(sublime.Region(0, self.view.size())).lstrip()
	items = content.rsplit(gemsAnswerTag, 1)
	if len(items)==2:
		answer = items[1].strip()
	else:
		answer = ''
	data = dict(
		content=content,
		answer=answer,
		filename=os.path.basename(fname),
		priority=priority,
	)
	response = gemsRequest('student_shares', data)
	sublime.message_dialog(response)
	if gemsTracking==False:
		gemsTracking = True
		sublime.set_timeout_async(gems_periodic_update, 5000)

# ------------------------------------------------------------------
class gemsNeedHelp(sublime_plugin.TextCommand):
	def run(self, edit):
		gems_share(self, edit, priority=2)

# ------------------------------------------------------------------
class gemsGotIt(sublime_plugin.TextCommand):
	def run(self, edit):
		gems_share(self, edit, priority=1)

# ------------------------------------------------------------------
class gemsGetBoardContent(sublime_plugin.ApplicationCommand):
	def run(self):
		response = gemsRequest('student_gets', {})
		if response is None:
			return
		json_obj = json.loads(response)
		if json_obj == []:
			sublime.message_dialog("Whiteboard is empty.")
			return

		feedback_dir = os.path.join(gemsFOLDER, 'FEEDBACK')
		if not os.path.exists(feedback_dir):
			os.mkdir(feedback_dir)
		old_dir = os.path.join(gemsFOLDER, 'OLD')
		if not os.path.exists(old_dir):
			os.mkdir(old_dir)

		for board in json_obj:
			content = board['Content']
			filename = board['Filename']
			mesg = ''
			if board['Type'] == 'feedback':
				local_file = os.path.join(feedback_dir, filename)
				mesg = 'Teacher has some feedback for you.'
			else:
				local_file = os.path.join(gemsFOLDER, filename)
				if os.path.exists(local_file):
					with open(local_file) as f:
						moved_file = os.path.join(old_dir, filename)
						with open(moved_file, 'w', encoding='utf-8') as newf:
							newf.write(f.read())
					mesg = 'Move existing file, {}, to {}.'.format(filename,old_dir)
			with open(local_file, 'w', encoding='utf-8') as f:
				f.write(content)
			if sublime.active_window().id() == 0:
				sublime.run_command('new_window')
			sublime.active_window().open_file(local_file)
			if mesg != '':
				sublime.message_dialog(mesg)

# ------------------------------------------------------------------------------
# ------------------------------------------------------------------------------
# These functionalities below are identical to those of teachers
# ------------------------------------------------------------------------------
# ------------------------------------------------------------------------------
def gemsRequest(path, data, authenticated=True, method='POST', verbal=True):
	global gemsFOLDER, gemsSERVER, gemsSERVER_TIME

	try:
		with open(gemsFILE, 'r') as f:
			info = json.loads(f.read())
	except:
		info = dict()

	if 'Folder' not in info:
		if verbal:
			sublime.message_dialog("Please set a local folder to store working files.")
		return None

	if 'CourseId' not in info:
		sublime.message_dialog("Please set the course id.")
		return None

	if 'Server' not in info:
		if verbal:
			sublime.message_dialog("Please connect to the server first.")
		return None

	if gemsSERVER == '' or time.time() - gemsSERVER_TIME > 5400:
		sublime.run_command('gems_connect')
		if gemsSERVER == '':
			sublime.message_dialog('Unable to connect. Check server address or course id.')
			return

	if authenticated:
		if 'Uid' not in info:
			sublime.message_dialog("Please register.")
			return None
		data['name'] = info['Name']
		data['password'] = info['Password']
		data['uid'] = info['Uid']
		gemsFOLDER = info['Folder']

	url = urllib.parse.urljoin(gemsSERVER, path)
	load = urllib.parse.urlencode(data).encode('utf-8')
	req = urllib.request.Request(url, load, method=method)
	try:
		with urllib.request.urlopen(req, None, gemsTIMEOUT) as response:
			return response.read().decode(encoding="utf-8")
	except urllib.error.HTTPError as err:
		if verbal:
			sublime.message_dialog("{0}".format(err))
	except urllib.error.URLError as err:
		if verbal:
			sublime.message_dialog("{0}\nCannot connect to server.".format(err))
	print('Error making request')
	return None

# ------------------------------------------------------------------
class gemsSetLocalFolder(sublime_plugin.ApplicationCommand):
	def run(self):
		try:
			with open(gemsFILE, 'r') as f:
				info = json.loads(f.read())
		except:
			info = dict()
		if 'Folder' not in info:
			info['Folder'] = os.path.join(os.path.expanduser('~'), 'GEM')
		if sublime.active_window().id() == 0:
			sublime.run_command('new_window')
		sublime.active_window().show_input_panel("This folder will be used to store working files.",
			info['Folder'],
			self.set,
			None,
			None)

	def set(self, folder):
		folder = folder.strip()
		if len(folder) > 0:
			try:
				with open(gemsFILE, 'r') as f:
					info = json.loads(f.read())
			except:
				info = dict()
			info['Folder'] = folder
			if not os.path.exists(folder):
				try:
					os.mkdir(folder)
					os.mkdir(os.path.join(folder,'FEEDBACK'))
					with open(gemsFILE, 'w') as f:
						f.write(json.dumps(info, indent=4))
				except:
					sublime.message_dialog('Could not create {}.'.format(folder))
			else:
				with open(gemsFILE, 'w') as f:
					f.write(json.dumps(info, indent=4))
				sublime.message_dialog('Folder exists. Will use it to store working files.')
		else:
			sublime.message_dialog("Folder name cannot be empty.")

# ------------------------------------------------------------------
class gemsConnect(sublime_plugin.ApplicationCommand):
	def run(self):
		try:
			with open(gemsFILE, 'r') as f:
				info = json.loads(f.read())
		except:
			info = dict()

		if 'CourseId' not in info:
			sublime.message_dialog("Please set the course id.")
			return None

		if 'Server' not in info:
			sublime.message_dialog("Please set server address.")
			return None

		global gemsSERVER, gemsSERVER_TIME
		url = urllib.parse.urljoin(info['Server'], 'ask')
		load = urllib.parse.urlencode({'who':info['CourseId']}).encode('utf-8')
		req = urllib.request.Request(url, load)
		try:
			with urllib.request.urlopen(req, None, gemsTIMEOUT) as response:
				server = response.read().decode(encoding="utf-8")
				try:
					with open(gemsFILE, 'r') as f:
						info = json.loads(f.read())
				except:
					info = dict()
				if not server.startswith('http://'):
					sublime.message_dialog('Unable to get address.')
					return
				gemsSERVER = server
				gemsSERVER_TIME = time.time()
				sublime.status_message('Connected')
		except urllib.error.HTTPError as err:
			sublime.message_dialog("{0}".format(err))
		except urllib.error.URLError as err:
			sublime.message_dialog("{0}\nCannot connect to server.".format(err))

# ------------------------------------------------------------------
class gemsCompleteRegistration(sublime_plugin.ApplicationCommand):
	def run(self):
		try:
			with open(gemsFILE, 'r') as f:
				info = json.loads(f.read())
		except:
			info = dict()

		if 'CourseId' not in info:
			sublime.message_dialog('Please enter course id.')
			return

		if 'Name' not in info:
			sublime.message_dialog('Please enter assigned username.')
			return

		response = gemsRequest(
			'complete_registration',
			{'role':'student', 'name':info['Name'], 'course_id':info['CourseId']},
			authenticated=False,
		)
		if response is None:
			sublime.message_dialog('Response is None. Failed to complete registration.')
			return

		if response == 'Failed' or response.count(',') != 1:
			sublime.message_dialog('Failed to complete registration.')
		else:
			uid, password = response.split(',')
			info['Uid'] = int(uid)
			info['Password'] = password.strip()
			sublime.message_dialog('{} is registered.'.format(info['Name']))
		with open(gemsFILE, 'w') as f:
			f.write(json.dumps(info, indent=4))

# ------------------------------------------------------------------
class gemsSetServerAddress(sublime_plugin.ApplicationCommand):
	def run(self):
		try:
			with open(gemsFILE, 'r') as f:
				info = json.loads(f.read())
		except:
			info = dict()

		if 'Server' not in info:
			info['Server'] = ''
		if sublime.active_window().id() == 0:
			sublime.run_command('new_window')
		sublime.active_window().show_input_panel("Set server address.  Press Enter:",
			info['Server'],
			self.set,
			None,
			None)

	def set(self, addr):
		addr = addr.strip()
		if len(addr) > 0:
			try:
				with open(gemsFILE, 'r') as f:
					info = json.loads(f.read())
			except:
				info = dict()
			if not addr.startswith('http://'):
				addr = 'http://' + addr
			info['Server'] = addr
			with open(gemsFILE, 'w') as f:
				f.write(json.dumps(info, indent=4))
		else:
			sublime.message_dialog("Server address cannot be empty.")

# ------------------------------------------------------------------
class gemsSetCourseId(sublime_plugin.ApplicationCommand):
	def run(self):
		try:
			with open(gemsFILE, 'r') as f:
				info = json.loads(f.read())
		except:
			info = dict()

		if 'CourseId' not in info:
			info['CourseId'] = ''
		if sublime.active_window().id() == 0:
			sublime.run_command('new_window')
		sublime.active_window().show_input_panel("Set course id.  Press Enter:",
			info['CourseId'],
			self.set,
			None,
			None)

	def set(self, cid):
		cid = cid.strip()
		if len(cid) > 0:
			try:
				with open(gemsFILE, 'r') as f:
					info = json.loads(f.read())
			except:
				info = dict()
			info['CourseId'] = cid
			with open(gemsFILE, 'w') as f:
				f.write(json.dumps(info, indent=4))
			sublime.message_dialog('Course id is set to ' + cid)
		else:
			sublime.message_dialog("Server address cannot be empty.")

# ------------------------------------------------------------------
class gemsSetName(sublime_plugin.ApplicationCommand):
	def run(self):
		try:
			with open(gemsFILE, 'r') as f:
				info = json.loads(f.read())
		except:
			info = dict()
		if 'Name' not in info:
			info['Name'] = ''
		if sublime.active_window().id() == 0:
			sublime.run_command('new_window')
		sublime.active_window().show_input_panel("Set assgined username.  Press Enter:",
			info['Name'],
			self.set,
			None,
			None)

	def set(self, name):
		name = name.strip()
		if len(name) > 0:
			try:
				with open(gemsFILE, 'r') as f:
					info = json.loads(f.read())
			except:
				info = dict()
			info['Name'] = name
			with open(gemsFILE, 'w') as f:
				f.write(json.dumps(info, indent=4))
			sublime.message_dialog('Assigned name is set to ' + name)
		else:
			sublime.message_dialog("Name cannot be empty.")

# ------------------------------------------------------------------
class gemsUpdate(sublime_plugin.WindowCommand):
	def run(self):
		package_path = os.path.join(sublime.packages_path(), "GEMStudent");
		try:
			version = open(os.path.join(package_path, "VERSION")).read()
		except:
			version = 0
		if sublime.ok_cancel_dialog("Current version is {}. Click OK to update.".format(version)):
			if not os.path.isdir(package_path):
				os.mkdir(package_path)
			module_file = os.path.join(package_path, "GEMStudent.py")
			menu_file = os.path.join(package_path, "Main.sublime-menu")
			version_file = os.path.join(package_path, "version.go")
			urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/GEM/master/src/GEMStudent/GEMStudent.py", module_file)
			urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/GEM/master/src/GEMStudent/Main.sublime-menu", menu_file)
			urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/GEM/master/src/version.go", version_file)
			with open(version_file) as f:
				lines = f.readlines()
			for line in lines:
				if line.strip().startswith('const VERSION ='):
					prefix, version = line.strip().split('const VERSION =')
					version = version.strip().strip('"')
					break
			os.remove(version_file)
			with open(os.path.join(package_path, "VERSION"), 'w') as f:
				f.write(version)
			sublime.message_dialog("GEM has been updated to version {}.".format(version))

# ------------------------------------------------------------------
