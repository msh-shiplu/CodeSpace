To use this software to share code in class, you will need to (1) install [Sublime Text](https://www.sublimetext.com/download) and (2) install a specific plug in for Sublime Text.


### Student's installation

(1) Open Sublime Text

(2) Click Show Console in the View menu.

(3) Copy this code:
```
import os; package_path = os.path.join(sublime.packages_path(), "GEMStudent"); os.mkdir(package_path) if not os.path.isdir(package_path) else print("dir exists"); module_file = os.path.join(package_path, "GEMStudent.py") ; menu_file = os.path.join(package_path, "Main.sublime-menu"); version_file = os.path.join(package_path, "version.go"); import urllib.request; urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/GEM/master/src/GEMStudent/GEMStudent.py", module_file); urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/GEM/master/src/GEMStudent/Main.sublime-menu", menu_file); urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/GEM/master/src/version.go", version_file)
```

(4) Paste copied code to Console and hit enter.

(5) In Sublime Text: (i) specify a folder on their computers to store local files, (ii) set the server address, which is shown when the server is run, and (iii) complete the registration by simply entering your given username.


### Teacher's installation

(1) Open Sublime Text

(2) Click Show Console in the View menu.

(3) Copy this code:
```
import os; package_path = os.path.join(sublime.packages_path(), "GEMTeacher"); os.mkdir(package_path) if not os.path.isdir(package_path) else print("dir exists"); module_file = os.path.join(package_path, "GEMTeacher.py") ; menu_file = os.path.join(package_path, "Main.sublime-menu"); version_file = os.path.join(package_path, "version.go"); import urllib.request; urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/GEM/master/src/GEMTeacher/GEMTeacher.py", module_file); urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/GEM/master/src/GEMTeacher/Main.sublime-menu", menu_file); urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/GEM/master/src/version.go", version_file); keymap_file = os.path.join(package_path, "Default.sublime-keymap"); urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/GEM/master/src/GEMTeacher/Default.sublime-keymap", keymap_file); 
```
(4) Paste copied code to Console and hit enter.

#### Running the GEM server

The instructor must run the GEM server.  The server can be run permanently or each time class starts.

Method 2 (deployment mode): Download the latest server ([Windows](https://www.dropbox.com/s/bjb8fvikjze20bu/gem.exe?dl=0), [MacOS](https://www.dropbox.com/s/vo3zn6pz8mhp083/gem?dl=0)) and make them *executable* on teacher's computer.  This command-line server needs to be run on the teacher's computer every time GEM is used in class.

(6) First-time configuration:
Add teachers
```
    ./gem -c config.json -add_teachers teachers.txt
```

Add students
```
    ./gem -c config.json -add_students students.txt
```

Run the server
```
    ./gem -c config.json
```

Examples of files containing configurations, teachers and students: [config.json](Examples/gem_config.json), 
[teachers.txt](Examples/teachers.txt), [students.txt](Examples/students.txt)

(7) When the server is run for the first time after teachers and students are added, teachers and students must configure their Sublime Text modules by going through 3 steps in Sublime Text: (i) specify a folder on their computers to store local files, (ii) set the server address, which is shown when the server is run, and (iii) complete the registration by simply entering their usernames, as specify in *teachers.txt* and *students.txt*.

These steps are done only once.  In subsequent usage, there is no need to go through these steps (even though the teacher's computer has a new IP address.)

#### Development mode

Install latest version of [Go](https://golang.org/dl/). Run these on the command line inside `src` folder.
* go mod init GEM
* go mod tidy
* go get github.com/mattn/go-sqlite3


(6) First-time configuration:
Add teachers
```
    ./go run *.go -c config.json -add_teachers teachers.txt
```

Add students
```
    ./go run *.go -c config.json -add_students students.txt
```

Run the server
```
    ./go run *.go -c config.json
```

### TA's installation

(1) Open Sublime Text

(2) Click Show Console in the View menu.

(3) Copy this code:
```
import os; package_path = os.path.join(sublime.packages_path(), "GEMAssistant"); os.mkdir(package_path) if not os.path.isdir(package_path) else print("dir exists"); module_file = os.path.join(package_path, "GEMAssistant.py") ; menu_file = os.path.join(package_path, "Main.sublime-menu"); version_file = os.path.join(package_path, "version.go"); import urllib.request; urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/GEM/master/src/GEMAssistant/GEMAssistant.py", module_file); urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/GEM/master/src/GEMAssistant/Main.sublime-menu", menu_file); urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/GEM/master/src/version.go", version_file); keymap_file = os.path.join(package_path, "Default.sublime-keymap"); urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/GEM/master/src/GEMAssistant/Default.sublime-keymap", keymap_file); 

```

(4) Paste copied code to Console and hit enter.

(5) in Sublime Text: (i) specify a folder on their computers to store local files, (ii) set the server address, which is shown when the server is run, and (iii) complete the registration by simply entering your given username.


