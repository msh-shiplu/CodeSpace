import requests
import os
import sys
import json
import random
from common import *

#-------------------------------------------------------------------
def register(filename, role):
	students = []
	with open(filename) as f:
		users = [ l.strip() for l in f.readlines() if l.strip() ]
	a = users[0]
	for user in users:
		res = request('complete_registration', {'name':user, 'role':role})
		uid, pwd, cid, nameserver = res.split(',')
		info = dict(
			CourseId = cid,
			NameServer = nameserver,
			Uid = int(uid),
			Password = pwd,
			Name = user,
			Server = Config['Server'],
			Ability = random.random(),
		)
		if role == 'teacher':
			outfile = os.path.join(TeacherDir, 'info_{}.json'.format(uid))
		elif role == 'student':
			outfile = os.path.join(StudentDir, 'info_{}.json'.format(uid))
		else:
			raise Exception('Unknown role.')
		with open(outfile, 'w') as f:
			f.write(json.dumps(info, sort_keys=True, indent=4))
		print('Saving info of {} to {}'.format(user, outfile))

#-------------------------------------------------------------------
if __name__ == '__main__':
	if len(sys.argv) != 3:
		print('\n\tUsage: python {} user_file.txt [student|teacher]\n'.format(sys.argv[0]))
		sys.exit(0)
	register(sys.argv[1], sys.argv[2])