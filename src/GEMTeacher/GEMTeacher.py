# GEMTeacher
# Author: Vinhthuy Phan, 2018
#

import sublime, sublime_plugin
import urllib.parse
import urllib.request
import os
import json
import socket
import webbrowser
import random
import re

gemtAnswerTag = 'ANSWER:'
gemtFOLDER = ''
gemtTIMEOUT = 7
gemtStudentSubmissions = {}
gemtFILE = os.path.join(os.path.dirname(os.path.realpath(__file__)), "info")
gemtSERVER = ''

# ------------------------------------------------------------------
# ------------------------------------------------------------------
# These functionalities are unique to the GEMTeacher module
# ------------------------------------------------------------------
# ------------------------------------------------------------------

# ------------------------------------------------------------------
class gemtStatistics(sublime_plugin.WindowCommand):
	def run(self):
		global gemtSERVER
		if gemtSERVER == '':
			sublime.run_command('gemt_connect')
			if gemtSERVER == '':
				sublime.message_dialog('Unable to connect. Check server address.')
				return
		passcode = gemtRequest('teacher_gets_passcode', {})
		with open(gemtFILE, 'r') as f:
			info = json.loads(f.read())
		data = urllib.parse.urlencode({'pc' : passcode, 'pid': 0})
		webbrowser.open(gemtSERVER + '/statistics?' + data)

# ------------------------------------------------------------------
class gemtReport(sublime_plugin.WindowCommand):
	def run(self):
		global gemtSERVER
		if gemtSERVER == '':
			sublime.run_command('gemt_connect')
			if gemtSERVER == '':
				sublime.message_dialog('Unable to connect. Check server address.')
				return

		passcode = gemtRequest('teacher_gets_passcode', {})
		with open(gemtFILE, 'r') as f:
			info = json.loads(f.read())
		data = urllib.parse.urlencode({'pc' : passcode})
		webbrowser.open(gemtSERVER + '/report?' + data)

# ------------------------------------------------------------------
class gemtViewActivities(sublime_plugin.WindowCommand):
	def run(self):
		global gemtSERVER
		if gemtSERVER == '':
			sublime.run_command('gemt_connect')
			if gemtSERVER == '':
				sublime.message_dialog('Unable to connect. Check server address.')
				return
		passcode = gemtRequest('teacher_gets_passcode', {})
		with open(gemtFILE, 'r') as f:
			info = json.loads(f.read())
		data = urllib.parse.urlencode({'pc' : passcode})
		webbrowser.open(gemtSERVER + '/view_activities?' + data)


# ------------------------------------------------------------------
def gemt_get_problem_info(fname):
	basename = os.path.basename(fname)
	with open(fname, 'r', encoding='utf-8') as fp:
		content = fp.read()
	items = content.split('\n',1)
	if len(items)==0 or (not items[0].startswith('#') and not items[0].startswith('//')):
		return content, '', 0, 0, 0, '', basename, False

	merit, effort, attempts, tag, exact_answer = 0, 0, 0, '', True
	first_line, body = items[0], items[1]
	if first_line.startswith('#'):
		prefix = '#'
		first_line = first_line.strip('# ')
	else:
		prefix = '//'
		first_line = first_line.strip('/ ')
	try:
		items = re.match('(\d+)\s+(\d+)\s+(\d+)(\s+(\w.*))?', first_line).groups()
		merit, effort, attempts, tag = int(items[0]), int(items[1]), int(items[2]), items[4]
		if tag is None:
			tag = ''
		tag = tag.strip()
		if tag.startswith('_manual_'):
			exact_answer = False
			tag = tag.split('_manual_')[1].strip()
	except:
		return content, '', 0, 0, 0, '', basename, False

	if merit < effort:
		return content, '', 0, 0, 0, '', basename, False

	body = '{} {} points, {} for effort. Maximum attempts: {}.\n{}'.format(
		prefix, merit, effort, attempts, body)
	answer = ''
	if os.path.exists(fname + '.answer'):
		with open(fname + '.answer', 'r', encoding='utf-8') as fp:
			answer = fp.read().strip()

	# items = body.split(gemtAnswerTag)
	# if len(items) > 2:
	# 	sublime.message_dialog('This problem has {} answers. There should be at most 1.'.format(len(items)-1))
	# 	raise Exception('Too many answers')

	# body = '{} {} points, {} for effort. Maximum attempts: {}.\n{}'.format(
	# 	prefix, merit, effort, attempts, items[0])
	# answer = ''
	# if len(items) == 2:
	# 	body += '\n{} '.format(gemtAnswerTag)
	# 	answer = items[1].strip()

	return body, answer, merit, effort, attempts, tag, basename, exact_answer


# ------------------------------------------------------------------
class gemtShare(sublime_plugin.TextCommand):
	def run(self, edit):
		fname = self.view.file_name()
		if fname is None:
			sublime.message_dialog('Content must be saved first.')
			return
		content, answer, merit, effort, attempts, tag, name, exact_answer = gemt_get_problem_info(fname)
		data = {
			'content': 		content,
			'answer':		answer,
			'merit':		merit,
			'effort':		effort,
			'attempts':		attempts,
			'tag':			tag,
			'filename':		name,
			'exact_answer':	exact_answer,
		}
		response = gemtRequest('teacher_broadcasts', data)
		if response is not None:
			if merit==0:
				mesg = 'Sharing not-graded content. '
			else:
				mesg = 'Sharing graded content. '
			mesg += response
			sublime.message_dialog(mesg)

# ------------------------------------------------------------------
class gemtDeactivateProblems(sublime_plugin.TextCommand):
	def run(self, edit):
		if self.view.file_name() is None:
			sublime.message_dialog('Unknown problem!')
			return
		filename = os.path.basename(self.view.file_name())
		if sublime.ok_cancel_dialog('Submission for this problem will no longer be possible. Click OK to confirm.'):
			response = gemtRequest('teacher_deactivates_problems', {'filename':filename})
			if response == '-1':
				sublime.message_dialog('Unknown or inactive problem!')
			elif response == '0':
				sublime.message_dialog('Problem is now inactive.')
			elif response == '1':
				global gemtSERVER
				with open(gemtFILE, 'r') as f:
					info = json.loads(f.read())
				passcode = gemtRequest('teacher_gets_passcode', {})
				p = urllib.parse.urlencode({'pc' : passcode, 'filename':filename})
				webbrowser.open(gemtSERVER + '/view_answers?' + p)

# ------------------------------------------------------------------
class gemtClearSubmissions(sublime_plugin.ApplicationCommand):
	def run(self):
		if sublime.ok_cancel_dialog('Do you want to clear all submissions and white boards?'):
			response = gemtRequest('teacher_clears_submissions', {})
			sublime.message_dialog(response)
# ----------------------------------------------------------------------


# ----------------------------------------------------------------------
# ----------------------------------------------------------------------
# These functionalities below are identical to the GEMAssistant module
# ----------------------------------------------------------------------
# ----------------------------------------------------------------------
def gemtRequest(path, data, authenticated=True, method='POST'):
	global gemtFOLDER, gemtSERVER

	try:
		with open(gemtFILE, 'r') as f:
			info = json.loads(f.read())
	except:
		info = dict()

	if 'Folder' not in info:
		sublime.message_dialog("Please set a local folder for keeping working files.")
		return None

	if 'CourseId' not in info:
		sublime.message_dialog("Please set the course id.")
		return None

	if 'Server' not in info:
		sublime.message_dialog("Please set server address.")
		return None

	if gemtSERVER == '':
		sublime.run_command('gemt_connect')
		if gemtSERVER == '':
			sublime.message_dialog('Unable to connect. Check server address or course id.')
			return

	if authenticated:
		if 'Name' not in info or 'Password' not in info:
			sublime.message_dialog("Please ask to setup a new teacher account and register.")
			return None
		data['name'] = info['Name']
		data['password'] = info['Password']
		data['uid'] = info['Uid']
		data['role'] = 'teacher'
		gemtFOLDER = info['Folder']

	url = urllib.parse.urljoin(gemtSERVER, path)
	load = urllib.parse.urlencode(data).encode('utf-8')
	req = urllib.request.Request(url, load, method=method)
	try:
		with urllib.request.urlopen(req, None, gemtTIMEOUT) as response:
			return response.read().decode(encoding="utf-8")
	except urllib.error.HTTPError as err:
		sublime.message_dialog("{0}".format(err))
	except urllib.error.URLError as err:
		sublime.message_dialog("{0}\nCannot connect to server.".format(err))
	print('Something is wrong')
	return None

# ------------------------------------------------------------------
class gemtViewBulletinBoard(sublime_plugin.ApplicationCommand):
	def run(self):
		response = gemtRequest('teacher_gets_passcode', {})
		if response.startswith('Unauthorized'):
			sublime.message_dialog('Unauthorized')
		else:
			global gemtSERVER
			p = urllib.parse.urlencode({'pc' : response})
			with open(gemtFILE, 'r') as f:
				info = json.loads(f.read())
			webbrowser.open(gemtSERVER + '/view_bulletin_board?' + p)

# ------------------------------------------------------------------
class gemtAddBulletin(sublime_plugin.TextCommand):
	def run(self, edit):
		this_file_name = self.view.file_name()
		if this_file_name is None:
			sublime.message_dialog('Error: file is empty.')
			return
		beg, end = self.view.sel()[0].begin(), self.view.sel()[0].end()
		content = '\n\n' + self.view.substr(sublime.Region(beg,end)) + '\n\n'
		if len(content) <= 20:
			sublime.message_dialog('Select more text to show on the bulletin board.')
			return
		response = gemtRequest('teacher_adds_bulletin_page', {'content':content})
		if response:
			sublime.message_dialog(response)

# ------------------------------------------------------------------
class gemtPutBack(sublime_plugin.TextCommand):
	def run(self, edit):
		fname = self.view.file_name()
		sid = os.path.basename(os.path.dirname(fname))
		response = gemtRequest('teacher_puts_back', {'sid':sid})
		if response:
			sublime.message_dialog(response)
		if response != 'Unknown submission.':
			self.view.window().run_command('close')

# ------------------------------------------------------------------
def remove_first_line(content):
	lines = content.split('\n')
	if len(lines) > 0:
		return '\n'.join(lines[1:])
	else:
		return content

# ------------------------------------------------------------------
def gemt_grade(self, edit, decision):
	fname = self.view.file_name()
	changed = False
	sid = os.path.basename(os.path.dirname(fname))
	if decision=='dismissed':
		content = ''
	else:
		content = self.view.substr(sublime.Region(0, self.view.size())).strip()
		if sid in gemtStudentSubmissions:
			v1 = gemtStudentSubmissions[sid].strip()
			v2 = content.strip()
			if remove_first_line(v1) != remove_first_line(v2):
				changed = True

	stop = False
	if sid == '0':
		stop = True
	try:
		int(sid)
	except:
		stop = True
	if stop:
		if decision=='dismissed':
			self.view.window().run_command('close')
		else:
			sublime.message_dialog('This is not a graded problem.')
		return

	data = dict(
		sid = sid,
		content = content,
		decision = decision,
		changed = changed,
	)
	response = gemtRequest('teacher_grades', data)
	if response:
		sublime.message_dialog(response)
		self.view.window().run_command('close')

# ------------------------------------------------------------------
class gemtUngrade(sublime_plugin.TextCommand):
	def run(self, edit):
		gemt_grade(self, edit, "ungraded")

class gemtGradeCorrect(sublime_plugin.TextCommand):
	def run(self, edit):
		gemt_grade(self, edit, "correct")

class gemtGradeIncorrect(sublime_plugin.TextCommand):
	def run(self, edit):
		gemt_grade(self, edit, "incorrect")

class gemtDismissed(sublime_plugin.TextCommand):
	def run(self, edit):
		gemt_grade(self, edit, "dismissed")

# ------------------------------------------------------------------
def gemt_gets(self, index, priority):
	global gemtStudentSubmissions
	response = gemtRequest('teacher_gets', {'index':index, 'priority':priority})
	if response is not None:
		sub = json.loads(response)
		if sub['Content'] != '':
			filename = sub['Filename']
			sid = str(sub['Sid'])
			dir = os.path.join(gemtFOLDER, sid)
			if not os.path.exists(dir):
				os.mkdir(dir)
			local_file = os.path.join(dir, filename)
			with open(local_file, 'w', encoding='utf-8') as fp:
				fp.write(sub['Content'])
			gemtStudentSubmissions[sid] = sub['Content']
			if sublime.active_window().id() == 0:
				sublime.run_command('new_window')
			sublime.active_window().open_file(local_file)
			if sub['Priority'] == 2:
				sublime.message_dialog('This student asked for help.')
		elif priority == 0:
			sublime.message_dialog('There are no submissions.')
		elif priority > 0:
			sublime.message_dialog('There are no submissions with priority {}.'.format(priority))
		elif index >= 0:
			sublime.message_dialog('There are no submission with index {}.'.format(index))
		
# ------------------------------------------------------------------
class gemtSeeQueue(sublime_plugin.ApplicationCommand):
	def run(self):
		response = gemtRequest('teacher_gets_queue', {})
		if response is not None:
			json_obj = json.loads(response)
			if json_obj is None:
				sublime.status_message("Queue is empty.")
			else:
				users = []
				for entry in json_obj:
					if entry['Priority'] == 2:
						status = '😥'
					elif entry['Priority'] == 1:
						status = '😎'
					else:
						status = ''
					users.append( '{} {}'.format(entry['Name'], status))
				if users:
					sublime.active_window().active_view().show_popup_menu(users, self.request_entry)
				else:
					sublime.status_message("Queue is empty.")

	# ---------------------------------------------------------
	def request_entry(self, selected):
		if selected < 0:
			return
		gemt_gets(self, selected, 1)

# ------------------------------------------------------------------
# Priorities: 1 (I got it), 2 (I need help),
# ------------------------------------------------------------------
class gemtGetPrioritized(sublime_plugin.ApplicationCommand):
	def run(self):
		gemt_gets(self, -1, 0)

class gemtGetFromNeedHelp(sublime_plugin.ApplicationCommand):
	def run(self):
		gemt_gets(self, -1, 2)

class gemtGetFromOk(sublime_plugin.ApplicationCommand):
	def run(self):
		gemt_gets(self, -1, 1)

# ------------------------------------------------------------------
class gemtSetLocalFolder(sublime_plugin.ApplicationCommand):
	def run(self):
		try:
			with open(gemtFILE, 'r') as f:
				info = json.loads(f.read())
		except:
			info = dict()
		if 'Folder' not in info:
			info['Folder'] = os.path.join(os.path.expanduser('~'), 'GEMT')
		if sublime.active_window().id() == 0:
			sublime.run_command('new_window')
		sublime.active_window().show_input_panel("Specify a folder on your computer to store working files.",
			info['Folder'],
			self.set,
			None,
			None)

	def set(self, folder):
		folder = folder.strip()
		if len(folder) > 0:
			try:
				with open(gemtFILE, 'r') as f:
					info = json.loads(f.read())
			except:
				info = dict()
			info['Folder'] = folder
			if not os.path.exists(folder):
				try:
					os.mkdir(folder)
					with open(gemtFILE, 'w') as f:
						f.write(json.dumps(info, indent=4))
				except:
					sublime.message_dialog('Could not create {}.'.format(folder))
			else:
				with open(gemtFILE, 'w') as f:
					f.write(json.dumps(info, indent=4))
				sublime.message_dialog('Folder exists. Will use it to store working files.')
		else:
			sublime.message_dialog("Folder name cannot be empty.")

# ------------------------------------------------------------------
class gemtConnect(sublime_plugin.ApplicationCommand):
	def run(self):
		try:
			with open(gemtFILE, 'r') as f:
				info = json.loads(f.read())
		except:
			info = dict()

		if 'CourseId' not in info:
			sublime.message_dialog("Please set the course id.")
			return None

		if 'Server' not in info:
			sublime.message_dialog("Please set server address.")
			return None

		global gemtSERVER
		url = urllib.parse.urljoin(info['Server'], 'ask')
		load = urllib.parse.urlencode({'who':info['CourseId']}).encode('utf-8')
		req = urllib.request.Request(url, load)
		try:
			with urllib.request.urlopen(req, None, gemtTIMEOUT) as response:
				server = response.read().decode(encoding="utf-8")
				try:
					with open(gemtFILE, 'r') as f:
						info = json.loads(f.read())
				except:
					info = dict()
				if not server.startswith('http://'):
					sublime.message_dialog('Unable to get address.')
					return
				gemtSERVER = server
				sublime.status_message('Connected')
		except urllib.error.HTTPError as err:
			sublime.message_dialog("{0}".format(err))
		except urllib.error.URLError as err:
			sublime.message_dialog("{0}\nCannot connect to server.".format(err))

# ------------------------------------------------------------------
class gemtCompleteRegistration(sublime_plugin.ApplicationCommand):
	def run(self):
		try:
			with open(gemtFILE, 'r') as f:
				info = json.loads(f.read())
		except:
			info = dict()

		if 'CourseId' not in info:
			sublime.message_dialog('Please enter course id.')
			return

		if 'Name' not in info:
			sublime.message_dialog('Please enter assigned username.')
			return

		response = gemtRequest(
			'complete_registration',
			{'role':'teacher', 'name':info['Name'], 'course_id':info['CourseId']},
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
		with open(gemtFILE, 'w') as f:
			f.write(json.dumps(info, indent=4))

# ------------------------------------------------------------------
class gemtSetServerAddress(sublime_plugin.ApplicationCommand):
	def run(self):
		try:
			with open(gemtFILE, 'r') as f:
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
				with open(gemtFILE, 'r') as f:
					info = json.loads(f.read())
			except:
				info = dict()
			if not addr.startswith('http://'):
				addr = 'http://' + addr
			info['Server'] = addr
			with open(gemtFILE, 'w') as f:
				f.write(json.dumps(info, indent=4))
			sublime.message_dialog('Server address is set to ' + addr)
		else:
			sublime.message_dialog("Server address cannot be empty.")

# ------------------------------------------------------------------
class gemtSetCourseId(sublime_plugin.ApplicationCommand):
	def run(self):
		try:
			with open(gemtFILE, 'r') as f:
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
				with open(gemtFILE, 'r') as f:
					info = json.loads(f.read())
			except:
				info = dict()
			info['CourseId'] = cid
			with open(gemtFILE, 'w') as f:
				f.write(json.dumps(info, indent=4))
			sublime.message_dialog('Course id is set to ' + cid)
		else:
			sublime.message_dialog("Course id cannot be empty.")

# ------------------------------------------------------------------
class gemtSetName(sublime_plugin.ApplicationCommand):
	def run(self):
		try:
			with open(gemtFILE, 'r') as f:
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
				with open(gemtFILE, 'r') as f:
					info = json.loads(f.read())
			except:
				info = dict()
			info['Name'] = name
			with open(gemtFILE, 'w') as f:
				f.write(json.dumps(info, indent=4))
			sublime.message_dialog('Assigned name is set to ' + name)
		else:
			sublime.message_dialog("Name cannot be empty.")

# ------------------------------------------------------------------
class gemtUpdate(sublime_plugin.WindowCommand):
	def run(self):
		package_path = os.path.join(sublime.packages_path(), "GEMTeacher");
		try:
			version = open(os.path.join(package_path, "VERSION")).read()
		except:
			version = 0
		if sublime.ok_cancel_dialog("Current version is {}. Click OK to update.".format(version)):
			if not os.path.isdir(package_path):
				os.mkdir(package_path)
			module_file = os.path.join(package_path, "GEMTeacher.py")
			menu_file = os.path.join(package_path, "Main.sublime-menu")
			keymap_file = os.path.join(package_path, "Default.sublime-keymap")
			version_file = os.path.join(package_path, "version.go")
			urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/GEM/master/src/GEMTeacher/GEMTeacher.py", module_file)
			urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/GEM/master/src/GEMTeacher/Main.sublime-menu", menu_file)
			urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/GEM/master/src/GEMTeacher/Default.sublime-keymap", keymap_file)
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
			sublime.message_dialog("GEM has been updated to version %s." % version)

# ------------------------------------------------------------------


