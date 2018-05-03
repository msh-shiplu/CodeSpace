#-------------------------------------------------------------------
# All students submit solutions to a problem on their white boards.
#-------------------------------------------------------------------
import requests
import os
import sys
import json
import datetime
import random
from common import *

#-------------------------------------------------------------------
# Get content from white board and then submit
#-------------------------------------------------------------------
def students_solve(answer=''):
	files = [ os.path.join(StudentDir, f) for f in os.listdir(StudentDir) ]
	data = {}
	for file in files:
		info = json.load(open(file))
		data['name'] = info['Name']
		data['password'] = info['Password']
		data['uid'] = info['Uid']
		data['role'] = 'student'
		if not info['Server'].startswith('http://'):
			info['Server'] = 'http://' + info['Server']

		print('Student', info['Name'], info['Uid'])
		# 1. Get the problems.
		url = urllib.parse.urljoin(Config['Server'], 'student_gets')
		load = urllib.parse.urlencode(data).encode('utf-8')
		req = urllib.request.Request(url, load, method='POST')
		with urllib.request.urlopen(req, None, 5) as response:
			whiteboard = response.read().decode(encoding="utf-8")
			problems = json.loads(whiteboard)
			for problem in problems:
				# 2. Get problem information and prepare for submission
				base, ext = problem['Filename'].split('.')
				data['filename'] = '{}_{}.{}'.format(base, problem['Pid'], ext)
				data['pid'] = int(problem['Pid'])
				data['content'] = 'solution from {} ({})'.format(info['Name'],info['Uid'])
				if random.random() < info['Ability']:
					data['content'] = 'Correct ' + data['content']
				else:
					data['content'] = 'Incorrect ' + data['content']
				data['priority'] = 1
				if answer == '':
					data['answer'] = ''
				else:
					data['answer'] = random.choice(list(answer))

				# 3. Submit problem
				url = urllib.parse.urljoin(Config['Server'], 'student_shares')
				load = urllib.parse.urlencode(data).encode('utf-8')
				submit_req = urllib.request.Request(url, load, method='POST')
				with urllib.request.urlopen(submit_req, None, 5) as submit_response:
					print('  Problem', problem['Pid'], submit_response.read().decode(encoding="utf-8"))


#-------------------------------------------------------------------
if __name__ == '__main__':
	if len(sys.argv) == 1:
		students_solve()
	elif len(sys.argv) == 2:
		students_solve(sys.argv[1])
	else:
		print('\n\tUsage: python {} [answer_choices_string]'.format(sys.argv[0]))

