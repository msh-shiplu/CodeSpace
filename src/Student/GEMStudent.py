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
import datetime
import webbrowser

gemsUpdateTimeout = 10000
gemsFILE = os.path.join(os.path.dirname(os.path.realpath(__file__)), "info")
gemsFOLDER = ''
gemsTIMEOUT = 7
gemsAnswerTag = 'ANSWER:'
gemsTracking = False
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
def gems_problem_info(fname):
	basename = os.path.basename(fname)
	if '.' in fname:
		prefix, ext = fname.rsplit('.',1)
	else:
		prefix, ext = fname, ''
	if prefix.count('_') < 1:
		return basename, 0
	prefix, pid = prefix.rsplit('_', 1)
	try:
		pid = int(pid)
	except:
		return basename, 0
	if ext == '':
		orginal_fname = prefix
	else:
		orginal_fname = prefix + '.' + ext
	return orginal_fname, pid

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
		if board_stat == 1:
			if sublime.active_window().id() == 0:
				sublime.run_command('new_window')
			sublime.active_window().run_command('gems_get_board_content')

		# Keep checking periodically
		update_timeout = gemsUpdateTimeout
		if submission_stat == 1:
			update_timeout = gemsUpdateTimeout // 2
		print('checking', submission_stat, board_stat, update_timeout)
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
	original_fname, pid = gems_problem_info(fname)
	if pid == 0:
		priority = 1
	if pid > 0 or sublime.ok_cancel_dialog('This file is not a graded problem. Do you want to send it?'):
		content = self.view.substr(sublime.Region(0, self.view.size())).lstrip()
		items = content.rsplit(gemsAnswerTag, 1)
		if len(items)==2:
			answer = items[1].strip()
		else:
			answer = ''
		data = dict(
			content=content,
			answer=answer,
			pid=pid,
			filename=original_fname,
			priority=priority,
		)
		response = gemsRequest('student_shares', data)
		sublime.message_dialog(response)
		if pid > 0 and gemsTracking==False:
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
		for board in json_obj:
			content = board['Content']
			attempts = board['Attempts']
			filename = board['Filename']
			pid = board['Pid']
			today = datetime.datetime.today()
			if '.' in filename:
				fname, ext = filename.rsplit('.',1)
			else:
				fname, ext = filename, ''
			new_fname = os.path.join(gemsFOLDER, '{}_{}.{}'.format(fname,pid,ext))
			if pid>0 and os.path.exists(new_fname):
				new_fname = os.path.join(gemsFOLDER, 'FEEDBACK.txt')
			with open(new_fname, 'w', encoding='utf-8') as f:
				f.write(content)
			if sublime.active_window().id() == 0:
				sublime.run_command('new_window')
			sublime.active_window().open_file(new_fname)

# ------------------------------------------------------------------------------
# ------------------------------------------------------------------------------
# These functionalities below are identical to those of teachers
# ------------------------------------------------------------------------------
# ------------------------------------------------------------------------------
def gemsRequest(path, data, authenticated=True, method='POST', verbal=True):
	global gemsFOLDER
	try:
		with open(gemsFILE, 'r') as f:
			info = json.loads(f.read())
	except:
		info = dict()

	if 'Folder' not in info:
		if verbal:
			sublime.message_dialog("Please set a local folder to store working files.")
		return None

	if 'Server' not in info:
		if verbal:
			sublime.message_dialog("Please connect to the server first.")
		return None

	if authenticated:
		if 'Uid' not in info:
			sublime.message_dialog("Please register.")
			return None
		data['name'] = info['Name']
		data['password'] = info['Password']
		data['uid'] = info['Uid']
		gemsFOLDER = info['Folder']

	url = urllib.parse.urljoin(info['Server'], path)
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
		if 'NameServer' not in info:
			self.ask_to_set_server_address(info)
		else:
			self.set_server_address_via_nameserver(info)

	# ------------------------------------------------------------------
	def set_server_address_via_nameserver(self, info):
		url = urllib.parse.urljoin(info['NameServer'], 'ask')
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
					sublime.message_dialog(server)
					return
				info['Server'] = server
				with open(gemsFILE, 'w') as f:
					f.write(json.dumps(info, indent=4))
				sublime.message_dialog('Connected to server at {}'.format(server))
		except urllib.error.HTTPError as err:
			sublime.message_dialog("{0}".format(err))
		except urllib.error.URLError as err:
			sublime.message_dialog("{0}\nCannot connect to name server.".format(err))

	# ------------------------------------------------------------------
	def ask_to_set_server_address(self, info):
		try:
			with open(gemsFILE, 'r') as f:
				info = json.loads(f.read())
		except:
			info = dict()
		if 'Server' not in info:
			info['Server'] = ''
		if sublime.active_window().id() == 0:
			sublime.run_command('new_window')
		sublime.active_window().show_input_panel("Set server address:",
			info['Server'],
			self.set_server_address,
			None,
			None)

	# ------------------------------------------------------------------
	def set_server_address(self, addr):
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
class gemsCompleteRegistration(sublime_plugin.ApplicationCommand):
	def run(self):
		try:
			with open(gemsFILE, 'r') as f:
				info = json.loads(f.read())
		except:
			info = dict()

		if 'Folder' not in info:
			sublime.message_dialog("Please set a local folder for keeping working files.")
			return None

		if 'Server' not in info:
			sublime.message_dialog("Please set server address.")
			return None

		mesg = 'Enter assigned_id'
		if 'Name' in info:
			mesg = '{} is already registered. Enter assigned_id:'.format(info['Name'])

		if 'Name' not in info:
			info['Name'] = ''

		if sublime.active_window().id() == 0:
			sublime.run_command('new_window')
		sublime.active_window().show_input_panel(mesg,info['Name'],self.process,None,None)

	# ------------------------------------------------------------------
	def process(self, data):
		name = data.strip()
		response = gemsRequest(
			'complete_registration',
			{'name':name.strip(), 'role':'student'},
			authenticated=False,
		)
		if response == 'Failed':
			sublime.message_dialog('Failed to complete registration.')
			return

		name_server = ""
		if response.count(',') == 3:
			uid, password, course_id, name_server = response.split(',')
			if not name_server.strip().startswith('http://'):
				name_server = 'http://' + name_server.strip()
		elif response.count(',') == 2:
			uid, password, course_id = response.split(',')
		else:
			sublime.message_dialog('Unable to complete registration.')
			return
		try:
			with open(gemsFILE, 'r') as f:
				info = json.loads(f.read())
		except:
			info = dict()
		info['Uid'] = int(uid)
		info['Password'] = password.strip()
		info['Name'] = name.strip()
		info['CourseId'] = course_id.strip()
		if name_server != "":
			info['NameServer'] = name_server
		with open(gemsFILE, 'w') as f:
			f.write(json.dumps(info, indent=4))
		sublime.message_dialog('{} registered for {}'.format(name, course_id))

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
