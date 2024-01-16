To use this software to share code in class, you will need to (1) install [Sublime Text](https://www.sublimetext.com/download) and (2) install a specific plug in for Sublime Text.


### Student's installation

(1) Open Sublime Text

(2) Click Show Console in the View menu.

(3) Copy this code:
```
import os; package_path = os.path.join(sublime.packages_path(), "GEMStudent"); os.mkdir(package_path) if not os.path.isdir(package_path) else print("dir exists"); module_file = os.path.join(package_path, "GEMStudent.py") ; menu_file = os.path.join(package_path, "Main.sublime-menu"); version_file = os.path.join(package_path, "version.go"); import urllib.request; urllib.request.urlretrieve("https://raw.githubusercontent.com/msh-shiplu/CodeSpace/master/src/GEMStudent/GEMStudent.py", module_file); urllib.request.urlretrieve("https://raw.githubusercontent.com/msh-shiplu/CodeSpace/master/src/GEMStudent/Main.sublime-menu", menu_file); # urllib.request.urlretrieve("https://raw.githubusercontent.com/msh-shiplu/CodeSpace/src/version.go", version_file)
```

(4) Paste copied code to Console and hit enter.

(5) In Sublime Text: (i) specify a folder on their computers to store local files, (ii) set the server address, which is shown when the server is run, and (iii) complete the registration by simply entering your given username.


### Teacher's installation
Web portal
```
    http://server_address:course_port/
```
For example
```
    http://192.168.86.189:8088/
```

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

##### Install Go
Install latest version of [Go](https://golang.org/dl/). Run these on the command line inside `src` folder.
* go mod init GEM
* go mod tidy

##### Install MySQL Server
The current version of Codespace need MySQL server pre-installed and running beforehand. The MySQL server can be installed using the following resources.
* [Debian (version 11)](https://www.devart.com/dbforge/mysql/install-mysql-on-debian/).
* [Ubuntu](https://ubuntu.com/server/docs/databases-mysql)
* [Windows](https://dev.mysql.com/doc/refman/8.2/en/windows-installation.html)

##### Create New MySQL User and Grant Previlege
Enter into MySQL CLI using the following command
```
myql -u root -p
```
Create a new user e.g. `gem`
```
CREATE USER 'gem'@'localhost' IDENTIFIED BY 'password';
```
Grant the user previleges for all databases (the system will create a database for each of the course).
```
GRANT ALL PRIVILEGES ON *.* TO 'gem'@'localhost' WITH GRANT OPTION;
```

##### Install MySQL Driver for Go
```
go install github.com/go-sql-driver/mysql@latest
```


##### First-time configuration:
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
#### Run [ChatGPTA](https://github.com/vtphan/chatGPTA/tree/main) Server (Optional)
To get feedback from ChatGPT, we need to run the ChatGPTA server independently and provide the server address to the config file.
Follow the [instructions](https://github.com/vtphan/chatGPTA/blob/main/README.md) in the repository to install the ChatGPTA server.