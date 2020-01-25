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

def print_help():
    print "usage: " + argv[0] + " [options....] "
    print "options:"
    print "-u, --user       Username of the key you'd like to create"
    quit()

try: 
    opts, args = getopt.getopt(argv[1:], "a:", ["access_key_id="])
except getopt.GetoptErrror:
    print_help()

for opt, arg in opts: 
    if opt in("-a", "--access_key_id"):
        access_key_id = arg

if not access_key_id: 
    print_help()

iam = boto3.client('iam')

response = iam.get_access_key_last_used(
        AccessKeyId = access_key_id
        #AccessKeyId = 'AKIAJB3AUX57TMSO2HUA'
        )

print(response['AccessKeyLastUsed'])

