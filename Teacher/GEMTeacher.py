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

gemtFILE = os.path.join(os.path.dirname(os.path.realpath(__file__)), "info")
gemtPostDir = os.path.join(os.path.dirname(os.path.realpath(__file__)), "Posts")
gemtProblemDivider = '<PROBLEM DIVIDER>'
gemtAnswerTag = 'ANSWER:'
gemtTIMEOUT = 7

# ------------------------------------------------------------------
def gemtRequest(path, data, authenticated=True, localhost=False):
	try:
		with open(gemtFILE, 'r') as f:
			info = json.loads(f.read())
	except:
		info = dict()

	if 'Server' not in info or localhost:
		info['Server'] = 'http://localhost:8080'

	if authenticated:
		if 'Name' not in info or 'Password' not in info:
			sublime.message_dialog("Please ask to setup a new teacher account and register.")
			return None
		data['name'] = info['Name']
		data['password'] = info['Password']
		data['uid'] = info['Uid']
		data['role'] = 'teacher'

	url = urllib.parse.urljoin(info['Server'], path)
	load = urllib.parse.urlencode(data).encode('utf-8')
	req = urllib.request.Request(url, load)
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
class gemtTest(sublime_plugin.WindowCommand):
	def run(self):
		response = gemtRequest('test', {})
		sublime.message_dialog(response)

# ------------------------------------------------------------------
def gemt_gets(self, index, priority):
	response = gemtRequest('teacher_gets', {'index':index, 'priority':priority})
	if response is not None:
		sub = json.loads(response)
		if sub['Sid'] > 0:
			ext = sub['Ext'] or 'txt'
			pid, sid, uid = sub['Pid'], sub['Sid'], sub['Uid']
			fname = 'gemt_{}_{}_{}.{}'.format(uid,pid,sid,ext)
			if not os.path.isdir(gemtPostDir):
				os.mkdir(gemtPostDir)
			fname = os.path.join(gemtPostDir, fname)
			with open(fname, 'w', encoding='utf-8') as fp:
				fp.write(sub['Content'])
			sublime.active_window().open_file(fname)
		else:
			sublime.message_dialog('No submission with index {} or priority {}.'.format(index,priority))

# ------------------------------------------------------------------
class gemtGetSub(sublime_plugin.WindowCommand):
	def run(self):
		gemt_gets(self, -1, -1)

class gemtGetSubFour(sublime_plugin.WindowCommand):
	def run(self):
		gemt_gets(self, -1, 4)

class gemtGetSubThree(sublime_plugin.WindowCommand):
	def run(self):
		gemt_gets(self, -1, 3)

class gemtGetSubTwo(sublime_plugin.WindowCommand):
	def run(self):
		gemt_gets(self, -1, 2)

class gemtGetSubOne(sublime_plugin.WindowCommand):
	def run(self):
		gemt_gets(self, -1, 1)

# ------------------------------------------------------------------
class gemtBroadcast(sublime_plugin.TextCommand):
	def run(self, edit):
		fname = self.view.file_name()
		if fname is None:
			sublime.message_dialog('Cannot broadcast unsaved content.')
			return
		ext = fname.rsplit('.',1)[-1]
		content = self.view.substr(sublime.Region(0, self.view.size())).lstrip()
		if content == '':
			sublime.message_dialog("File is empty.")
			return
		data = {
			'content': 			content,
			'ext': 				ext,
			'problem_divider': 	gemtProblemDivider,
			'answer_tag': 		gemtAnswerTag,
		}
		response = gemtRequest('teacher_broadcasts', data)
		if response is not None:
			sublime.status_message(response)

# ------------------------------------------------------------------
class gemtSetupNewTeacher(sublime_plugin.WindowCommand):
	def run(self):
		if sublime.ok_cancel_dialog("This can only be done if SublimeText and the server are running on localhost."):
			sublime.active_window().show_input_panel('Enter username:',
				'',
				self.process,
				None,
				None)

	def process(self, name):
		name = name.strip()
		response = gemtRequest('teacher_adds_ta', {'name':name}, authenticated=False, localhost=True)
		sublime.message_dialog(response)

# ------------------------------------------------------------------
class gemtRegister(sublime_plugin.WindowCommand):
	def run(self):
		if sublime.ok_cancel_dialog("Register a username that was temporarily added on localhost."):
			sublime.active_window().show_input_panel('Enter username:',
				'',
				self.process,
				None,
				None)

	def process(self, name):
		name = name.strip()
		response = gemtRequest('teacher_registers', {'name':name}, authenticated=False)
		if response == 'Failed':
			sublime.message_dialog('This name is not registered. Ask the teacher in charge')
		else:
			uid, password = response.split(',')
			try:
				with open(gemtFILE, 'r') as f:
					info = json.loads(f.read())
			except:
				info = dict()
			info['Uid'] = int(uid)
			info['Password'] = password
			info['Name'] = name
			with open(gemtFILE, 'w') as f:
				f.write(json.dumps(info, indent=4))
			sublime.message_dialog('{} registered'.format(name))

# ------------------------------------------------------------------
class gemtSetServerAddress(sublime_plugin.WindowCommand):
	def run(self):
		try:
			with open(gemtFILE, 'r') as f:
				info = json.loads(f.read())
		except:
			info = dict()
		if 'Server' not in info:
			info['Server'] = 'http://localhost:8080'
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
		else:
			sublime.message_dialog("Server address cannot be empty.")

# ------------------------------------------------------------------
