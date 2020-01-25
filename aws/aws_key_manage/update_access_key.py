#!/usr/bin/python
import boto3
import os
import subprocess
import time
import progressbar
import sys
import random
import getopt
from sys import argv

access_key_id = ""
status = ""
user = ""

def print_help():
    print "usage: " + argv[0] + " [options....] "
    print "options:"
    print "-u, --user               Username of the key you'd like to update"
    print "-a, --access_key_id      Access Key ID you'd like to update"
    print "-s, --status             Active, Inactive, Delete"
    quit()

try: 
    opts, args = getopt.getopt(argv[1:], "a:u:s:", ["access_key_id=", "user_name=", "status="])
except getopt.GetoptErrror:
    print_help()

for opt, arg in opts: 
    if opt in("-a", "--access_key_id"):
        access_key_id = arg
    if opt in("-u", "--user"):
        user = arg
    if opt in("-s", "--status"):
        status = arg

if not access_key_id: 
    print_help()
if not user: 
    print_help()
if not status: 
    print_help()

iam = boto3.client('iam')

response = iam.update_access_key(
        AccessKeyId = access_key_id,
        Status = status,
        UserName = user
        )

