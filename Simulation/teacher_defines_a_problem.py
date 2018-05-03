#-------------------------------------------------------------------
# Teacher posts a new problem
#-------------------------------------------------------------------

import requests
import os
import sys
import json
import datetime
import random
from common import *

learning_objectives = [
	'iterating a list',
	'accumulating result from a list',
	'trace a recursive call',
	'understand variable assignment',
]
#-------------------------------------------------------------------
def teacher_defines_a_problem(answer='', content='', tag='', merit=5, effort=2, attempts=3):
	pid = ''
	for i in range(3):
		pid += random.choice(list('ABCDEFGHIJKLMNOPQRSTUVWXYZ'))
	merit, effort, attempts = merit, effort, attempts
	if tag == '':
		tag = random.choice(learning_objectives)
	if content == '':
		content = '# {} points for correctness, {} points for effort. Maximum attempts: {}.\n{}'.format(merit,effort,attempts,"print('Hello world.')")
	data = dict(
		content = content,
		answers = answer,
		tags = tag,
		filenames = 'simulated_' + pid + '.py',
		merits = merit,
		efforts = effort,
		attempts = attempts,
		mode = 'unicast',
		divider = '',
		nic = '',
		nii = '',
	)
	response = teacher_request('teacher_broadcasts', data)
	print(response)

#-------------------------------------------------------------------
if __name__ == '__main__':
	if len(sys.argv) == 1:
		teacher_defines_a_problem()
	elif len(sys.argv) == 2:
		teacher_defines_a_problem(sys.argv[1])
	else:
		print('\n\tUsage: python {} [some answer] \n'.format(sys.argv[0]))


