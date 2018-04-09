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

gemtOrTag = '<GEM_OR>'
gemtAndTag = '<GEM_AND>'
gemtSeqTag = '<GEM_NEXT>'
gemtAnswerTag = 'ANSWER:'
gemtFOLDER = ''
gemtTIMEOUT = 7
gemtStudentSubmissions = {}
gemtFILE = os.path.join(os.path.dirname(os.path.realpath(__file__)), "info")

# ------------------------------------------------------------------
# ------------------------------------------------------------------
# These functionalities are unique to the GEMTeacher module
# ------------------------------------------------------------------
# ------------------------------------------------------------------

class gemtTest(sublime_plugin.WindowCommand):
	def run(self):
		response = gemtRequest('test', {})
		sublime.message_dialog(response)

# ------------------------------------------------------------------
class gemtShare(sublime_plugin.TextCommand):
	def run(self, edit):
		fname = self.view.file_name()
		if fname is None:
			sublime.message_dialog('Cannot share unsaved content.')
			return
		basename = os.path.basename(fname)
		# ext = basename.rsplit('.',1)[-1]
		content = self.view.substr(sublime.Region(0, self.view.size())).lstrip()
		if content == '':
			sublime.message_dialog("File is empty.")
			return
		data = {
			'content': 	content,
			'filename': basename,
			# 'ext': 		ext,
		}
		response = gemtRequest('teacher_shares', data)
		if response is not None:
			sublime.status_message(response)

# ------------------------------------------------------------------
def gemt_get_problem_info(fname):
	merit, effort, attempts = 0, 0, -1
	ext = fname.rsplit('.', 1)[-1]
	with open(fname, 'r', encoding='utf-8') as fp:
		content = fp.read()
	items = content.split('\n',1)
	if len(items)==1:
		sublime.message_dialog('Improper problem definition')
		raise Exception('Improper problem definition')
	first_line, body = items[0], items[1]
	prefix = '//'
	if first_line.startswith('#'):
		prefix = '#'
		first_line = first_line.strip('# ')
	elif first_line.startswith('/'):
		first_line = first_line.strip('/ ')
	try:
		items = first_line.split(' ')
		merit, effort = int(items[0]), int(items[1])
		if len(items) > 2:
			attempts = int(items[2])
	except:
		sublime.message_dialog('Improper problem definition')
		raise Exception('Improper problem definition')
	if merit < effort:
		sublime.message_dialog('Merit points ({}) should be higher than effort points {}.'.format(merit, effort))
		raise Exception('Merit points lower than effort points.')

	items = body.split(gemtAnswerTag)
	if len(items) > 2:
		sublime.message_dialog('This problem has {} answers. There should be at most 1.'.format(len(items)-1))
		raise Exception('Too many answers')

	body = '{} {} points, {} for effort. Maximum attempts: {}.\n{}'.format(
		prefix, merit, effort, attempts, items[0])
	answer = ''
	if len(items) == 2:
		body += '\n{} '.format(gemtAnswerTag)
		answer = items[1].strip()

	basename = os.path.basename(fname)
	return body, answer, str(merit), str(effort), str(attempts), basename, ext


# ------------------------------------------------------------------
# Input: file names
# Output:
# - files sorted by difficulty level
# - list of index of the next problem if solution is correct
# - list of index of the next problem if solution is incorrect
# ------------------------------------------------------------------
def gemt_get_next_problems(fns, mode):
	nic = [-1] * len(fns)
	nii = [-1] * len(fns)
	if mode != 'multicast_seq':
		return fns, nic, nii
	# Remove .py or .java
	names = [ f.rsplit('.', 1)[0] for f in fns ]
	for n in names:
		# Properly named files must have at least one "_"
		if n.count('_') < 1:
			return None, None, None
	try:
		nums = [ int(n.rsplit('_', 1)[1]) for n in names ]
		names = [ (nums[i],fns[i]) for i in range(len(fns)) ]
		names = sorted(names)
		# determine next if correct or fns[i]
		for i in range(len(names)):
			for j in range(i+1,len(names)):
				if names[j][0] > names[i][0]:
					nic[i] = j
					break
		# determine next if incorrect or fns[i]
		for i in range(len(names)-1):
			nii[i] = i+1
		# print(names)
		# print(nic)
		# print(nii)
		return [n[1] for n in names], nic, nii
	except:
		return None, None, None

# ------------------------------------------------------------------
# used by gemt_multicast and gemtUnicast to start a new problem
# ------------------------------------------------------------------
def gemt_broadcast(content, answers, merits, efforts, attempts, filenames, exts, tag, mode, nic='', nii=''):
	data = {
		'content': 			content,
		'answers':			answers,
		'filenames':		filenames,
		'exts': 			exts,
		'merits':			merits,
		'efforts':			efforts,
		'attempts':			attempts,
		'divider_tag':	 	tag,
		'mode':				mode,
		'nic': 				nic,
		'nii':				nii,
	}
	response = gemtRequest('teacher_broadcasts', data)
	if response is not None:
		sublime.status_message(response)

# ------------------------------------------------------------------
def gemt_multicast(self, edit, tag, mode, mesg):
	fnames = [ v.file_name() for v in sublime.active_window().views() ]
	fnames = [ fname for fname in fnames if fname is not None ]
	if len(fnames)>0 and sublime.ok_cancel_dialog(mesg):
		content, answers, merits, efforts, attempts, fns, exts = [],[],[],[],[],[],[]
		nic, nii = [], []
		for fname in fnames:
			c, an, m, e, at, fn, ex = gemt_get_problem_info(fname)
			content.append(c)
			answers.append(an)
			merits.append(m)
			efforts.append(e)
			attempts.append(at)
			fns.append(fn)
			exts.append(ex)

		fns, nic, nii = gemt_get_next_problems(fns, mode)
		if fns == None:
			sublime.message_dialog('Could not detect difficulty level in file names.  Example of a correctly named file at level 1: abc_1.py')
			return

		content = '\n{}\n'.format(tag).join(content)
		answers = '\n'.join(answers)
		merits = '\n'.join(merits)
		efforts = '\n'.join(efforts)
		attempts = '\n'.join(attempts)
		fns = '\n'.join(fns)
		exts = '\n'.join(exts)
		nic = '\n'.join([str(i) for i in nic])
		nii = '\n'.join([str(i) for i in nii])

		gemt_broadcast(
			content,
			answers,
			merits,
			efforts,
			attempts,
			fns,
			exts,
			tag,
			mode,
			nic,
			nii,
		)

# ------------------------------------------------------------------
class gemtUnicast(sublime_plugin.TextCommand):
	def run(self, edit):
		if sublime.ok_cancel_dialog('Starting a new problem will close up active problems.  Do you want to proceed?'):
			fname = self.view.file_name()
			if fname is None:
				sublime.message_dialog('Content must be saved first.')
				return
			content, answers, merits, efforts, attempts, fns, exts = gemt_get_problem_info(fname)
			gemt_broadcast(content, answers, merits, efforts, attempts, fns, exts, tag='', mode='unicast')

# ------------------------------------------------------------------
class gemtMulticastOr(sublime_plugin.TextCommand):
	def run(self, edit):
		if sublime.ok_cancel_dialog('Starting new problems will close up active problems.  Do you want to proceed?'):
			gemt_multicast(
				self,
				edit,
				gemtOrTag,
				'multicast_or',
				'Send problems *randomly* to students, where problems are defined in all non-empty tabs in this window?',
			)

# ------------------------------------------------------------------
class gemtMulticastAnd(sublime_plugin.TextCommand):
	def run(self, edit):
		if sublime.ok_cancel_dialog('Starting problems will close up active problems.  Do you want to proceed?'):
			gemt_multicast(
				self,
				edit,
				gemtAndTag,
				'multicast_and',
				'Send problems *simultaneously* to students, where problems are defined in all non-empty tabs in this window?',
			)

# ------------------------------------------------------------------
class gemtMulticastSeq(sublime_plugin.TextCommand):
	def run(self, edit):
		if sublime.ok_cancel_dialog('Starting problems will close up active problems.  Do you want to proceed?'):
			gemt_multicast(
				self,
				edit,
				gemtSeqTag,
				'multicast_seq',
				'Send problems *sequentially* to students, where problems are defined in all non-empty tabs in this window?',
			)

# ------------------------------------------------------------------
class gemtDeactivateProblems(sublime_plugin.WindowCommand):
	def run(self):
		if sublime.ok_cancel_dialog('Do you want to close active problems?  No more submissions are possible until a new problem is started.'):
			passcode = gemtRequest('teacher_gets_passcode', {})
			json_response = gemtRequest('teacher_deactivates_problems', {})
			with open(gemtFILE, 'r') as f:
				info = json.loads(f.read())
			active_pids = json.loads(json_response)
			mesg = "Problems closed."
			if len(active_pids) > 0:
				mesg += " Answers to be summarized."
			sublime.message_dialog(mesg)
			for pid in active_pids:
				p = urllib.parse.urlencode({'pc' : passcode, 'pid':pid})
				webbrowser.open(info['Server'] + '/view_answers?' + p)

# ------------------------------------------------------------------
class gemtClearSubmissions(sublime_plugin.WindowCommand):
	def run(self):
		if sublime.ok_cancel_dialog('Do you want to clear all submissions and white boards?'):
			response = gemtRequest('teacher_clears_submissions', {})
			sublime.message_dialog(response)



# ----------------------------------------------------------------------
# ----------------------------------------------------------------------
# These functionalities below are identical to the GEMAssistant module
# ----------------------------------------------------------------------
# ----------------------------------------------------------------------
def gemtRequest(path, data, authenticated=True, method='POST'):
	global gemtFOLDER
	try:
		with open(gemtFILE, 'r') as f:
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
		gemtFOLDER = info['Folder']

	url = urllib.parse.urljoin(info['Server'], path)
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
class gemtViewBulletinBoard(sublime_plugin.WindowCommand):
	def run(self):
		response = gemtRequest('teacher_gets_passcode', {})
		if response.startswith('Unauthorized'):
			sublime.message_dialog('Unauthorized')
		else:
			p = urllib.parse.urlencode({'pc' : response})
			with open(gemtFILE, 'r') as f:
				info = json.loads(f.read())
			webbrowser.open(info['Server'] + '/view_bulletin_board?' + p)

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
def gemt_grade(self, edit, decision):
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
		if pid in gemtStudentSubmissions and content.strip()!=gemtStudentSubmissions[pid].strip():
			changed = True
	data = dict(
		stid = stid,
		pid = pid,
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
			if '.' in filename:
				ext = filename.rsplit('.',1)[1]
			else:
				ext = 'txt'
			pid, sid, uid = sub['Pid'], sub['Sid'], sub['Uid']
			fname = 'gemt{}_{}_{}.{}'.format(uid,pid,sid,ext)
			fname = os.path.join(gemtFOLDER, fname)
			with open(fname, 'w', encoding='utf-8') as fp:
				fp.write(sub['Content'])
			gemtStudentSubmissions[pid] = sub['Content']
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
class gemtGetPrioritized(sublime_plugin.WindowCommand):
	def run(self):
		gemt_gets(self, -1, 0)

class gemtGetFromNeedHelp(sublime_plugin.WindowCommand):
	def run(self):
		gemt_gets(self, -1, 2)

class gemtGetFromOk(sublime_plugin.WindowCommand):
	def run(self):
		gemt_gets(self, -1, 1)

# ------------------------------------------------------------------
class gemtSetServerAddress(sublime_plugin.WindowCommand):
	def run(self):
		try:
			with open(gemtFILE, 'r') as f:
				info = json.loads(f.read())
		except:
			info = dict()
		if 'Server' not in info:
			info['Server'] = ''
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
class gemtSetLocalFolder(sublime_plugin.WindowCommand):
	def run(self):
		try:
			with open(gemtFILE, 'r') as f:
				info = json.loads(f.read())
		except:
			info = dict()
		if 'Folder' not in info:
			info['Folder'] = os.path.join(os.path.expanduser('~'), 'GEMT')
		sublime.active_window().show_input_panel("This folder will be used to store working files.",
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
class gemtCompleteRegistration(sublime_plugin.WindowCommand):
	def run(self):
		try:
			with open(gemtFILE, 'r') as f:
				info = json.loads(f.read())
		except:
			info = dict()
		if 'Folder' not in info:
			sublime.message_dialog("Please set a local folder for keeping working files.")
			return None

		if 'Server' not in info:
			sublime.message_dialog("Please set server address.")
			return None

		mesg = 'Enter your assigned name'
		if 'Name' in info:
			mesg = '{} is already registered. Enter a new name or Esc:'.format(info['Name'])

		if 'Name' not in info:
			info['Name'] = ''

		if sublime.ok_cancel_dialog("Register an assigned username."):
			sublime.active_window().show_input_panel(
				mesg,
				info['Name'],
				self.process,
				None,
				None,
			)

	def process(self, name):
		name = name.strip()
		response = gemtRequest(
			'complete_registration',
			{'name':name, 'role':'teacher'},
			authenticated=False,
		)
		if response == 'Failed':
			sublime.message_dialog('Failed to complete registration.')
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
