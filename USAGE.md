This page explains how to use GEM, as a teacher, in the class room.

The basic classroom workflow includes: (i) the teacher broadcasts a file, which is opened in the current tab; (ii) students get the file from their virtual white boards, work on them and submit their work; (iii) if the work is gradable, the teacher grades and possibly sends feedback to students.

### Teacher broadcasting a file.

Click on "Broadcast" on the GEM menu in Sublime Text.  This will send whichever file is open in the current tab.

If the file is not *gradable*, students can get it from their white board, but they can't submit for grading.

If the file is gradable, students can submit their work on it.  To make a file *gradable*, the first line of the file must satisfy a specific syntax. Here's an example:

```
# 5 2 3 Loop and iteration
```

This means the problem has 5 points for merit (for correct solution), 2 points for effort (for trying at least once), and students have 3 attempts to try.  Further, the learning objective for this problem is "Loop and iteration".

For Java/C++/Go files, we should use "//" instead of "#" to by syntatically correct.

### Automatic grading

If a problem has an answer (i.e. the last line has the "ANSWER:" tag followed by a 1-line answer), student submissions will be graded automatically.  If a student's answer matches the specified answer exactly, the solution is automatically marked correct. If they don't match, the solution is marked incorrect.

```
# 5 2 1 Cardinality

What is the cardinality of this set: {0, 1, 2, 3}?

A. 1
B. 2
C. 3
D. 4

ANSWER: D
```

There are cases where teachers want to look at solutions whose answers do not exactly match the specified one.  If such is the case, the teacher can specify the keyword *_manual_* that comes after the set of 3 numbers (specifying merit, effort, attempt).  This keyword indicates that students' answers that are not exactly matched will be looked at by the instructor, instead of being marked "incorrect" automatically.

Students' answers that are matched exactly are still automatically marked as "correct".

```
# 5 2 1 _manual_ Cardinality

What is the cardinality of this set: {0, 1, 2, 3}?

A. 1
B. 2
C. 3
D. 4

ANSWER: D
```

### Getting a student's submission and Grading

By selecting an appropriate item in the GEM menu, the teacher get a student's submission and then can give a grading of "correct" or "incorrect" to it.  The teacher can also "dismiss" the submission without grading it.

### Time's Up

This feature disallows students to submit answers to the current problem.  If the problem has answers (see above), GEM will automatically summarizes the answers.

### Excerpting selected content to share on the bulletin board

Selecting a specific content, the teacher can share this content to the "bulletin board" so all students can see.  This feature is particular useful for teaching assistants, who help the teacher with grading but do not have their computers connected to the projector.

### Keyboard shortcuts 

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



