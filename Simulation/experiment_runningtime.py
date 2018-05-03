from teacher_defines_a_problem import teacher_defines_a_problem
from students_solve import students_solve
from teacher_grades import teacher_grades
from common import *

content = '''
# 5 points for correctness, 1 point for effort. Maximum attempts: 1 

def foo(L):
	sum = 0
	for x in L:
		for y in L:
			sum = x * y
	return sum

Select one of these as the running time of foo:

A 	O(n)
B 	O(n^2)
C 	O(n log n)
D 	O(1)
'''

objective = 'Analyze running time of nested loops.'

teacher_defines_a_problem(content=content, tag=objective, merit=5, effort=1, attempts=1)
students_solve('ABCD')
teacher_grades()
