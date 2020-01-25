#!/usr/bin/python
import os
import subprocess 
import time 
import progressbar
import sys
import random
import getopt
from sys import argv
user_key = ""
aws_env = list()
def get_aws_env():
    global aws_env
    cred_file = os.path.expanduser("~")
    cred_file += "/.aws/credentials"
    r_cred = open(cred_file, 'r')
    for row in r_cred:
        if "aws_access_key_id" in row or "aws_secret_access_key" in row or row == "":
            continue
        elif "[" in row and "]" in row:
            aws_env.append(row.strip())
        else:
            continue
def get_access_key():
    for x in aws_env:
        users = subprocess.check_output("aws iam list-users --output text --profile " + x.translate(None, '[]') + " | awk '{print $NF}'", shell=True)
        user_split = users.splitlines()
        for u in user_split:
            subprocess.call("aws iam list-access-keys --user " + u.strip() + " --output text --profile " + x.translate(None, '[]') + " | grep " + user_key.strip(), shell=True)
def print_help():
    print "usage: " + argv[0] + " [options....] "
    print "options:"
    print "-k, --key       The Access key you are looking for"
    quit()
try:
    opts, args = getopt.getopt(argv[1:], "k:", ["key="])
except getopt.GetoptError: 
    print_help()
for opt, arg in opts: 
    if opt in ("-k", "--key"): 
        user_key = arg
if not user_key: 
    print_help()
get_aws_env()
get_access_key()
