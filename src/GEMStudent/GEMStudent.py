# GEMStudent
# Author: Vinhthuy Phan, 2018
#
import imp
from urllib import response
import sublime
import sublime_plugin
import urllib.parse
import urllib.request
import os
import json
import time
import datetime
import webbrowser
import threading
import re

gemsUpdateIntervalLong = 20000  # Update interval
gemsUpdateIntervalShort = 10000  # When submission is being looked at
gemsCodeSnapshotInterval = 15000
gemsFILE = os.path.join(os.path.dirname(os.path.realpath(__file__)), "info")
gemsFOLDER = ""
gemsTIMEOUT = 7
gemsAnswerTag = "ANSWER:"
gemsTracking = False
gemsSERVER = ""
gemsSERVER_TIME = 0
gemsConnected = False
gemsUpdateMessage = {
    1: "Your {} submission for problem {} is being looked at.",
    2: "Teacher did not grade your {} submission for problem {}.",
    3: "Good effort!!!  However, the teacher did not think your {} solution for problem {} was correct.",
    4: "Your {} solution for problem {} was correct.",
    5: "Your {} solution for problem {} was correct. You are now elligible to help your classmates.",
}
gemsFeedbackFolder = "FEEDBACKS"
gemsCurrentHelpSubId = None
gemsHelpRequestMessage = [
    "You have fetched a help request entry.",
    "There is no pending help request.",
    "You are not yet elligible to help",
]
gemsCurrentFiles = set()

gemsSnapshotTracking = False
gemsActiveFiles = {}
lastSentCodes = {}
isRegistered = False

# Set of filename: to check whether the student got a feedback against a problem
gotFeedback = {}
# Set of filename: to check whether the student already submitted for this problem.
submitted = set()
# Stores the submitted and gotFeedback in case sublime text exits
gemsSubFile = os.path.join(os.path.dirname(os.path.realpath(__file__)), "sub.p")
gemsBackFeedbackTimeout = 60 * 10  # 10 minutes
gemsBackFeedbackStatus = {}
gemsBackFeedbackTimers = {}

# ------------------------------------------------------------------


# def plugin_loaded():
#     updateActiveProblems()
#     sendCodeSnapshot()


def sendCodeSnapshot():
    global gemsSnapshotTracking
    global gemsActiveFiles
    global lastSentCodes
    global gemsTracking
    try:
        for window in sublime.windows():
            for view in window.views():
                filename = view.file_name()
                if filename in gemsActiveFiles:
                    problem_id = gemsActiveFiles[filename]
                    code = view.substr(sublime.Region(0, view.size())).rstrip()
                    if (
                        problem_id in lastSentCodes
                        and lastSentCodes[problem_id] == code
                    ):
                        continue
                    data = {
                        "code": code,
                        "problem_id": problem_id,
                        "event": "at_regular_interval",
                    }
                    gemsRequest("code_snapshot", data)
                    lastSentCodes[problem_id] = code
                    print("Code snapshot sent!")

                    if gemsTracking == False:
                        gemsTracking = True
                        sublime.set_timeout_async(gems_periodic_update, 5000)
    except:
        gemsSnapshotTracking = False
        return

    if gemsSnapshotTracking:
        sublime.set_timeout_async(sendCodeSnapshot, gemsCodeSnapshotInterval)


def updateActiveProblems():
    global gemsSnapshotTracking
    global gemsActiveFiles
    try:
        response = gemsRequest("get_global_info", {})
        response = json.loads(response)
        gemsActiveFiles = dict()
        for problem in response["ActiveProblems"]:
            gemsActiveFiles[os.path.join(gemsFOLDER, problem["Filename"])] = problem[
                "ProblemID"
            ]
    except:
        gemsSnapshotTracking = False
    if gemsSnapshotTracking:
        sublime.set_timeout_async(updateActiveProblems, 1000 * 60)


# ------------------------------------------------------------------


class gemsReviewFeedbackForMe(sublime_plugin.ApplicationCommand):
    def run(self):
        global gemsSERVER
        with open(gemsFILE, "r") as f:
            info = json.loads(f.read())
        p = urllib.parse.urlencode(
            {
                "viewtype": "forme",
                "password": info["Password"],
                "uid": info["Uid"],
                "role": "student",
            }
        )
        webbrowser.open(gemsSERVER + "/student_views_feedback?" + p)


# ------------------------------------------------------------------


class gemsGetCodeSpace(sublime_plugin.ApplicationCommand):
    def run(self):
        global gemsSERVER
        with open(gemsFILE, "r") as f:
            info = json.loads(f.read())
        p = urllib.parse.urlencode(
            {"password": info["Password"], "uid": info["Uid"], "role": "student"}
        )
        webbrowser.open(gemsSERVER + "/get_codespace?" + p)


# ------------------------------------------------------------------


class gemsAttendanceReport(sublime_plugin.ApplicationCommand):
    def run(self):
        response = gemsRequest("student_checks_in", {})
        json_obj = json.loads(response)
        dates = set()
        for d in json_obj:
            dates.add(datetime.datetime.fromtimestamp(d).strftime("%Y-%m-%d"))

        with open(gemsFILE, "r") as f:
            info = json.loads(f.read())
        report_file = os.path.join(info["Folder"], "Attendance.txt")
        with open(report_file, "w", encoding="utf-8") as f:
            f.write("Your attendance was taken on these dates:\n")
            for d in sorted(dates, reverse=True):
                f.write("{}\n".format(d))
        sublime.run_command("new_window")
        sublime.active_window().open_file(report_file)


# ------------------------------------------------------------------


class gemsPointsReport(sublime_plugin.ApplicationCommand):
    def run(self):
        response = gemsRequest("student_gets_report", {})
        if response != None:
            json_obj = json.loads(response)
            report = {}
            total_points = 0
            for i in json_obj:
                Date, Points, Filename = i["Date"], i["Points"], i["Filename"]
                if Date not in report:
                    report[Date] = []
                if "." in Filename:
                    prefix, ext = Filename.rsplit(".", 1)
                else:
                    prefix = Filename
                report[Date].append((Points, prefix))
                total_points += Points

            with open(gemsFILE, "r") as f:
                info = json.loads(f.read())
            report_file = os.path.join(info["Folder"], "Points.txt")
            with open(report_file, "w", encoding="utf-8") as f:
                f.write("Total points: {}\n".format(total_points))
                for d, v in reversed(sorted(report.items())):
                    date = datetime.datetime.fromtimestamp(d).strftime("%Y-%m-%d")
                    for entry in v:
                        f.write("{}\t{}\t{}\n".format(date, entry[0], entry[1]))

            sublime.run_command("new_window")
            sublime.active_window().open_file(report_file)


# ------------------------------------------------------------------


def get_word_from_num(n):
    if n > 10 and n < 20:
        return str(n) + "th"
    if n % 10 == 1:
        return str(n) + "st"
    if n % 10 == 2:
        return str(n) + "nd"
    if n % 10 == 3:
        return str(n) + "rd"
    return str(n) + "th"


def gems_periodic_update():
    global gemsTracking
    response = gemsRequest("student_periodic_update", {}, verbal=False)
    if response is None:
        print("Response is None. Stop tracking.")
        gemsTracking = False
        return
    try:
        (
            submission_stat,
            board_stat,
            thank_stat,
            snapshot_feedback_stat,
            attempt_number,
            filename,
        ) = response.split(";")
        submission_stat = int(submission_stat)
        board_stat = int(board_stat)
        thank_stat = int(thank_stat)
        attempt_number = int(attempt_number)
        snapshot_feedback_stat = int(snapshot_feedback_stat)
        # Display messages if necessary
        mesg = ""
        if submission_stat > 0 and submission_stat in gemsUpdateMessage:
            mesg = gemsUpdateMessage[submission_stat].format(
                get_word_from_num(attempt_number), filename[: filename.rfind(".")]
            )
        if thank_stat == 1:
            mesg += "\nYour classmate Thanked You for your help."
        if board_stat == 1:
            mesg += "\nTeacher placed new material on your board."
        elif board_stat == 2:
            mesg += "\nYou have feedback on your board."
        mesg = mesg.strip()
        if mesg != "":
            sublime.message_dialog(mesg)

        if snapshot_feedback_stat == 1:
            resp = gemsRequest("get_snapshot_feedback", {})
            if resp is not None:
                resp = json.loads(resp)
                sublime.message_dialog(
                    "You recieved feedback from " + resp["Provider"] + "."
                )
                feedbackFolder = os.path.join(gemsFOLDER, gemsFeedbackFolder)
                if not os.path.exists(feedbackFolder):
                    os.makedirs(feedbackFolder)
                filename = os.path.join(
                    feedbackFolder,
                    "feedback-" + str(resp["FeedbackID"]) + "-" + resp["ProblemName"],
                )
                if resp["ProblemName"].endswith(".py"):
                    comment = "#"
                else:
                    comment = "//"
                with open(filename, "w", encoding="utf-8") as f:
                    f.write(resp["Feedback"])
                if sublime.active_window().id() == 0:
                    sublime.run_command("new_window")
                view = sublime.active_window().open_file(filename)
                # view.set_read_only(True)
                print("Code snapshot feedback recieved!")

        # Open board pages and feedback automatically
        # if board_stat == 1:
        # 	if sublime.active_window().id() == 0:
        # 		sublime.run_command('new_window')
        # 	sublime.active_window().run_command('gems_get_board_content')

        # Keep checking periodically
        update_timeout = gemsUpdateIntervalLong
        if submission_stat == 1:
            update_timeout = gemsUpdateIntervalShort
        # print('checking', submission_stat, board_stat, update_timeout)
        sublime.set_timeout_async(gems_periodic_update, update_timeout)
    except:
        gemsTracking = False


# ------------------------------------------------------------------


def gems_share(self, edit, priority):
    global gemsTracking
    fname = self.view.file_name()
    if fname is None:
        sublime.message_dialog("Cannot share unsaved content.")
        return
    content = self.view.substr(sublime.Region(0, self.view.size())).lstrip()
    filename = os.path.basename(fname)
    match = re.search(r"feedback-\d+-(.*)", filename)
    if match:
        filename = match.group(1)
    items = content.rsplit(gemsAnswerTag, 1)
    if len(items) == 2:
        answer = items[1].strip()
    else:
        answer = ""

    if priority == 1 and os.path.exists(fname + ".test_cases"):
        with open(fname + ".test_cases", "r") as f:
            testCases = f.read()
    else:
        testCases = ""

    data = dict(
        content=content,
        answer=answer,
        testcases=testCases,
        filename=filename,
        priority=priority,
    )
    response = gemsRequest("student_shares", data)
    sublime.message_dialog(response)
    # if priority == 1:
    #     if filename in gotFeedback:
    #         for feedback_filename in gotFeedback[filename]:
    #             ask_for_back_feedback(filename, feedback_filename)

    #     with open(gemsSubFile, "wb+") as f:
    #         pickle.dump({"gotFeedback": gotFeedback,
    #                     "submitted": submitted}, f)
    if gemsTracking == False:
        gemsTracking = True
        sublime.set_timeout_async(gems_periodic_update, 5000)


# ------------------------------------------------------------------


class gemsNeedHelp(sublime_plugin.TextCommand):
    def run(self, edit):
        # gems_share(self, edit, priority=2)
        fname = self.view.file_name()
        if fname is None:
            sublime.message_dialog("Cannot share unsaved content.")
            return
        # if os.path.basename(fname) not in gemsCurrentFiles:
        # 	sublime.message_dialog("Invalid file")
        # 	return
        sublime.message_dialog(
            "Walk me through your thought process for what you are trying to accomplish."
        )
        sublime.active_window().show_input_panel(
            "Click Enter to send help request:",
            "",
            self.get_need_help_with,
            None,
            self.get_need_help_with,
        )

    # def get_trying_what(self, trying_what_message=""):
    # 	self.trying_what_message = trying_what_message
    # 	sublime.active_window().show_input_panel("Explain what you need help with:", "", self.get_need_help_with, None, self.get_need_help_with)

    def get_need_help_with(self, need_help_with_message=""):
        global gemsTracking
        fname = self.view.file_name()
        content = self.view.substr(sublime.Region(0, self.view.size())).lstrip()

        filename = os.path.basename(fname)
        match = re.search(r"feedback-\d+-(.*)", filename)
        if match:
            filename = match.group(1)
        data = dict(
            content=content,
            filename=filename,
            # trying_what=self.trying_what_message,
            need_help_with=need_help_with_message,
        )
        response = gemsRequest("student_ask_help", data)
        sublime.message_dialog(response)
        if gemsTracking == False:
            gemsTracking = True
            sublime.set_timeout_async(gems_periodic_update, 5000)


# ------------------------------------------------------------------


class gemsGotIt(sublime_plugin.TextCommand):
    def run(self, edit):
        gems_share(self, edit, priority=1)


# ------------------------------------------------------------------


class gemsGetBoardContent(sublime_plugin.ApplicationCommand):
    def run(self):
        response = gemsRequest("student_gets", {})
        if response is None:
            return
        json_obj = json.loads(response)
        if json_obj == []:
            sublime.message_dialog("Whiteboard is empty.")
            return

        feedback_dir = os.path.join(gemsFOLDER, "FEEDBACK")
        if not os.path.exists(feedback_dir):
            os.mkdir(feedback_dir)
        old_dir = os.path.join(gemsFOLDER, "OLD")
        if not os.path.exists(old_dir):
            os.mkdir(old_dir)
        # start sending code snapshots if not started yet or turned off
        global gemsSnapshotTracking
        if gemsSnapshotTracking == False:
            gemsSnapshotTracking = True
            updateActiveProblems()
            sendCodeSnapshot()

        for board in json_obj:
            content = board["Content"]
            filename = board["Filename"]
            self.filename = filename
            mesg = ""
            if board["Type"] in ["feedback", "peer_feedback"]:
                global gemsBackFeedbackTimers
                global gemsBackFeedbackStatus
                problem_filename = filename[filename.find("-", 9) + 1 :]
                local_file = os.path.join(feedback_dir, filename)
                mesg = "You have feedback"
                # if problem_filename not in gotFeedback:
                #     gotFeedback[problem_filename] = []
                # gotFeedback[problem_filename].append(local_file)
                gemsBackFeedbackTimers[
                    (problem_filename, local_file)
                ] = threading.Timer(
                    gemsBackFeedbackTimeout,
                    ask_for_back_feedback,
                    [problem_filename, local_file],
                )
                gemsBackFeedbackTimers[(problem_filename, local_file)].start()
                gemsBackFeedbackStatus[(problem_filename, local_file)] = True
                # with open(gemsSubFile, "wb+") as f:
                #     pickle.dump({"gotFeedback": gotFeedback,
                #                 "submitted": submitted}, f)
            else:
                local_file = os.path.join(gemsFOLDER, filename)
                if os.path.exists(local_file):
                    with open(local_file) as f:
                        moved_file = os.path.join(old_dir, filename)
                        with open(moved_file, "w", encoding="utf-8") as newf:
                            newf.write(f.read())
                    mesg = "Move existing file, {}, to {}.".format(filename, old_dir)
            with open(local_file, "w", encoding="utf-8") as f:
                f.write(content)
            if sublime.active_window().id() == 0:
                sublime.run_command("new_window")
            sublime.active_window().open_file(local_file)
            if board["Type"] == "peer_feedback":
                sublime.active_window().active_view().set_read_only(True)
            # elif board['Type'] == 'new':
            # 	gemsCurrentFiles.add(filename)
            # if mesg != '':
            # 	sublime.message_dialog(mesg)

    def send_thank_you(self, message):
        message = message.lower().strip()
        if message != "yes" and message != "no":
            self.force_for_thank_you()
            return

        data = {"useful": message, "message_id": self.message_id}
        gemsRequest("student_send_thank_you", data)

    def force_for_thank_you(self):
        decision = sublime.yes_no_cancel_dialog(
            "Please confirm: was this useful feedback?", "Yes", "No"
        )
        if decision == sublime.DIALOG_YES:
            message = "yes"
        elif decision == sublime.DIALOG_NO:
            message = "no"
        elif decision == sublime.DIALOG_CANCEL:
            message = "cancel"
        else:
            self.force_for_thank_you()
            return

        data = {"useful": message, "message_id": self.message_id}
        gemsRequest("student_send_thank_you", data)


# ------------------------------------------------------------------------------
# ------------------------------------------------------------------------------
# These functionalities below are identical to those of teachers
# ------------------------------------------------------------------------------
# ------------------------------------------------------------------------------
def gemsRequest(path, data, authenticated=True, method="POST", verbal=True):
    global gemsFOLDER, gemsSERVER, gemsSERVER_TIME

    try:
        with open(gemsFILE, "r") as f:
            info = json.loads(f.read())
    except:
        info = dict()

    if "Folder" not in info:
        if verbal:
            sublime.message_dialog("Please set a local folder to store working files.")
        return None

    if "CourseId" not in info:
        sublime.message_dialog("Please set the course id.")
        return None

    if "Server" not in info:
        if verbal:
            sublime.message_dialog("Please connect to the server first.")
        return None

    if gemsSERVER == "" or time.time() - gemsSERVER_TIME > 5400:
        sublime.run_command("gems_connect")
        if gemsSERVER == "":
            sublime.message_dialog(
                "Unable to connect. Check server address or course id."
            )
            return

    if authenticated:
        if "Uid" not in info:
            sublime.message_dialog("Please register.")
            return None
        data["name"] = info["Name"]
        data["password"] = info["Password"]
        data["uid"] = info["Uid"]
        data["role"] = "student"
        gemsFOLDER = info["Folder"]

    url = urllib.parse.urljoin(gemsSERVER, path)
    load = urllib.parse.urlencode(data).encode("utf-8")
    req = urllib.request.Request(url, load, method=method)
    try:
        with urllib.request.urlopen(req, None, gemsTIMEOUT) as response:
            return response.read().decode(encoding="utf-8")
    except urllib.error.HTTPError as err:
        if verbal:
            sublime.message_dialog("{0}".format(err))
    except urllib.error.URLError as err:
        if verbal:
            sublime.message_dialog("{0}\nCannot connect to server.".format(err))
    print("Error making request")
    return None


# ------------------------------------------------------------------


class gemsSetLocalFolder(sublime_plugin.ApplicationCommand):
    def run(self):
        try:
            with open(gemsFILE, "r") as f:
                info = json.loads(f.read())
        except:
            info = dict()
        if "Folder" not in info:
            info["Folder"] = os.path.join(os.path.expanduser("~"), "GEM")
        if sublime.active_window().id() == 0:
            sublime.run_command("new_window")
        sublime.active_window().show_input_panel(
            "This folder will be used to store working files.",
            info["Folder"],
            self.set,
            None,
            None,
        )

    def set(self, folder):
        folder = folder.strip()
        if len(folder) > 0:
            try:
                with open(gemsFILE, "r") as f:
                    info = json.loads(f.read())
            except:
                info = dict()
            info["Folder"] = folder
            if not os.path.exists(folder):
                try:
                    os.mkdir(folder)
                    os.mkdir(os.path.join(folder, "FEEDBACK"))
                    with open(gemsFILE, "w") as f:
                        f.write(json.dumps(info, indent=4))
                except:
                    sublime.message_dialog("Could not create {}.".format(folder))
            else:
                with open(gemsFILE, "w") as f:
                    f.write(json.dumps(info, indent=4))
                sublime.message_dialog(
                    "Folder exists. Will use it to store working files."
                )
        else:
            sublime.message_dialog("Folder name cannot be empty.")


# ------------------------------------------------------------------


class gemsConnect(sublime_plugin.ApplicationCommand):
    def run(self):
        try:
            with open(gemsFILE, "r") as f:
                info = json.loads(f.read())
        except:
            info = dict()

        if "CourseId" not in info:
            sublime.message_dialog("Please set the course id.")
            return None

        if "Server" not in info:
            sublime.message_dialog("Please set server address.")
            return None

        global gemsSERVER, gemsSERVER_TIME
        url = urllib.parse.urljoin(info["Server"], "ask")
        load = urllib.parse.urlencode({"who": info["CourseId"]}).encode("utf-8")
        req = urllib.request.Request(url, load)
        try:
            with urllib.request.urlopen(req, None, gemsTIMEOUT) as response:
                server = response.read().decode(encoding="utf-8")
                try:
                    with open(gemsFILE, "r") as f:
                        info = json.loads(f.read())
                except:
                    info = dict()
                if not server.startswith("http://"):
                    sublime.message_dialog("Unable to get address.")
                    return
                gemsSERVER = server
                gemsSERVER_TIME = time.time()
                sublime.status_message("Connected")
        except urllib.error.HTTPError as err:
            sublime.message_dialog("{0}".format(err))
        except urllib.error.URLError as err:
            sublime.message_dialog("{0}\nCannot connect to server.".format(err))


# ------------------------------------------------------------------


class gemsCompleteRegistration(sublime_plugin.ApplicationCommand):
    def run(self):
        try:
            with open(gemsFILE, "r") as f:
                info = json.loads(f.read())
        except:
            info = dict()

        if "CourseId" not in info:
            sublime.message_dialog("Please enter course id.")
            return

        if "Name" not in info:
            sublime.message_dialog("Please enter assigned username.")
            return

        response = gemsRequest(
            "complete_registration",
            {"role": "student", "name": info["Name"], "course_id": info["CourseId"]},
            authenticated=False,
        )
        if response is None:
            sublime.message_dialog("Response is None. Failed to complete registration.")
            return

        if response == "Failed" or response.count(",") != 1:
            sublime.message_dialog("Failed to complete registration.")
        else:
            uid, password = response.split(",")
            info["Uid"] = int(uid)
            info["Password"] = password.strip()
            global isRegistered
            isRegistered = True
            sublime.message_dialog("{} is registered.".format(info["Name"]))
        with open(gemsFILE, "w") as f:
            f.write(json.dumps(info, indent=4))


# ------------------------------------------------------------------


class gemsSetServerAddress(sublime_plugin.ApplicationCommand):
    def run(self):
        try:
            with open(gemsFILE, "r") as f:
                info = json.loads(f.read())
        except:
            info = dict()

        if "Server" not in info:
            info["Server"] = ""
        if sublime.active_window().id() == 0:
            sublime.run_command("new_window")
        sublime.active_window().show_input_panel(
            "Set server address.  Press Enter:",
            info["Server"],
            self.set,
            None,
            self.on_cancel,
        )

    def on_cancel(self):
        print("line 417")
        try:
            with open(gemsFILE, "r") as f:
                info = json.loads(f.read())
        except:
            info = dict()

        if "Server" not in info:
            info["Server"] = ""
        sublime.active_window().show_input_panel(
            "Set server address.  Press Enter:",
            info["Server"],
            self.set,
            None,
            self.on_cancel,
        )

    def set(self, addr):
        addr = addr.strip()
        if len(addr) > 0:
            try:
                with open(gemsFILE, "r") as f:
                    info = json.loads(f.read())
            except:
                info = dict()
            if not addr.startswith("http://"):
                addr = "http://" + addr
            info["Server"] = addr
            with open(gemsFILE, "w") as f:
                f.write(json.dumps(info, indent=4))
        else:
            sublime.message_dialog("Server address cannot be empty.")


# ------------------------------------------------------------------


class gemsSetCourseId(sublime_plugin.ApplicationCommand):
    def run(self):
        try:
            with open(gemsFILE, "r") as f:
                info = json.loads(f.read())
        except:
            info = dict()

        if "CourseId" not in info:
            info["CourseId"] = ""
        if sublime.active_window().id() == 0:
            sublime.run_command("new_window")
        sublime.active_window().show_input_panel(
            "Set course id.  Press Enter:", info["CourseId"], self.set, None, None
        )

    def set(self, cid):
        cid = cid.strip()
        if len(cid) > 0:
            try:
                with open(gemsFILE, "r") as f:
                    info = json.loads(f.read())
            except:
                info = dict()
            info["CourseId"] = cid
            with open(gemsFILE, "w") as f:
                f.write(json.dumps(info, indent=4))
            sublime.message_dialog("Course id is set to " + cid)
        else:
            sublime.message_dialog("Server address cannot be empty.")


# ------------------------------------------------------------------


class gemsSetName(sublime_plugin.ApplicationCommand):
    def run(self):
        try:
            with open(gemsFILE, "r") as f:
                info = json.loads(f.read())
        except:
            info = dict()
        if "Name" not in info:
            info["Name"] = ""
        if sublime.active_window().id() == 0:
            sublime.run_command("new_window")
        sublime.active_window().show_input_panel(
            "Set assgined username.  Press Enter:", info["Name"], self.set, None, None
        )

    def set(self, name):
        name = name.strip()
        if len(name) > 0:
            try:
                with open(gemsFILE, "r") as f:
                    info = json.loads(f.read())
            except:
                info = dict()
            info["Name"] = name
            with open(gemsFILE, "w") as f:
                f.write(json.dumps(info, indent=4))
            sublime.message_dialog("Assigned name is set to " + name)
        else:
            sublime.message_dialog("Name cannot be empty.")


# ------------------------------------------------------------------


class gemsAddTestCase(sublime_plugin.ApplicationCommand):
    def run(self):
        sublime.active_window().show_input_panel(
            "Input: ", "", self.set_input, None, None
        )

    def set_input(self, test_input):
        test_input = test_input.strip()
        if len(test_input) > 0:
            self.test_input = test_input
            sublime.active_window().show_input_panel(
                "Expected Output: ", "", self.set_output, None, None
            )
        else:
            sublime.message_dialog("Input cannot be empty.")

    def set_output(self, expected_output):
        expected_output = expected_output.strip()
        if len(expected_output) > 0:
            test_cases = []
            testCaseFile = (
                sublime.active_window().active_view().file_name() + ".test_cases"
            )
            testCaseFile = os.path.join(
                os.path.dirname(os.path.realpath(__file__)), testCaseFile
            )
            try:
                with open(testCaseFile, "r") as f:
                    test_cases = json.loads(f.read())
            except:
                test_cases = []
            test_cases.append({"input": self.test_input, "output": expected_output})

            # with open(testCaseFile, 'w') as f:
            # 	f.write(json.dumps(test_cases, indent=4))
            json.dump(test_cases, open(testCaseFile, "w"))
            sublime.message_dialog("Test case saved successfully.")
        else:
            sublime.message_dialog("Expected output cannot be empty.")


class gemsGetTestCase(sublime_plugin.TextCommand):
    def run(self, edit):
        fname = self.view.file_name()
        fname = os.path.basename(fname)
        if fname is None:
            sublime.message_dialog("Cannot share unsaved content.")
            return

        response = gemsRequest("get_testcase", {"file_name": fname})
        if response is None:
            return

        json_obj = json.loads(response)
        if json_obj == "":
            sublime.message_dialog(
                "No test case found for this problem. Try again later."
            )
            return

        content = ""
        for i, tc in enumerate(json_obj, 1):
            input = tc["input"]
            output = tc["output"]
            content += (
                "Input "
                + str(i)
                + "\n"
                + input
                + "\nExpected Output for Input "
                + str(i)
                + "\n"
                + output
                + "\n\n\n"
            )

        local_file = os.path.join(gemsFOLDER, fname + ".test_case.txt")
        with open(local_file, "w", encoding="utf-8") as fp:
            fp.write(content)
        if sublime.active_window().id() == 0:
            sublime.run_command("new_window")
        sublime.active_window().open_file(local_file)


# ------------------------------------------------------------------


class gemsUpdate(sublime_plugin.WindowCommand):
    def run(self):
        package_path = os.path.join(sublime.packages_path(), "GEMStudent")
        try:
            version = open(os.path.join(package_path, "VERSION")).read()
        except:
            version = 0
        if sublime.ok_cancel_dialog(
            "Current version is {}. Click OK to update.".format(version)
        ):
            if not os.path.isdir(package_path):
                os.mkdir(package_path)
            module_file = os.path.join(package_path, "GEMStudent.py")
            menu_file = os.path.join(package_path, "Main.sublime-menu")
            version_file = os.path.join(package_path, "version.go")
            urllib.request.urlretrieve(
                "https://raw.githubusercontent.com/msh-shiplu/CodeSpace/2.0/src/GEMStudent/GEMStudent.py",
                module_file,
            )
            urllib.request.urlretrieve(
                "https://raw.githubusercontent.com/msh-shiplu/CodeSpace/2.0/src/GEMStudent/Main.sublime-menu",
                menu_file,
            )
            urllib.request.urlretrieve(
                "https://raw.githubusercontent.com/msh-shiplu/CodeSpace/2.0/src/version.go",
                version_file,
            )
            with open(version_file) as f:
                lines = f.readlines()
            for line in lines:
                if line.strip().startswith("const VERSION ="):
                    prefix, version = line.strip().split("const VERSION =")
                    version = version.strip().strip('"')
                    break
            os.remove(version_file)
            with open(os.path.join(package_path, "VERSION"), "w") as f:
                f.write(version)
            sublime.message_dialog(
                "GEM has been updated to version {}.".format(version)
            )


# ------------------------------------------------------------------


class gemsGetFriendCode(sublime_plugin.TextCommand):
    def is_enabled(self):
        global gemsCurrentHelpSubId
        if gemsCurrentHelpSubId is not None:
            return False
        return True

    def run(self, edit):
        global gemsCurrentHelpSubId
        global gemsHelpRequestMessage
        if gemsCurrentHelpSubId is not None:
            sublime.message_dialog("You already have a submission to help")
            return

        filename = self.view.file_name()
        # print(filename)
        filename = os.path.basename(filename)

        data = {"filename": filename}
        response = gemsRequest("student_get_help_code", data)
        if response is None:
            sublime.message_dialog("Could not load any help submission")
            return
        response = json.loads(response)
        content = response["Content"]
        filename = response["Filename"]
        status = response["Status"]
        if status > 0:
            sublime.message_dialog(gemsHelpRequestMessage[status])
            return

        gemsCurrentHelpSubId = response["Sid"]
        # print("Submission ID", gemsCurrentHelpSubId)
        helpFolder = os.path.join(gemsFOLDER, "HelpSubmissions/")
        if not os.path.exists(helpFolder):
            os.mkdir(helpFolder)

        local_file = os.path.join(helpFolder, filename)
        # print(helpFolder, local_file, filename)
        with open(local_file, "w", encoding="utf-8") as f:
            f.write(content)
        if sublime.active_window().id() == 0:
            sublime.run_command("new_window")
        sublime.active_window().open_file(local_file)
        # sublime.message_dialog("")
        sublime.active_window().active_view().set_read_only(True)
        sublime.message_dialog(
            "Press Enter to send feedback. Press Esc to return without feedback."
        )
        sublime.active_window().show_input_panel(
            "Feedback:", "", self.send_help_message, None, self.return_without_feedback
        )

    def send_help_message(self, message):
        global gemsCurrentHelpSubId

        if message is None or message == "":
            self.return_without_feedback()
            # sublime.message_dialog("Help message can not be empty!")
            # return
        data = {"submission_id": gemsCurrentHelpSubId, "message": message}
        response = gemsRequest("student_send_help_message", data)
        gemsCurrentHelpSubId = None
        sublime.active_window().run_command("close")
        sublime.message_dialog(response)

    def return_without_feedback(self):
        global gemsCurrentHelpSubId
        data = {"submission_id": gemsCurrentHelpSubId}
        response = gemsRequest("student_return_without_feedback", data)
        gemsCurrentHelpSubId = None
        sublime.active_window().run_command("close")
        sublime.message_dialog(response)


# class gemsReturnWithoutFeedback(sublime_plugin.WindowCommand):

# 	def is_enabled(self):
# 		global gemsCurrentHelpSubId
# 		if gemsCurrentHelpSubId is None:
# 			return False
# 		return True

# 	def run(self):
# 		global gemsCurrentHelpSubId
# 		# print("Sub id in return", gemsCurrentHelpSubId)
# 		if gemsCurrentHelpSubId is None:
# 			sublime.message_dialog("You don't have any submission to return")
# 			return

# 		data = {"submission_id": gemsCurrentHelpSubId}
# 		response = gemsRequest("student_return_without_feedback", data)
# 		gemsCurrentHelpSubId = None
# 		sublime.message_dialog(response)
# 		self.window.run_command("close")
# 		self.window.run_command("hide_panel", {"cancel": True})

# class gemsEventListeners(sublime_plugin.EventListener):

#     def on_pre_close(self, view):
#         if isRegistered == False:
#             return
#         filename = os.path.basename(view.file_name())
#         fn_splits = filename.split("-")
#         if filename is not None and len(fn_splits) > 2 and fn_splits[0] == "feedback" and fn_splits[1].isdigit():
#             view.window().focus_view(view)
#             resp = sublime.yes_no_cancel_dialog(
#                 "Thank you, this helps!!!", "Yes", "No")
#             data = {"feedback_id": int(fn_splits[1])}
#             if resp == sublime.DIALOG_YES:
#                 data["feedback"] = "yes"
#             else:
#                 data["feedback"] = "no"
#             data["role"] = "student"
#             gemsRequest("save_snapshot_back_feedback", data)


def check_message_feedback(feedback_id):
    response = gemsRequest("has_message_feedback", {"feedback_id": feedback_id})
    if response == "yes":
        return True
    return False


class gemsEventListeners(sublime_plugin.EventListener):
    def on_pre_close(self, view):
        if isRegistered == False or view is None or view.file_name() is None:
            return
        filename = os.path.basename(view.file_name())
        fn_splits = filename.split("-")
        if (
            filename is not None
            and len(fn_splits) > 2
            and fn_splits[0] == "feedback"
            and fn_splits[1].isdigit()
        ):
            view.window().focus_view(view)
            feedback_dir = os.path.join(gemsFOLDER, "FEEDBACK")
            local_file = os.path.join(feedback_dir, filename)
            fname = "".join(fn_splits[2:])
            ask_for_back_feedback(fname, local_file, True)

    # def on_deactivated_async(self, view):
    #     if view.file_name() != None:
    #         filename = os.path.basename(view.file_name())
    #         fn_splits = filename.split("-")
    #         if filename is not None and len(fn_splits) > 2 and fn_splits[0] == "feedback" and fn_splits[1].isdigit():
    #             feedback_dir = os.path.join(gemsFOLDER, 'FEEDBACK')
    #             local_file = os.path.join(feedback_dir, filename)
    #             fname = "".join(fn_splits[2:])
    #             if local_file not in feedback_resp:
    #                 feedback_resp.append(local_file)
    # ask_for_back_feedback(view, fname, local_file, True)
    # sublime.active_window().open_file(view.file_name())
    # print("Filename: ",view.file_name())
    # print("Active window: ",sublime.active_window().id() )
    # if sublime.active_window().id() == 0:
    #     sublime.run_command('new_window')
    # view = sublime.active_window().open_file(view.file_name())
    # view.run_command('enter_insert_mode')


def ask_for_back_feedback(filename, feedback_filename, fromEvent=False):
    feedback_id = os.path.basename(feedback_filename).split("-")[1]

    if check_message_feedback(feedback_id):
        return
        # sublime.active_window().open_file(feedback_filename)

    # show Yes Option as first.
    resp = sublime.yes_no_cancel_dialog(
        "Was this feedback helpful? Please answer Yes or No", "No", "Yes"
    )
    # consider the first option as YES
    if resp == sublime.DIALOG_NO:
        send_student_back_feedback(filename, "yes", feedback_filename)
    elif resp == sublime.DIALOG_YES:
        send_student_back_feedback(filename, "no", feedback_filename)
    elif resp == sublime.DIALOG_CANCEL:
        ask_for_back_feedback(filename, feedback_filename, True)
    if fromEvent == False:
        sublime.active_window().active_view().close()


def send_student_back_feedback(filename, response, feedback_filename):
    global gemsBackFeedbackTimers
    # global gotFeedback
    global gemsBackFeedbackStatus

    if (filename, feedback_filename) in gemsBackFeedbackTimers:
        gemsBackFeedbackTimers[(filename, feedback_filename)].cancel()
        gemsBackFeedbackStatus[(filename, feedback_filename)] = False
    # gotFeedback[filename].remove(feedback_filename)
    feedback_filename = os.path.basename(feedback_filename)
    feedback_id = feedback_filename.split("-")[1]
    data = dict(
        filename=filename,
        feedback=response,
        feedback_id=feedback_id,
        role="student",
    )
    gemsRequest("save_snapshot_back_feedback", data)


# ------------------------------------------------------------------


class gemsViewHelpRequests(sublime_plugin.ApplicationCommand):
    def run(self):
        with open(gemsFILE, "r") as f:
            info = json.loads(f.read())
            p = urllib.parse.urlencode(
                {"password": info["Password"], "uid": info["Uid"], "role": "student"}
            )
        webbrowser.open(gemsSERVER + "/help_requests?" + p)


# ------------------------------------------------------------------


class gemsViewExercises(sublime_plugin.ApplicationCommand):
    def run(self):
        global gemsSERVER
        with open(gemsFILE, "r") as f:
            info = json.loads(f.read())
        p = urllib.parse.urlencode(
            {"password": info["Password"], "uid": info["Uid"], "role": "student"}
        )
        webbrowser.open(gemsSERVER + "/view_exercises?" + p)


class gemsPeerTutoring(sublime_plugin.ApplicationCommand):
    def run(self):
        global gemsSERVER
        with open(gemsFILE, "r") as f:
            info = json.loads(f.read())
        p = urllib.parse.urlencode(
            {"password": info["Password"], "uid": info["Uid"], "role": "student"}
        )
        fname = sublime.active_window().active_view().file_name()
        if fname is None:
            sublime.message_dialog("Unable to find filename!")
            return
        filename = os.path.basename(fname)
        match = re.search(r"feedback-\d+-(.*)", filename)
        if match:
            filename = match.group(1)
        data = dict(
            filename=filename,
            role="student",
        )
        response = gemsRequest(gemsSERVER + "/peer_tutoring", data)
        if response == "redirect":
            webbrowser.open(gemsSERVER + "/view_exercises?" + p)
        else:
            sublime.message_dialog(response)
