To use this software to share code in class, you will need to (1) install [Sublime Text 3](https://www.sublimetext.com/3) and (2) install a specific plug in for Sublime Text.

To install Sublime Text 3, [go here.](https://www.sublimetext.com/3)

### Teacher's installation

(1) Open Sublime Text

(2) Click Show Console in the View menu.

(3) Copy this code:
```
import os; package_path = os.path.join(sublime.packages_path(), "GEMTeacher"); os.mkdir(package_path) if not os.path.isdir(package_path) else print("dir exists"); module_file = os.path.join(package_path, "GEMTeacher.py") ; menu_file = os.path.join(package_path, "Main.sublime-menu"); version_file = os.path.join(package_path, "version.go"); import urllib.request; urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/GEM/master/src/GEMTeacher/GEMTeacher.py", module_file); urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/GEM/master/src/GEMTeacher/Main.sublime-menu", menu_file); urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/GEM/master/src/version.go", version_file)
```
(4) Paste copied code to Console and hit enter.

(5) Download the latest server ([Windows](http://umdrive.memphis.edu/vphan/public/GEM/gem.exe), [MacOS](http://umdrive.memphis.edu/vphan/public/GEM/gem)) and make them *executable* on teacher's computer.  This command-line server needs to be run on the teacher's computer every time GEM is used in class.

(6) First-time configuration:
Add teachers
```
    ./gem -c config.json -add_teacher teachers.txt
```

Add students
```
    ./gem -c config.json -add_student students.txt
```

Run the server
```
    ./gem -c config.json
```

(7) When the server is run for the first time after teachers and students are added, teachers and students must configure their Sublime Text modules by going through 3 steps in Sublime Text: (i) specify a local folder on their computers, (ii) set the server address, which is shown when the server is run, and (iii) complete the registration by simply entering their usernames, as specify in *teachers.txt* and *students.txt*.


### TA's installation

+ Open Sublime Text
+ Click Show Console in the View menu.
+ Copy this code:
```
import os; package_path = os.path.join(sublime.packages_path(), "GEMAssistant"); os.mkdir(package_path) if not os.path.isdir(package_path) else print("dir exists"); module_file = os.path.join(package_path, "GEMAssistant.py") ; menu_file = os.path.join(package_path, "Main.sublime-menu"); version_file = os.path.join(package_path, "version.go"); import urllib.request; urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/GEM/master/src/GEMAssistant/GEMAssistant.py", module_file); urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/GEM/master/src/GEMAssistant/Main.sublime-menu", menu_file); urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/GEM/master/src/version.go", version_file)
```
+ Paste copied code to Console and hit enter.


### Student's installation

+ Open Sublime Text
+ Click Show Console in the View menu.
+ Copy this code:
```
import os; package_path = os.path.join(sublime.packages_path(), "GEMStudent"); os.mkdir(package_path) if not os.path.isdir(package_path) else print("dir exists"); module_file = os.path.join(package_path, "GEMStudent.py") ; menu_file = os.path.join(package_path, "Main.sublime-menu"); version_file = os.path.join(package_path, "version.go"); import urllib.request; urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/GEM/master/src/GEMStudent/GEMStudent.py", module_file); urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/GEM/master/src/GEMStudent/Main.sublime-menu", menu_file); urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/GEM/master/src/version.go", version_file)
```
+ Paste copied code to Console and hit enter.
