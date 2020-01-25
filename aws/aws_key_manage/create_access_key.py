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

user_name = ""

def print_help():
    print "usage: " + argv[0] + " [options....] "
    print "options:"
    print "-u, --user       Username of the key you'd like to create"
    quit()

try: 
    opts, args = getopt.getopt(argv[1:], "u:", ["user="])
except getopt.GetoptErrror:
    print_help()

for opt, arg in opts: 
    if opt in("-u", "--user"):
        user_name = arg

if not user_name: 
    print_help()

iam = boto3.client('iam')

response = iam.create_access_key(
        UserName=user_name
        )

print(response['AccessKey'])
