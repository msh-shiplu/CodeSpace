To use this software to share code in class, you will need to (1) install [Sublime Text 3](https://www.sublimetext.com/3) and (2) install a specific plug in for Sublime Text.

To install Sublime Text 3, [go here.](https://www.sublimetext.com/3)

### Teacher's installation

+ Open Sublime Text
+ Click Show Console in the View menu.
+ Copy this code:
```
import os; package_path = os.path.join(sublime.packages_path(), "GEMTeacher"); os.mkdir(package_path) if not os.path.isdir(package_path) else print("dir exists"); module_file = os.path.join(package_path, "GEMTeacher.py") ; menu_file = os.path.join(package_path, "Main.sublime-menu"); version_file = os.path.join(package_path, "version.go"); import urllib.request; urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/GEM/master/src/GEMTeacher/GEMTeacher.py", module_file); urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/GEM/master/src/GEMTeacher/Main.sublime-menu", menu_file); urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/GEM/master/src/version.go", version_file)
```
+ Paste copied code to Console and hit enter.

#### Add teachers

```
    ./go run *.go -config config.json -add_teacher teachers.txt
```

#### Add students

```
    ./go run *.go -config config.json -add_student students.txt
```

#### Run the server in class

```
    ./go run *.go -config config.json
```

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
