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

aws_region = ['us-east-1', 'us-east-2', 'us-west-1', 'us-west-2', 'ap-south-1', 'ap-northeast-2', 'ap-southeast-1', 'ap-southeast-2', 'ap-northeast-1', 'ca-central-1', 'eu-central-1', 'eu-west-1', 'eu-west-2', 'sa-east-1']

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
        for r in aws_region:
            security_groups = subprocess.call("aws ec2 describe-security-groups --region " + r + " --profile "
                                                      + x.translate(None, '[]') + " | grep " + user_key.strip(), shell=True)
            if security_groups == 0:
                sg_out = subprocess.check_output("aws ec2 describe-security-groups --region " + r + " --profile "
                                                 + x.translate(None, '[]') + " --group-ids " + user_key.strip(), shell=True)
                print sg_out
                print "Security Group in Region: " + r
                quit()

def print_help():
    print "usage: " + argv[0] + " [options....] "
    print "options:"
    print "-sg, --security-group       The Access key you are looking for"
    quit()

try:
    opts, args = getopt.getopt(argv[1:], "s:", ["security-group="])
except getopt.GetoptError: 
    print_help()

for opt, arg in opts: 
    if opt in ("-s", "--security-group"):
        user_key = arg

if not user_key: 
    print_help()

print "Searching seucurity group " + user_key.strip()

get_aws_env()
get_access_key()
