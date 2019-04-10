import os
import os.path
import shutil
import re

re_apikey = re.compile(r"apikey=([^&#]*)", re.IGNORECASE)
re_code = re.compile(r"code=([^&#]*)", re.IGNORECASE)
re_verifier = re.compile(r"verifier=([^&#]*)", re.IGNORECASE)

NOT_FOUND_SKIP_RULES = [
    "\"message\": \"Not Found\",",
    "\"status\": \"404\""
]

def copytree(src, dst):
    if not os.path.exists(dst):
        os.makedirs(dst)
    for item in os.listdir(src):
        s = os.path.join(src, item)
        d = os.path.join(dst, item)
        if os.path.isdir(s):
            copytree(s, d)
        else:
            if not os.path.exists(d) or os.stat(s).st_mtime - os.stat(d).st_mtime > 1:
                shutil.copy2(s, d)


def read_content(src_path, id):
    body_path = os.path.join(src_path, id, "resp_body.json")
    if not os.path.isfile(body_path):
        return []
    with open(body_path, 'r') as f:
        content = f.readlines()
        return content


def check_skip_by_code(src_path, id, code):
    if code == "404":
        content = read_content(src_path, id)
        for line in content:
            line = line.strip()
            if line in NOT_FOUND_SKIP_RULES:
                return True
    return False


def clear_url(url):
    url = re.sub(re_apikey, "apikey=*", url)
    url = re.sub(re_code, "code=*", url)
    url = re.sub(re_verifier, "verifier=*", url)
    return url
