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

For Java/C++/Go files, you can use "//" instead of "#".

### Automatic grading

In the example below, a correct answer (D) will be automatically graded correct by GEM.  However, if the answer is incorrect, the instructor will have to look at the student's submission and grades it.

```
# 5 2 1 Cardinality

What is the cardinality of this set: {0, 1, 2, 3}?

A. 1
B. 2
C. 3
D. 4

ANSWER: D
```

Here's an example, where GEM can automatically grades correct and incorrect answers without borthering the instructor.  The difference in this example is the keyword *multiple_choice* that comes after the set of 3 numbers (specifying merit, effort, attempt).  This keyword indicates that the answer should be simple enough that an answer that is different from the correct answer (D) should be automatically graded incorrect.

```
# 5 2 1 multiple_choice Cardinality

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







