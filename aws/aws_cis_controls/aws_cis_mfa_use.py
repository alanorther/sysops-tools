#!/usr/bin/python
import boto3
import csv
import os
import getopt
from sys import argv

aws_env = "default"
aws_region = "us-east-1"
aws_secret_id = "" 
aws_secret_key = ""

def print_help():
    print "usage: " + argv[0] + " [options....] "
    print "options:"
    print "-e, --env       Environment you want to use - optional"
    print "-r, --region    Region you want to use - optional"
    print "-h, --help      Print this message"
    quit()


def get_creds(): 
    global aws_env
    global aws_secret_id 
    global aws_secret_key
    counter = 0
    cred_file = os.path.expanduser("~")
    cred_file += "/.aws/credentials"
    r_cred = open(cred_file, 'r')
    for row in r_cred:
        if counter == 2: 
            split_key = row.split("=")
            aws_secret_key = split_key[1].strip()
            counter = 0
        if counter == 1:
            split_id = row.split("=")
            aws_secret_id = split_id[1].strip()
            counter = 2
        if aws_env in row:
            counter = 1
        if counter == 0:
            pass

try: 
    opts, args = getopt.getopt(argv[1:], "he:r:", ["help", "env=", "region="])
except getopt.GetoptErrror:
    print_help()

for opt, arg in opts: 
    if opt in("-e", "--env"):
        aws_env = arg
    if opt in("-r", "--region"):
        aws_region = arg
    if opt in("-h", "--help"):
        print_help()

def get_cred_report():
        status = ""
	client = boto3.client(
	        'iam',
	        aws_access_key_id=aws_secret_id,
	        aws_secret_access_key=aws_secret_key,
	)
        while client.generate_credential_report()['State'] != "COMPLETE":
            time.sleep(2)
            x += 1

            if x > 10:
                status = "Fail: no CredentialReport available."
                break

        if "Fail" in status:
            return status

        response = client.get_credential_report()
        report = []
        reader = csv.DictReader(response['Content'].splitlines(), delimiter=',')
        for row in reader: 
            report.append(row)

         # Verify if root key's never been used, if so add N/A
	try:
	    if report[0]['access_key_1_last_used_date']:
	        pass
	except:
	    report[0]['access_key_1_last_used_date'] = "N/A"
	try:
	    if report[0]['access_key_2_last_used_date']:
	        pass
	except:
	    report[0]['access_key_2_last_used_date'] = "N/A"
	return report
	
	try:
	    credreport = client.get_credential_report()
	except:
	    print "Unable to retrieve credential report" 

get_creds()
cred_report = get_cred_report()
offenders = []
good_job = []

for i in range(len(cred_report)):
    if cred_report[i]['password_enabled'] == "true":
        if cred_report[i]['mfa_active'] == "true":
            good_job.append(str(cred_report[i]['arn']))
        if cred_report[i]['mfa_active'] == "false":
            offenders.append(str(cred_report[i]['arn']))


print "Good Job!"
for i in range(len(good_job)):
    print good_job[i]

print "No good!"
for i in range(len(offenders)):
    print offenders[i]


