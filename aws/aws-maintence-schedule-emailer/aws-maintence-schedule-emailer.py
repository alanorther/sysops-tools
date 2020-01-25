#!/usr/bin/python

from collections import defaultdict

import boto3

"""
A tool for retrieving information about EC2 instances that are pending maintenance.
"""

ec2client = boto3.client('ec2')

filters=[{'Name': 'event.code', 'Value': ['system-reboot', 'instance-reboot', 'system-maintenance', 'instance-retirement', 'instance-stop']}]

response = ec2client.describe_instance_status()

for instance in response["InstanceStatuses"]:

    #print instance

    if 'Events' in instance:
        print "%s %s" % (instance['InstanceId'], instance['Events'])

quit()