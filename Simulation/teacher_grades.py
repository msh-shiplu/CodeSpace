#-------------------------------------------------------------------
# Teacher grades submitted solutions.
#------------------------------------------------------------------
import requests
import os
import sys
import json
import datetime
import random
from common import *

#-------------------------------------------------------------------
def teacher_grades():
	data = dict(index=-1, priority=0)
	while True:
		sub = json.loads(teacher_request('teacher_gets', data))
		if sub['Uid'] == 0:
			break
		if sub['Content'].startswith('Correct'):
			decision = 'correct'
		else:
			decision = 'incorrect'
		grade_data = dict(
			content = sub['Content'],
			decision = decision,
			changed = False,
			pid = sub['Pid'],
			sid = sub['Sid'],
			stid = sub['Uid'],
		)
		output = teacher_request('teacher_grades', grade_data)
		print('Grading {}: {}'.format(sub['Uid'], output))

#-------------------------------------------------------------------
if __name__ == '__main__':
	teacher_grades()