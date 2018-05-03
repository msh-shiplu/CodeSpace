import json
import urllib.parse
import urllib.request
import os
import random

# ConfigFile = 'config_home.json'
ConfigFile = 'config_um.json'
Config = None
TeacherDir = 'Teachers'
StudentDir = 'Students'

#-------------------------------------------------------------------
def init():
	global Config
	Config = json.load(open(ConfigFile))
	if not Config['Server'].startswith('http://'):
		Config['Server'] = 'http://' + Config['Server']

init()

#-------------------------------------------------------------------
def request(path, data):
	url = urllib.parse.urljoin(Config['Server'], path)
	load = urllib.parse.urlencode(data).encode('utf-8')
	req = urllib.request.Request(url, load, method='POST')
	with urllib.request.urlopen(req, None, 5) as response:
		return response.read().decode(encoding="utf-8")

#-------------------------------------------------------------------
def teacher_request(path, data, teacher=None):
	files = [ os.path.join(TeacherDir, f) for f in os.listdir(TeacherDir) ]
	if teacher is None:
		file = random.choice(files)
		info = json.load(open(file))
	data['name'] = info['Name']
	data['password'] = info['Password']
	data['uid'] = info['Uid']
	data['role'] = 'teacher'
	url = urllib.parse.urljoin(Config['Server'], path)
	load = urllib.parse.urlencode(data).encode('utf-8')
	req = urllib.request.Request(url, load, method='POST')
	with urllib.request.urlopen(req, None, 5) as response:
		return response.read().decode(encoding="utf-8")

#-------------------------------------------------------------------

