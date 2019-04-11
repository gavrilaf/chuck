"""Chuck scenarios cleaner

Usage:
  sc_copy.py copy <from> <to> (auto|ask)
  sc_copy.py new <name> <from> <to>
  sc_copy.py (-h | --help)

Options:
  -h --help     Show this screen.
"""

import sys
import os.path
import csv
from utils import *
from docopt import docopt

MODE_AUTO = 1
MODE_INTERACTIVE = 2

STRIP_CHARS = "\t, "


def ask_for_skip_line(id, code, method, url):
    return False


def skip_line(mode, src_path, id, code, method, url):
    return check_skip_by_code(src_path, id, code)


def copy_scenario(mode, src_path, dest_path, name):
    print("Copy scenario {} from {} to {}".format(name, src_path, dest_path))

    with open(os.path.join(src_path, 'index.txt'), 'r') as f:
        reader = csv.reader(f, delimiter = '\t')
        index = list(reader)

    result = []
    ids_map = {}
    line_indx = 0
    id_indx = 1
    for line in index:
        code = line[1].strip(STRIP_CHARS)
        id = line[2].strip(STRIP_CHARS)
        method = line[3].strip(STRIP_CHARS)
        url = line[4].strip(STRIP_CHARS)

        if skip_line(mode, src_path, id, code, method, url):
            print("Skip line {}, status code check, code: {}, method: {}, url: {}".format(line_indx, code, method, url))
            line_indx += 1
            continue

        new_id = "r_{}".format(id_indx)
        ids_map[id] = new_id
        id_indx += 1

        new_url = clear_url(url)

        line_indx += 1
        result.append((code, new_id, method, new_url))

    if len(result) == 0:
        print("Empty scenario {}".format(name))
        return

    dest_path = os.path.join(dest_path, name)
    os.makedirs(dest_path, exist_ok = True)

    with open(os.path.join(dest_path, 'index.txt'), 'w') as f:
        for r in result:
            line = "F,\t{},\t{},\t{},\t{}\n".format(r[0], r[1], r[2], r[3])
            f.write(line)

    for old_id, new_id in ids_map.items():
        copytree(os.path.join(src_path, old_id), os.path.join(dest_path, new_id))


def copy_scenarios(mode, src_path, dest_path):
    print("Copy scenarios from {} to {}".format(src_path, dest_path))
    for dirName, subdirList, fileList in os.walk(src_path):
        if "index.txt" in fileList:
            subdirList.clear()
            sc_name = os.path.basename(dirName)
            copy_scenario(mode, dirName, dest_path, sc_name)


def main(args):
    print("sc-copy main")

    if args["copy"]:
        mode = MODE_INTERACTIVE if args["ask"] else MODE_AUTO
        src = args["<from>"]
        dest = args["<to>"]
        copy_scenarios(mode, src, dest)
    else:
        print("Doesn't supported yet")

    # src = "./../../log-intg/2019_4_10_9_40_52"
    # dest = "../../cleaned"
    # copy_scenarios(mode, src, dest)


if __name__ == "__main__":
    args = docopt(__doc__, version='v0.1')
    main(args)