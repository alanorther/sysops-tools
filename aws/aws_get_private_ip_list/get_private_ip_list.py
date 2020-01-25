#!/usr/bin/python

import os
import subprocess 
import time 
import sys
import random
import getopt

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


def get_ip_addr():
    for x in aws_env:
        for r in aws_region:
            ip = subprocess.check_output("aws ec2 describe-addresses --region " +
                                         r + " --query Addresses[*].PrivateIpAddress --profile " + x.translate(None, '[]'), shell=True)
            if any(str.isdigit(c) for c in ip):
                print "Environment: " + x + " Region: " + r
                print ip.translate(None, '"[],\b\t\w')


get_aws_env()
get_ip_addr()

