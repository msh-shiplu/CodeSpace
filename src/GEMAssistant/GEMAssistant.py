# Code4Brownies - Instructor module
# Author: Vinhthuy Phan, 2015-2017
#

import sublime, sublime_plugin
import urllib.parse
import urllib.request
import os
import json
import socket
import webbrowser
import random

gemaFILE = os.path.join(os.path.dirname(os.path.realpath(__file__)), "info")
gemaFOLDER = ''
gemaTIMEOUT = 7
gemaHighestPriority = 2
gemaConnected = False
gemaStudentSubmissions = {}

# ------------------------------------------------------------------
# These functionalities are unique to the GEMAssistant module
# ------------------------------------------------------------------
class gemaPutBack(sublime_plugin.TextCommand):
	def run(self, edit):
		fname = self.view.file_name()
		basename = os.path.basename(fname)
		if not basename.startswith('gemt') or basename.count('_') < 2:
			sublime.message_dialog('This is not a student submission.')
			return
		if '.' in basename:
			prefix, ext = basename.rsplit('.', 1)
		else:
			prefix = basename
		prefix = prefix[4:]
		stid, pid, sid = prefix.split('_')
		try:
			stid = int(stid)
			pid = int(pid)
			sid = int(sid)
		except:
			sublime.message_dialog('This is not a student submission.')
			return
		content = self.view.substr(sublime.Region(0, self.view.size())).strip()
		data = dict(
			sid = sid,
			stid = stid,
			pid = pid,
			content = content,
			priority = gemaHighestPriority,
		)
		response = gemaRequest('teacher_puts_back', data)
		if response:
			sublime.message_dialog(response)
		self.view.window().run_command('close')

# ------------------------------------------------------------------
# These functionalities should be identical to the GEMTeacher module
# ------------------------------------------------------------------
def gemaRequest(path, data, authenticated=True, method='POST'):
	global gemaFOLDER, gemaConnected
	if not gemaConnected:
		sublime.run_command('gema_connect')
		gemaConnected = True

	try:
		with open(gemaFILE, 'r') as f:
			info = json.loads(f.read())
	except:
		info = dict()

	if 'Folder' not in info:
		sublime.message_dialog("Please set a local folder for keeping working files.")
		return None

	if 'Server' not in info:
		sublime.message_dialog("Please sett the server address.")
		return None

	if authenticated:
		if 'Name' not in info or 'Password' not in info:
			sublime.message_dialog("Please ask to setup a new teacher account and register.")
			return None
		data['name'] = info['Name']
		data['password'] = info['Password']
		data['uid'] = info['Uid']
		data['role'] = 'teacher'
		gemaFOLDER = info['Folder']

	url = urllib.parse.urljoin(info['Server'], path)
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
			p = urllib.parse.urlencode({'pc' : response})
			with open(gemaFILE, 'r') as f:
				info = json.loads(f.read())
			webbrowser.open(info['Server'] + '/view_bulletin_board?' + p)

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
def gema_grade(self, edit, decision):
	fname = self.view.file_name()
	basename = os.path.basename(fname)
	if not basename.startswith('gemt') or basename.count('_') < 2:
		sublime.message_dialog('This is not a student submission.')
		return
	if '.' in basename:
		prefix, ext = basename.rsplit('.', 1)
	else:
		prefix = basename
	prefix = prefix[4:]
	stid, pid, sid = prefix.split('_')
	try:
		stid = int(stid)
		pid = int(pid)
		sid = int(sid)
	except:
		sublime.message_dialog('This is not a student submission.')
		return
	if pid == 0:
		sublime.message_dialog('This is not a graded problem.')
		return
	changed = False
	if decision=='dismissed':
		content = ''
	else:
		content = self.view.substr(sublime.Region(0, self.view.size())).strip()
		if pid in gemaStudentSubmissions and content.strip()!=gemaStudentSubmissions[pid].strip():
			changed = True
	data = dict(
		stid = stid,
		pid = pid,
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
			if '.' in filename:
				ext = filename.rsplit('.',1)[1]
			else:
				ext = 'txt'
			pid, sid, uid = sub['Pid'], sub['Sid'], sub['Uid']
			fname = 'gemt{}_{}_{}.{}'.format(uid,pid,sid,ext)
			fname = os.path.join(gemaFOLDER, fname)
			with open(fname, 'w', encoding='utf-8') as fp:
				fp.write(sub['Content'])
			gemaStudentSubmissions[pid] = sub['Content']
			if sublime.active_window().id() == 0:
				sublime.run_command('new_window')
			sublime.active_window().open_file(fname)
		elif priority == 0:
			sublime.message_dialog('There are no submissions.')
		elif priority > 0:
			sublime.message_dialog('There are no submissions with priority {}.'.format(priority))
		elif index >= 0:
			sublime.message_dialog('There are no submission with index {}.'.format(index))

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
		sublime.active_window().show_input_panel("This folder will be used to store working files.",
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
			with urllib.request.urlopen(req, None, gemaTIMEOUT) as response:
				server = response.read().decode(encoding="utf-8")
				try:
					with open(gemaFILE, 'r') as f:
						info = json.loads(f.read())
				except:
					info = dict()
				if not server.startswith('http://'):
					sublime.message_dialog(server)
					return
				info['Server'] = server
				with open(gemaFILE, 'w') as f:
					f.write(json.dumps(info, indent=4))
				sublime.message_dialog('Connected to server at {}'.format(server))
		except urllib.error.HTTPError as err:
			sublime.message_dialog("{0}".format(err))
		except urllib.error.URLError as err:
			sublime.message_dialog("{0}\nCannot connect to name server.".format(err))

	# ------------------------------------------------------------------
	def ask_to_set_server_address(self, info):
		try:
			with open(gemaFILE, 'r') as f:
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
class gemaCompleteRegistration(sublime_plugin.ApplicationCommand):
	def run(self):
		try:
			with open(gemaFILE, 'r') as f:
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
		response = gemaRequest(
			'complete_registration',
			{'name':name.strip(), 'role':'teacher'},
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
			with open(gemaFILE, 'r') as f:
				info = json.loads(f.read())
		except:
			info = dict()

		info['Uid'] = int(uid)
		info['Password'] = password.strip()
		info['Name'] = name.strip()
		info['CourseId'] = course_id.strip()
		if name_server != "":
			info['NameServer'] = name_server
		with open(gemaFILE, 'w') as f:
			f.write(json.dumps(info, indent=4))
		sublime.message_dialog('{} registered for {}'.format(name, course_id))

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
class gemaUpdate(sublime_plugin.WindowCommand):
	def run(self):
		package_path = os.path.join(sublime.packages_path(), "GEMAssistant");
		try:
			version = open(os.path.join(package_path, "VERSION")).read()
		except:
			version = 0
		if sublime.ok_cancel_dialog("Current version is %s. Click OK to update.".format(version)):
			if not os.path.isdir(package_path):
				os.mkdir(package_path)
			module_file = os.path.join(package_path, "GEMAssistant.py")
			menu_file = os.path.join(package_path, "Main.sublime-menu")
			version_file = os.path.join(package_path, "version.go")
			try:
				urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/GEM/master/src/GEMAssistant/GEMAssistant.py", module_file)
				urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/GEM/master/src/GEMAssistant/Main.sublime-menu", menu_file)
				urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/GEM/master/src/version.go", version_file)
				lines = open(version_file).readlines()
				for line in lines:
					if line.strip().startswith('const VERSION ='):
						prefix, version = line.strip().split('const VERSION =')
						version = float(version)
						break
				os.remove(version_file)
				with open(os.path.join(package_path, "VERSION"), 'w') as f:
					f.write(version)
				sublime.message_dialog("GEM has been updated to version %s." % version)
			except:
				sublime.message_dialog("A problem occurred during update.")

# ------------------------------------------------------------------

