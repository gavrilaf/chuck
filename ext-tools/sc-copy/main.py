"""Chuck scenarios cleaner

Usage:
  sc_copy.py copy <from> <to>
  sc_copy.py new <from> <to>
  sc_copy.py (-h | --help)

Options:
  -h --help     Show this screen.
"""

import os.path
import csv
from utils import *
from docopt import docopt
import colorama
from colorama import Fore

STRIP_CHARS = "\t, "

MODE_AUTO = 1
MODE_ASK = 2
MODE_SKIP_ALL = 3
MODE_COPY_ALL = 4


skip_mode = MODE_AUTO


def ask_for_skip_line(id, code, method, url):
    global skip_mode

    while True:
        prompt = Fore.GREEN + "Copy {} : {}, code = {}, id = {}, [(Y)es/(N)o/(S)kip all/(C)opy all]: ".format(method, url, code, id)
        cp = input(prompt).lower()
        if cp == "y":
            return True
        elif cp == "n":
            return False
        elif cp == "s":
            skip_mode = MODE_SKIP_ALL
            return False
        elif cp == "c":
            skip_mode = MODE_COPY_ALL
            return True


def is_skip_line(src_path, id, code, method, url):
    if skip_mode == MODE_AUTO:
        return check_skip_by_code(src_path, id, code)
    elif skip_mode == MODE_ASK:
        return not ask_for_skip_line(id, code, method, url)
    elif skip_mode == MODE_SKIP_ALL:
        return True
    elif skip_mode == MODE_COPY_ALL:
        return False
    else:
        print("Unknown code")
        exit(1)


def copy_scenario(mode, src_path, dest_path, name):
    print(Fore.GREEN + "\nCopy scenario {}".format(name))

    global skip_mode
    skip_mode = mode

    with open(os.path.join(src_path, 'index.txt'), 'r') as f:
        reader = csv.reader(f, delimiter = '\t')
        index = list(reader)

    result = []
    ids_map = {}
    line_indx = 0
    id_indx = 1
    processed = set()

    for line in index:
        code = line[1].strip(STRIP_CHARS)
        id = line[2].strip(STRIP_CHARS)
        method = line[3].strip(STRIP_CHARS)
        url = line[4].strip(STRIP_CHARS)

        key = method + url
        if key in processed:
            print(Fore.YELLOW + "Duplicated line {}, {}, {} : {}".format(line_indx, code, method, url))
            line_indx += 1
            continue
        else:
            processed.add(key)

        if is_skip_line(src_path, id, code, method, url):
            print(Fore.YELLOW + "Skip line {}, {}, {} : {}".format(line_indx, code, method, url))
            line_indx += 1
            continue

        new_id = "r_{}".format(id_indx)
        ids_map[id] = new_id
        id_indx += 1

        new_url = clear_url(url)

        line_indx += 1
        result.append((code, new_id, method, new_url))

    if len(result) == 0:
        print(Fore.RED + "Empty scenario {}\n".format(name))
        return

    dest_path = os.path.join(dest_path, name)
    os.makedirs(dest_path, exist_ok = True)

    with open(os.path.join(dest_path, 'index.txt'), 'w') as f:
        for r in result:
            line = "F,\t{},\t{},\t{},\t{}\n".format(r[0], r[1], r[2], r[3])
            f.write(line)

    for old_id, new_id in ids_map.items():
        copy_stub(os.path.join(src_path, old_id), os.path.join(dest_path, new_id))


def copy_scenarios(src_path, dest_path):
    print(Fore.GREEN + "Copy scenarios from {} to {}\n".format(src_path, dest_path))
    for dirName, subdirList, fileList in os.walk(src_path):
        if "index.txt" in fileList:
            subdirList.clear()
            sc_name = os.path.basename(dirName)
            copy_scenario(MODE_AUTO, dirName, dest_path, sc_name)


def create_scenario(src_path, dest_path):
    dirs = os.listdir(src_path)
    dirs.remove(".DS_Store")
    if len(dirs) == 1:
        log_name = dirs[0]
    else:
        print(Fore.GREEN + "Logs {}\n", ", ".join(dirs))
        log_name = input(Fore.GREEN + "Enter log folder: ")

    scenario_name = input(Fore.GREEN + "Enter new scenario name: ")
    src_path = os.path.join(src_path, log_name)

    print(Fore.GREEN + "Copy new scenario {} based on {} to {}\n".format(scenario_name, src_path, dest_path))
    copy_scenario(MODE_ASK, src_path, dest_path, scenario_name)


def main(args):
    if args["copy"]:
        src = args["<from>"]
        dest = args["<to>"]
        copy_scenarios(src, dest)
    elif args["new"]:
        src = args["<from>"]
        dest = args["<to>"]
        create_scenario(src, dest)
    else:
        print(Fore.RED + "Unknown mode")


if __name__ == "__main__":
    colorama.init()
    args = docopt(__doc__, version='v0.1')
    main(args)
    colorama.deinit()