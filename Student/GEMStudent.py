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

gemsFILE = os.path.join(os.path.dirname(os.path.realpath(__file__)), "info")
gemsFOLDER = os.path.join(os.path.expanduser('~'), 'GEM')

# ------------------------------------------------------------------------------
def gemsRequest(path, data, headers={}, is_json=False, authenticated=True):
	try:
		with open(gemsFILE, 'r') as f:
			info = json.loads(f.read())
	except:
		info = dict()

	if 'Server' not in info:
		sublime.message_dialog("Please set server address.")
		return None
	data['server'] = info['Server']

	if authenticated:
		if 'Name' not in info or 'Password' not in info:
			sublime.message_dialog("Please ask to setup a new teacher account and register.")
			return None
		data['name'] = info['Name']
		data['password'] = info['Password']
		data['uid'] = info['Uid']

	url = urllib.parse.urljoin(info['Server'], path)
	if is_json:
		load = json.dumps(data).encode('utf-8')
	else:
		load = urllib.parse.urlencode(data).encode('utf-8')
	req = urllib.request.Request(url, load, headers=headers)
	try:
		with urllib.request.urlopen(req, None, TIMEOUT) as response:
			return response.read().decode(encoding="utf-8")
	except urllib.error.HTTPError as err:
		sublime.message_dialog("{0}".format(err))
	except urllib.error.URLError as err:
		sublime.message_dialog("{0}\nCannot connect to server.".format(err))
	print('Something is wrong')
	return None

# ------------------------------------------------------------------
class gemsRegisterCommand(sublime_plugin.WindowCommand):
	def run(self):
		sublime.active_window().show_input_panel('Enter username:',
			'',
			self.process,
			None,
			None)

	def process(self, name):
		name = name.strip()
		response = gemsRequest('register_student', {'name':name}, authenticated=False)
		if response == 'exist':
			sublime.message_dialog('{} exists. Choose a different name.'.format(name))
		else:
			uid, password = response.split(',')
			try:
				with open(gemsFILE, 'r') as f:
					info = json.loads(f.read())
			except:
				info = dict()
			info['Uid'] = int(uid)
			info['Password'] = password
			info['Name'] = name
			with open(gemsFILE, 'w') as f:
				f.write(json.dumps(info, indent=4))
			sublime.message_dialog('{} registered'.format(name))

# ------------------------------------------------------------------
class gemsSetServerAddress(sublime_plugin.WindowCommand):
	def run(self):
		try:
			with open(gemsFILE, 'r') as f:
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

# ------------------------------------------------------------------------------
