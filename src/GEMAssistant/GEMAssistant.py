# GEMAssistant
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

gemaAnswerTag = 'ANSWER:'
gemaFOLDER = ''
gemaTIMEOUT = 7
gemaStudentSubmissions = {}
gemaConnected = False
gemaFILE = os.path.join(os.path.dirname(os.path.realpath(__file__)), "info")
gemaSERVER = ''

# ----------------------------------------------------------------------
# ----------------------------------------------------------------------
# These functionalities below are identical to the GEMAssistant module
# ----------------------------------------------------------------------
# ----------------------------------------------------------------------
def gemaRequest(path, data, authenticated=True, method='POST'):
	global gemaFOLDER, gemaSERVER

	try:
		with open(gemaFILE, 'r') as f:
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
		sublime.message_dialog("Please sett the server address.")
		return None

	if gemaSERVER == '':
		sublime.run_command('gema_connect')
		if gemaSERVER == '':
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
		gemaFOLDER = info['Folder']

	url = urllib.parse.urljoin(gemaSERVER, path)
	load = urllib.parse.urlencode(data).encode('utf-8')
	req = urllib.request.Request(url, load, method=method)
	try:
		with urllib.request.urlopen(req, None, gemaTIMEOUT) as response:
			return response.read().decode(encoding="utf-8")
	except urllib.error.HTTPError as err:
		sublime.message_dialog("{0}".format(err))
	except urllib.error.URLError as err:
		sublime.message_dialog("{0}\nCannot connect to server.".format(err))
	print('Something is wrong')
	return None

# ------------------------------------------------------------------
class gemaViewBulletinBoard(sublime_plugin.ApplicationCommand):
	def run(self):
		response = gemaRequest('teacher_gets_passcode', {})
		if response.startswith('Unauthorized'):
			sublime.message_dialog('Unauthorized')
		else:
			global gemaSERVER
			p = urllib.parse.urlencode({'pc' : response})
			with open(gemaFILE, 'r') as f:
				info = json.loads(f.read())
			webbrowser.open(gemaSERVER + '/view_bulletin_board?' + p)

# ------------------------------------------------------------------
class gemaAddBulletin(sublime_plugin.TextCommand):
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
		response = gemaRequest('teacher_adds_bulletin_page', {'content':content})
		if response:
			sublime.message_dialog(response)

# ------------------------------------------------------------------
class gemaPutBack(sublime_plugin.TextCommand):
	def run(self, edit):
		fname = self.view.file_name()
		sid = os.path.basename(os.path.dirname(fname))
		response = gemaRequest('teacher_puts_back', {'sid':sid})
		if response:
			sublime.message_dialog(response)
		if response != 'Unknown submission.':
			self.view.window().run_command('close')

# ------------------------------------------------------------------
def gema_grade(self, edit, decision):
	fname = self.view.file_name()
	changed = False
	sid = os.path.basename(os.path.dirname(fname))
	if decision=='dismissed':
		content = ''
	else:
		content = self.view.substr(sublime.Region(0, self.view.size())).strip()
		if sid in gemaStudentSubmissions:
			if gemaStudentSubmissions[sid].strip() != content.strip():
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
	response = gemaRequest('teacher_grades', data)
	if response:
		sublime.message_dialog(response)
		self.view.window().run_command('close')

# ------------------------------------------------------------------
class gemaUngrade(sublime_plugin.TextCommand):
	def run(self, edit):
		gema_grade(self, edit, "ungraded")

class gemaGradeCorrect(sublime_plugin.TextCommand):
	def run(self, edit):
		gema_grade(self, edit, "correct")

class gemaGradeIncorrect(sublime_plugin.TextCommand):
	def run(self, edit):
		gema_grade(self, edit, "incorrect")

class gemaDismissed(sublime_plugin.TextCommand):
	def run(self, edit):
		gema_grade(self, edit, "dismissed")

# ------------------------------------------------------------------
def gema_gets(self, index, priority):
	global gemaStudentSubmissions
	response = gemaRequest('teacher_gets', {'index':index, 'priority':priority})
	if response is not None:
		sub = json.loads(response)
		if sub['Content'] != '':
			filename = sub['Filename']
			sid = str(sub['Sid'])
			dir = os.path.join(gemaFOLDER, sid)
			if not os.path.exists(dir):
				os.mkdir(dir)
			local_file = os.path.join(dir, filename)
			with open(local_file, 'w', encoding='utf-8') as fp:
				fp.write(sub['Content'])
			gemaStudentSubmissions[sid] = sub['Content']
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
class gemaSeeQueue(sublime_plugin.ApplicationCommand):
	def run(self):
		response = gemaRequest('teacher_gets_queue', {})
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
		gema_gets(self, selected, 1)

# ------------------------------------------------------------------
# Priorities: 1 (I got it), 2 (I need help),
# ------------------------------------------------------------------
class gemaGetPrioritized(sublime_plugin.ApplicationCommand):
	def run(self):
		gema_gets(self, -1, 0)

class gemaGetFromNeedHelp(sublime_plugin.ApplicationCommand):
	def run(self):
		gema_gets(self, -1, 2)

class gemaGetFromOk(sublime_plugin.ApplicationCommand):
	def run(self):
		gema_gets(self, -1, 1)

# ------------------------------------------------------------------
class gemaSetLocalFolder(sublime_plugin.ApplicationCommand):
	def run(self):
		try:
			with open(gemaFILE, 'r') as f:
				info = json.loads(f.read())
		except:
			info = dict()
		if 'Folder' not in info:
			info['Folder'] = os.path.join(os.path.expanduser('~'), 'GEMA')
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
				with open(gemaFILE, 'r') as f:
					info = json.loads(f.read())
			except:
				info = dict()
			info['Folder'] = folder
			if not os.path.exists(folder):
				try:
					os.mkdir(folder)
					with open(gemaFILE, 'w') as f:
						f.write(json.dumps(info, indent=4))
				except:
					sublime.message_dialog('Could not create {}.'.format(folder))
			else:
				with open(gemaFILE, 'w') as f:
					f.write(json.dumps(info, indent=4))
				sublime.message_dialog('Folder exists. Will use it to store working files.')
		else:
			sublime.message_dialog("Folder name cannot be empty.")

# ------------------------------------------------------------------
class gemaConnect(sublime_plugin.ApplicationCommand):
	def run(self):
		try:
			with open(gemaFILE, 'r') as f:
				info = json.loads(f.read())
		except:
			info = dict()

		if 'CourseId' not in info:
			sublime.message_dialog("Please set the course id.")
			return None

		if 'Server' not in info:
			sublime.message_dialog("Please set server address.")
			return None

		global gemaSERVER
		url = urllib.parse.urljoin(info['Server'], 'ask')
		load = urllib.parse.urlencode({'who':info['CourseId']}).encode('utf-8')
		req = urllib.request.Request(url, load)
		try:
			with urllib.request.urlopen(req, None, gemaTIMEOUT) as response:
				server = response.read().decode(encoding="utf-8")
				try:
					with open(gemaFILE, 'r') as f:
						info = json.loads(f.read())
				except:
					info = dict()
				if not server.startswith('http://'):
					sublime.message_dialog('Unable to get address.')
					return
				gemaSERVER = server
				sublime.status_message('Connected')
		except urllib.error.HTTPError as err:
			sublime.message_dialog("{0}".format(err))
		except urllib.error.URLError as err:
			sublime.message_dialog("{0}\nCannot connect to server.".format(err))

# ------------------------------------------------------------------
class gemaCompleteRegistration(sublime_plugin.ApplicationCommand):
	def run(self):
		try:
			with open(gemaFILE, 'r') as f:
				info = json.loads(f.read())
		except:
			info = dict()

		if 'CourseId' not in info:
			sublime.message_dialog('Please enter course id.')
			return

		if 'Name' not in info:
			sublime.message_dialog('Please enter assigned username.')
			return

		response = gemaRequest(
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
		with open(gemaFILE, 'w') as f:
			f.write(json.dumps(info, indent=4))

# ------------------------------------------------------------------
class gemaSetServerAddress(sublime_plugin.ApplicationCommand):
	def run(self):
		try:
			with open(gemaFILE, 'r') as f:
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
				with open(gemaFILE, 'r') as f:
					info = json.loads(f.read())
			except:
				info = dict()
			if not addr.startswith('http://'):
				addr = 'http://' + addr
			info['Server'] = addr
			with open(gemaFILE, 'w') as f:
				f.write(json.dumps(info, indent=4))
		else:
			sublime.message_dialog("Server address cannot be empty.")

# ------------------------------------------------------------------
class gemaSetCourseId(sublime_plugin.ApplicationCommand):
	def run(self):
		try:
			with open(gemaFILE, 'r') as f:
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
				with open(gemaFILE, 'r') as f:
					info = json.loads(f.read())
			except:
				info = dict()
			info['CourseId'] = cid
			with open(gemaFILE, 'w') as f:
				f.write(json.dumps(info, indent=4))
			sublime.message_dialog('Course id is set to ' + cid)
		else:
			sublime.message_dialog("Course id cannot be empty.")

# ------------------------------------------------------------------
class gemaSetName(sublime_plugin.ApplicationCommand):
	def run(self):
		try:
			with open(gemaFILE, 'r') as f:
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
				with open(gemaFILE, 'r') as f:
					info = json.loads(f.read())
			except:
				info = dict()
			info['Name'] = name
			with open(gemaFILE, 'w') as f:
				f.write(json.dumps(info, indent=4))
			sublime.message_dialog('Assigned name is set to ' + name)
		else:
			sublime.message_dialog("Name cannot be empty.")

# ------------------------------------------------------------------
class gemaUpdate(sublime_plugin.WindowCommand):
	def run(self):
		package_path = os.path.join(sublime.packages_path(), "GEMAssistant");
		try:
			version = open(os.path.join(package_path, "VERSION")).read()
		except:
			version = 0
		if sublime.ok_cancel_dialog("Current version is {}. Click OK to update.".format(version)):
			if not os.path.isdir(package_path):
				os.mkdir(package_path)
			module_file = os.path.join(package_path, "GEMAssistant.py")
			menu_file = os.path.join(package_path, "Main.sublime-menu")
			keymap_file = os.path.join(package_path, "Default.sublime-keymap")
			version_file = os.path.join(package_path, "version.go")
			urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/GEM/master/src/GEMAssistant/GEMAssistant.py", module_file)
			urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/GEM/master/src/GEMAssistant/Main.sublime-menu", menu_file)
			urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/GEM/master/src/GEMAssistant/Default.sublime-keymap", keymap_file)
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


