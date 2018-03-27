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
gemaStudentSubmissions = {}

# ------------------------------------------------------------------
def gemaRequest(path, data, authenticated=True, localhost=False, method='POST'):
	global gemaFOLDER
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
class gemaViewBulletinBoard(sublime_plugin.WindowCommand):
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
		content = '\n' + self.view.substr(sublime.Region(beg,end)) + '\n'
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
		basename = os.path.basename(fname)
		if not basename.startswith('gema') or basename.count('_')!=2:
			sublime.message_dialog('This is not a student submission.')
			return
		prefix, ext = basename.rsplit('.', 1)
		prefix = prefix[4:]
		stid, pid, sid = prefix.split('_')
		try:
			stid = int(stid)
			sid = int(sid)
		except:
			sublime.message_dialog('This is not a student submission.')
			return
		try:
			pid = int(pid)
		except:
			pid = 0
		content = self.view.substr(sublime.Region(0, self.view.size())).strip()
		data = dict(
			sid = sid,
			stid = stid,
			pid = pid,
			content = content,
			ext = ext,
			priority = gemsHighestPriority,
		)
		response = gemaRequest('teacher_puts_back', data)
		if response:
			sublime.message_dialog(response)

# ------------------------------------------------------------------
def gema_grade(self, edit, decision):
	fname = self.view.file_name()
	basename = os.path.basename(fname)
	if not basename.startswith('gema') or basename.count('_')!=2:
		sublime.message_dialog('This is not a student submission.')
		return
	prefix, ext = basename.rsplit('.', 1)
	prefix = prefix[4:]
	stid, pid, sid = prefix.split('_')
	try:
		stid = int(stid)
		sid = int(sid)
	except:
		sublime.message_dialog('This is not a student submission.')
		return
	try:
		pid = int(pid)
	except:
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
		content = content,
		ext = ext,
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
def gema_rand_chars(n):
	letters = 'abcdefghijklmkopqrstuvwxyzABCDEFGHIJKLMLOPQRSTUVWXYZ'
	return ''.join(random.choice(letters) for i in range(n))

# ------------------------------------------------------------------
def gema_gets(self, index, priority):
	global gemaStudentSubmissions
	response = gemaRequest('teacher_gets', {'index':index, 'priority':priority})
	if response is not None:
		sub = json.loads(response)
		if sub['Content'] != '':
			ext = sub['Ext'] or 'txt'
			pid, sid, uid = sub['Pid'], sub['Sid'], sub['Uid']
			if pid == 0:
				pid = gema_rand_chars(4)
			fname = 'gema{}_{}_{}.{}'.format(uid,pid,sid,ext)
			fname = os.path.join(gemaFOLDER, fname)
			with open(fname, 'w', encoding='utf-8') as fp:
				fp.write(sub['Content'])
			gemaStudentSubmissions[pid] = sub['Content']
			sublime.active_window().open_file(fname)
		else:
			sublime.message_dialog('No submission with index {} or priority {}.'.format(index,priority))

# ------------------------------------------------------------------
# Priorities: 1 (I got it), 2 (I need help),
# ------------------------------------------------------------------
class gemtGetPrioritized(sublime_plugin.WindowCommand):
	def run(self):
		gemt_gets(self, -1, 0)

class gemaGetFromNeedHelp(sublime_plugin.WindowCommand):
	def run(self):
		gema_gets(self, -1, 2)

class gemaGetFromOk(sublime_plugin.WindowCommand):
	def run(self):
		gema_gets(self, -1, 1)

# ------------------------------------------------------------------
class gemaSetServerAddress(sublime_plugin.WindowCommand):
	def run(self):
		try:
			with open(gemaFILE, 'r') as f:
				info = json.loads(f.read())
		except:
			info = dict()
		if 'Server' not in info:
			info['Server'] = 'http://x.x.x.x:8080'
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
class gemaSetLocalFolder(sublime_plugin.WindowCommand):
	def run(self):
		try:
			with open(gemaFILE, 'r') as f:
				info = json.loads(f.read())
		except:
			info = dict()

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
class gemaRegister(sublime_plugin.WindowCommand):
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

		mesg = 'Enter a new name'
		if 'Name' in info:
			mesg = '{} is already registered. Enter a new name or Esc:'.format(info['Name'])

		if 'Name' not in info:
			info['Name'] = ''

		if sublime.ok_cancel_dialog("Register a username that has been added by the teacher."):
			sublime.active_window().show_input_panel(
				mesg,
				info['Name'],
				self.process,
				None,
				None,
			)

	def process(self, name):
		name = name.strip()
		response = gemaRequest('ta_registers', {'name':name}, authenticated=False)
		if response == 'Failed':
			sublime.message_dialog('This name is not registered. Ask the teacher to add it.')
		else:
			uid, password = response.split(',')
			try:
				with open(gemaFILE, 'r') as f:
					info = json.loads(f.read())
			except:
				info = dict()
			info['Uid'] = int(uid)
			info['Password'] = password
			info['Name'] = name
			with open(gemaFILE, 'w') as f:
				f.write(json.dumps(info, indent=4))
			sublime.message_dialog('{} registered'.format(name))
