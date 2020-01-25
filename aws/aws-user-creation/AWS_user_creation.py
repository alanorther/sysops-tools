#!/usr/bin/python
# Created by Alan M. Orther - alan.orther@verve.com
# Start date June 20th 2017.
# This script is used to create an AWS user, S3 folder, and output the
# credentials to use for our Verve partners or internal users.

import os
import subprocess
import re
from sys import argv
from time import sleep

# Global Variables
username = ""
username_correct = ""
company_name = ""
folder_name = ""
aws_profile = ""
temp_file = "/tmp/AWS_tools_creation.tmp"
temp_aws_policy_file = ""
aws_username_list = []
data = []
iam_policy_document = ""
verve_partner = True
s3_bucket = "verve-opsdata"
s3_folder_key = "data/partners/"
US_partner = True


# Define general error and quit.
def error_message():
    print "*** ERROR ~ ERROR ~ ERROR ***\n"
    print "Please use one of the formats below:\n"
    print "AWS_tools.py"
    print "OR"
    print "AWS_tools.py <AWS Profile Name>\n"
    print "(Generally the AWS profile names are in ~/.aws/credentials)\n"
    print "REQUIRES AWS CLI INSTALLED\n"
    delete_all_temp_files()
    quit()

# Use argument as the AWS profile used.
if len(argv) > 2:
    print "\n   You have added too many arguments.\n"
    error_message()

# Unwrap the commandline arguments.
if len(argv) == 2:
    script_name, aws_profile = argv

# Create a help menu.
if aws_profile in ['-h', '--h', '-help', '--help']:
    print "\nPlease use one of the formats below:\n"
    print "AWS_tools.py"
    print "OR"
    print "AWS_tools.py <AWS Profile Name>\n"
    print "(Generally the AWS profile names are in ~/.aws/credentials)"
    print "REQUIRES AWS CLI INSTALLED\n"
    quit()


# Prompt for a username, clean and check syntax
def get_username():
    global username, username_correct, company_name, verve_partner, \
        s3_bucket, s3_folder_key, folder_name, US_partner
    print "\nCREATE AN AWS USER AND FOLDER IN S3."
    print "This script requires the AWS CLI tool to be installed locally"
    internal_or_partner = raw_input("\nIs this for a Verve user or external "
                                    "partner user? (v/P)").lower() or "partner"

    # Change global variables for internal verve user.
    if internal_or_partner in ['v', 'verve', 'verve user']:
        verve_partner = False
        s3_bucket = "verve-home"
        s3_folder_key = ""

    # Request the username or company based on verve_partner value.
    if verve_partner is False:
        username = raw_input("\nPlease input the username: "
                             " ").lower()
    else:
        username = raw_input("\nPlease input the company name: "
                             " ").lower()

    # Remove unwanted characters and replace spaces with '_'.
    username = re.sub(r'[`|~|!|@|#|$|%|^|&|*|(|)|+|?|>|<|,|\\|{|}|"|\'|:|;'
                      r'|=|/|\[|\]]', r'', username)
    username = username.strip()
    username = username.replace(' ', '_')
    company_name = username
    folder_name = re.sub(r'[.]', r'', username)

    # If this is a external Verve partner, add the 'ext-partner' prefix.
    if verve_partner is True:
        username = "ext-partner-" + username
    username_correct = raw_input("\nThe username that will be used is: " +
                                 username + "\n" + "Is that okay? (Y/n/quit)"
                                                   " ").lower() or "yes"

    # Prompt if the displayed username is correct. If not, repeat.
    if username_correct in ['y', 'yes', 'yeah', 'yep', 'okay', 'f-yeah!']:
        print "Okay. Using " + username
    elif username_correct in ['q', 'quit']:
        quit()
    else:
        get_username()

    # Prompt if the partner is in US or UK.
    # Not needed anymore.
    #if verve_partner is True:
    #    US_or_UK = raw_input("\nIs this company in the US or UK? [US "
    #                         "Default] ").lower()
    #
    #    if US_or_UK in ['uk']:
    #        s3_bucket = "verve-opsdata-uk"
    #        s3_folder_key = "data/partner/"

# Get list of AWS usernames based on if profile arg was used.
def get_aws_usernames():
    global aws_profile
    print "\nRetrieving all AWS usernames and comparing for duplicates....."
    sleep(1)

    # Retrieve list of users based on if AWS profile arg was given.
    if len(argv) == 1:
        subprocess.call('aws iam list-users | grep UserName'
                        '> /tmp/AWS_tools_creation.tmp', shell=True)
    else:
        # Create a string to use in the subprocess.call
        subprocess_call_string = "aws --profile %s iam list-users | grep " \
                                 "UserName > /tmp/AWS_tools_creation.tmp " \
                                 % aws_profile

        # Run the string that was created above.
        subprocess.call(subprocess_call_string, shell=True)


# Create list from the AWS username temp file.
def create_user_list():
    sleep(1)
    with open("/tmp/AWS_tools_creation.tmp", 'r') as f:
        the_file = f.readlines()
        for n in the_file:
            n = n.replace('"UserName":', '')
            n = n.replace('"', '')
            n = n.replace(',', '')
            n = n.strip()
            aws_username_list.append(n)

    # Delete the temp file.
    os.remove("/tmp/AWS_tools_creation.tmp")


# Compare list to the username that was entered. If same quit.
def compare_users_to_entered_name():
    sleep(1)
    if username in aws_username_list:
        print "\n   *** THIS USERNAME IS ALREADY BEING USED ***\n"
        quit()


# Create AWS user.
def create_aws_user():
    global username
    print "\nCreating AWS user %s ....." % username
    sleep(1)
    if len(argv) == 1:
        # Create a string to use in the subprocess.call
        subprocess_call_string = "aws iam create-user --user-name %s" \
                                 % username

        # Run the string that was created above.
        subprocess.call(subprocess_call_string, shell=True)
    else:
        # Create a string to use in the subprocess.call
        subprocess_call_string = "aws --profile %s iam create-user " \
                                 "--user-name %s" % (aws_profile, username)

        # Run the string that was created above.
        subprocess.call(subprocess_call_string, shell=True)


# Create Secret Access Key for AWS user.
def create_aws_access_key():
    global username, data
    sleep(1)
    # Create temp file name and location.
    access_key_temp_file = "/tmp/AWS_tools_creation_%s_access_key.tmp" % \
                           username

    # Create the key and strip out the unneeded JSON.
    print "\nCreating the AWS Access Key for the user %s ....." % username
    sleep(1)

    # If AWS profile arg is not given use this command. If not, use the other.
    if len(argv) == 1:
        # Create a string to use in the subprocess.call
        call = 'aws iam create-access-key --user-name %s | ' \
               'grep "UserName\|SecretAccessKey\|AccessKeyId" ' \
               '| sed s/\\"UserName\\":// | sed s/' \
               '\\"SecretAccessKey\\":// | sed s/\\"AccessKeyId\\"' \
               ':// | sed s/\\"//g | sed s/,//g | sed s/\ //g >' \
               ' %s' % (username, access_key_temp_file)

        # Run the string that was created above.
        subprocess.call(call, shell=True)

    else:
        # Create a string to use in the subprocess.call
        call = 'aws iam create-access-key --profile %s --user-name %s | ' \
               'grep "UserName\|SecretAccessKey\|AccessKeyId" ' \
               '| sed s/\\"UserName\\":// | sed s/' \
               '\\"SecretAccessKey\\":// | sed s/\\"AccessKeyId\\"' \
               ':// | sed s/\\"//g | sed s/,//g | sed s/\ //g >' \
               ' %s' % (aws_profile, username, access_key_temp_file)

        # Run the string that was created above.
        subprocess.call(call, shell=True)

    # Open the file that was created above, put into list, and remove '\n'.
    with open(access_key_temp_file, 'r') as myfile:
        data = myfile.readlines()
        data = [i.replace('\n', '') for i in data]

    # Delete the temp file.
    os.remove(access_key_temp_file)


# If user is an internal/Verve user, use default S3 folder or create custom.
def default_or_custom_folder():
    global s3_bucket, s3_folder_key, folder_name

    if verve_partner is False:
        folder_question_string = "\nWould you like to use the default " \
                                 "folder s3://%s/%s%s/ ? (y/N)" % \
                                 (s3_bucket, s3_folder_key, folder_name)

        folder_question = raw_input(folder_question_string) or "no"
        folder_question = folder_question.lower()

        # Remove unwanted characters.
        if folder_question in ['n', 'no', 'naw', 'nope', 'f-no']:
            s3_bucket = raw_input("\nWhat bucket would you like to use?"
                                  " Example: verve-home ? ").lower()
            s3_bucket = re.sub(r'[`|~|!|@|#|$|%|^|&|*|(|)|+|?|>|<|.|,|\\|'
                                   r'{|}|"|\'|:|;|=|/|\[|\]]', r'', s3_bucket)

            if folder_question in ['n', 'no', 'naw', 'nope', 'f-no']:
                s3_folder_key_string = "\nWhat folder do you want to make " \
                                       "in the %s bucket? Example: " \
                                       "data/some_folder/ ? " % s3_bucket

                s3_folder_key = raw_input(s3_folder_key_string)
                s3_folder_key = s3_folder_key.lower()
                s3_folder_key = re.sub(r'[`|~|!|@|#|$|%|^|&|*|(|)|+|?|>|<|.'
                                       r'|,|\\|{|}|"|\'|:|;|=|/|\[|\]]', r'',
                                       s3_folder_key)
                if s3_folder_key.endswith('/'):
                    pass
                else:
                    s3_folder_key = s3_folder_key + '/'


# Create the S3 folders.
def create_s3_folder():
    print "\nCreating S3 folder..... (Confirm " \
          "that you are using the AWS profile that has access to the " \
          "appropriate bucket)"
    sleep(1)

    if verve_partner is False:
        # If AWS profile arg is not given use this command. If not,
        # use the other.
        if len(argv) == 1:
            # Create a string to use in the subprocess.call
            call = 'aws s3api put-object --bucket %s ' \
                   '--key %s%s/temp/' % (s3_bucket, s3_folder_key, folder_name)

            # Run the string that was created above.
            subprocess.call(call, shell=True)

        else:
            # Create a string to use in the subprocess.call
            call = 'aws --profile %s s3api put-object --bucket %s ' \
                   '--key %s%s/temp/' % (aws_profile, s3_bucket, s3_folder_key,
                                    folder_name)

            # Run the string that was created above.
            subprocess.call(call, shell=True)
    else:
        # If AWS profile arg is not given use this command. If not,
        # use the other.
        if len(argv) == 1:
            # Create a string to use in the subprocess.call
            call1 = 'aws s3api put-object --bucket %s ' \
                   '--key %s%s/in/' % (s3_bucket, s3_folder_key,
                                       folder_name)

            call2 = 'aws s3api put-object --bucket %s ' \
                    '--key %s%s/out/' % (s3_bucket, s3_folder_key,
                                        folder_name)

            # Run the string that was created above.
            subprocess.call(call1, shell=True)
            sleep(1)
            subprocess.call(call2, shell=True)

        else:
            # Create a string to use in the subprocess.call
            call1 = 'aws --profile %s s3api put-object --bucket %s ' \
                   '--key %s%s/in/' % (aws_profile, s3_bucket, s3_folder_key,
                                       folder_name)

            call2 = 'aws --profile %s s3api put-object --bucket %s ' \
                   '--key %s%s/out/' % (aws_profile, s3_bucket, s3_folder_key,
                                       folder_name)

            # Run the string that was created above.
            subprocess.call(call1, shell=True)
            sleep(1)
            subprocess.call(call2, shell=True)


# Create the IAM policy document that will be used to assign the AWS policy.
def create_iam_policy_document():
    global iam_policy_document, temp_aws_policy_file

    print "\nCreating the AWS IAM policy document....."
    sleep(1)

    # Create the value for the temp_aws_policy_file
    temp_aws_policy_file = "/tmp/AWS_tools_creation_aws_policy_%s.tmp" % company_name

    # Add the company name to the AWS policy JSON document.
    if verve_partner is False:
        iam_policy_document = '''{
        "Version": "2012-10-17",
        "Statement": [
            {
                "Sid": "AllowGroupToSeeBucketListInTheConsole",
                "Effect": "Allow",
                "Action": [
                    "s3:ListAllMyBuckets",
                    "s3:GetBucketLocation"
                ],
                "Resource": "arn:aws:s3:::*"
            },
            {
                "Sid": "AllowRootLevelListingOfTheBucket",
                "Effect": "Allow",
                "Action": [
                    "s3:ListBucket"
                ],
                "Resource": "arn:aws:s3:::%s"
            },
            {
                "Sid": "AllowListBucketOfASpecificUserPrefix",
                "Effect": "Allow",
                "Action": [
                    "s3:ListBucket"
                ],
                "Resource": "arn:aws:s3:::%s",
                "Condition": {
                    "StringLike": {
                        "s3:prefix": [
                            "%s%s/*"
                        ]
                    }
                }
            },
            {
                "Sid": "AllowReadWriteDelete",
                "Effect": "Allow",
                "Action": [
                    "s3:GetObject",
                    "s3:PutObject",
                    "s3:DeleteObject"
                ],
                "Resource": "arn:aws:s3:::%s/%s%s/*"
            }
            ]
        }
        ''' % (s3_bucket, s3_bucket, s3_folder_key, company_name, s3_bucket,
               s3_folder_key, company_name)
    else:
        iam_policy_document = '''{
                "Version": "2012-10-17",
                "Statement": [
                    {
                        "Sid": "AllowGroupToSeeBucketListInTheConsole",
                        "Effect": "Allow",
                        "Action": [
                            "s3:ListAllMyBuckets",
                            "s3:HeadBucket"
                        ],
                        "Resource": "*"
                    },
                    {
                        "Sid": "AllowListBucketOfASpecificUserPrefix",
                        "Effect": "Allow",
                        "Action": [
                            "s3:ListBucket",
                            "s3:GetBucketLocation"
                        ],
                        "Resource": "arn:aws:s3:::%s",
                        "Condition": {
                            "StringLike": {
                                "s3:prefix": [
                                    "%s%s/*"
                                ]
                            }
                        }
                    },
                    {
                        "Sid": "AllowRead",
                        "Effect": "Allow",
                        "Action": [
                            "s3:GetObject"
                        ],
                        "Resource": "arn:aws:s3:::%s/%s%s/out/*"
                    },
                    {
                        "Sid": "AllowWrite",
                        "Effect": "Allow",
                        "Action": [
                            "s3:PutObject",
                            "s3:DeleteObject"
                        ],
                        "Resource": "arn:aws:s3:::%s/%s%s/in/*"
                    }
        ]
}
                ''' % (s3_bucket, s3_folder_key, company_name,
                       s3_bucket, s3_folder_key, company_name, s3_bucket,
                       s3_folder_key, company_name)

    # Create and write to the temp file.
    with open(temp_aws_policy_file, "w") as my_file:
            my_file.write(iam_policy_document)


# Put the policy into AWS so the new user has access to the new S3 folder.
def put_user_policy_aws():
    print "\nAssigning the IAM policy to the %s user....." % username
    sleep(1)

    # If AWS profile arg is not given use this command. If not, use the other.
    if len(argv) == 1:
        # Create a string to use in the subprocess.call
        call = 'aws iam put-user-policy --user-name %s --policy-name ' \
               'AllowAccess-%s --policy-document file://%s' % \
               (username, company_name, temp_aws_policy_file)

        # Run the string that was created above.
        subprocess.call(call, shell=True)

    else:
        # Create a string to use in the subprocess.call
        call = 'aws --profile %s iam put-user-policy --user-name %s ' \
               '--policy-name AllowAccess-%s --policy-document file://%s' % \
               (aws_profile, username, company_name, temp_aws_policy_file)

        # Run the string that was created above.
        subprocess.call(call, shell=True)


# Delete any temp files created by this script.
def delete_all_temp_files():
    print "\nCleaning the temp files from your system....."
    sleep(1)

    # Create a string to use in the subprocess.call
    subprocess_call_string = 'rm -rf /tmp/AWS_tools_creation*.tmp'

    # Run the string that was created above.
    subprocess.call(subprocess_call_string, shell=True)


# Print out the new user and S3 folder information.
def print_info():
    if verve_partner is False:
        print "\n===========================================================\n"
        print "  -= DATA CREATED FOR THE NEW USER AND S3 FOLDER =-\n"
        print "  -= Manually add a S3 lifecycle to delete files after 45 " \
              "days on the folder: \ns3://%s/%s%s/temp/ .\n" % \
              (s3_bucket, s3_folder_key, folder_name)
        print "Username: %s" % data[0]
        print "Access Key ID: %s" % data[2]
        print "Secret Access Key: %s" % data[1]
        print ""
        print "S3 Folder:"
        print "s3://%s/%s%s/" % (s3_bucket, s3_folder_key, folder_name)
        print ""
    else:
        print "\n===========================================================\n"
        print "  -= DATA CREATED FOR THE NEW USER AND S3 FOLDER =-\n"
        print "Username: %s" % data[0]
        print "Access Key ID: %s" % data[2]
        print "Secret Access Key: %s" % data[1]
        print ""
        print "S3 Upload Folder:"
        print "s3://%s/%s%s/in/" % (s3_bucket, s3_folder_key, folder_name)
        print ""
        print "S3 Download Folder:"
        print "s3://%s/%s%s/out/" % (s3_bucket, s3_folder_key, folder_name)
        print ""


# Run the script.
get_username()
get_aws_usernames()
create_user_list()
compare_users_to_entered_name()
create_aws_user()
create_aws_access_key()
default_or_custom_folder()
create_s3_folder()
create_iam_policy_document()
put_user_policy_aws()
delete_all_temp_files()
print_info()
quit()
