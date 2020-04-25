This page explains how to use GEM, as a teacher, in the class room.

The basic classroom workflow includes: (i) the teacher broadcasts a file, which is opened in the current tab; (ii) students get the file from their virtual white boards, work on them and submit their work; (iii) if the work is gradable, the teacher grades and possibly sends feedback to students.

### Teacher broadcasting a file.

Click on "Broadcast" on the GEM menu in Sublime Text.  This will send whichever file is open in the current tab.

If the file is not *gradable*, students can get it from their white board, but they can't submit for grading.

If the file is gradable, students can submit their work on it.  To make a file *gradable*, the first line of the file must satisfy a specific syntax. Here's an example:

```
# 5 1 2 simple recursion
```

This means the problem has 5 points for a correct grade, 1 point an incorrect grade, and students have 2 attempts to try.  Further, the learning objective for this problem is "simple recursion".

For Java/C++/Go files, we should use "//" instead of "#" to by syntatically correct:

```
// 5 1 2 simple recursion
```

### Grading and Partial Credits

When a submission is marked correct, 5 points are automatically recorded for the submission.

When a submission is marked inccorect, 1 points are automatically recorded for the submission.

The instructor/TA can revise points for effort by modyfing graded submission.  Specifically, the first
line of a submission looks like this:

```
# 5 points, 1 for effort. Maximum attempts: 2.
```

If instructors and TAs can revise points for effort if they see fit. For example, for instance, to give 3 points for effort (or partial credits), instructors and TA can replace 1 with 3 the first line of the submission as follows:

```
# 5 points, 3 for effort. Maximum attempts: 2.
```

Here is a concrete example: [exercise1.py](Examples/exercise1.py)

### Multiple-choice questions and automatic grading

A multiple choice question is automatically graded by default. The last line of the file should be "ANSWER: ", where students can give their answer. There must be an answer file associated with each question.

Example: Suppose the file [exercise2.txt](Examples/exercise2.txt) has this content:


```
# 5 2 1 Cardinality

What is the cardinality of this set: {0, 1, 2, 3}?

A. 1
B. 2
C. 3
D. 4

ANSWER: 
```

There should be another file called [exercise2.txt.answer](Examples/exercise2.txt.answer) with the following content:

```
D
```

In this case student submissions are automatically submitted. Answers are case sensitive. 


There are cases where teachers want to look at solutions whose answers do not exactly match the specified one.  If such is the case, the teacher can specify the keyword *_manual_* that comes after the set of 3 numbers (specifying merit, effort, attempt).  This keyword indicates that students' answers that are not exactly matched will be looked at by the instructor, instead of being marked "incorrect" automatically.

Students' answers that are matched exactly are still automatically marked as "correct".

```
# 5 2 1 _manual_ Cardinality

What is the cardinality of this set: {0, 1, 2, 3}?

A. 1
B. 2
C. 3
D. 4

ANSWER:
```

### Getting a student's submission and Grading

By selecting an appropriate item in the GEM menu, the teacher get a student's submission and then can give a grading of "correct" or "incorrect" to it.  The teacher can also "dismiss" the submission without grading it.

### Time's Up

Selecting "time's up" in the active tab that has a problem will disallow students to submit answers to the problem.  If the problem has answers (see above), GEM will automatically summarizes the answers.

### Excerpting selected content to share on the bulletin board

Selecting a specific content, the teacher can share this content to the "bulletin board" so all students can see.  This feature is particular useful for teaching assistants, who help the teacher with grading but do not have their computers connected to the projector.

### Keyboard shortcuts for teachers

In Sublime Text, select Preferences --> Key Bindings. This will open two files: (1) Default/Default.sublime-keymap and (2) User/Default.sublime-keymap.

Add (and save) the following to the file User/Default.sublime-keymap:

```
[
    { "keys": ["ctrl+1"], "command": "gemt_get_prioritized" },
    { "keys": ["ctrl+2"], "command": "gemt_grade_correct" },
    { "keys": ["ctrl+3"], "command": "gemt_grade_incorrect" },
]
```

The first (Ctrl+1) is a shortcut for getting a submision (prioritized those who need help).

The second and third (Ctrl+2 and Ctrl+3) are shortcuts for marking a submission "correct" or "incorrect".

These shortcuts can be customized, but make sure they do not "override" existing shortcuts.


### Keyboard shortcuts for teaching assistants

In Sublime Text, select Preferences --> Key Bindings. This will open two files: (1) Default/Default.sublime-keymap and (2) User/Default.sublime-keymap.

Add (and save) the following to the file User/Default.sublime-keymap:

```
[
    { "keys": ["ctrl+4"], "command": "gema_get_prioritized" },
    { "keys": ["ctrl+5"], "command": "gema_grade_correct" },
    { "keys": ["ctrl+6"], "command": "gema_grade_incorrect" },
]
```

The first (Ctrl+4) is a shortcut for getting a submision (prioritized those who need help).

The second and third (Ctrl+5 and Ctrl+6) are shortcuts for marking a submission "correct" or "incorrect".

These shortcuts can be customized, but make sure they do not "override" existing shortcuts.


