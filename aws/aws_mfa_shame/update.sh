#!/bin/bash

git add ../.
git commit -m "blah"
git push origin master
go get -u github.com/VerveWireless/sysops-tools/aws/aws_mfa_shame/core
