#!/bin/sh

REG=$1

if [ $# != 1 ]; then
  echo "No region detected, describing all vpcs"
  for reg in us-east-1 us-west-1 us-west-2; do 
    echo "############## $reg ####################"
    aws --output text --region $reg ec2 describe-vpcs | grep Name
  done
else
  aws --output text --region $REG ec2 describe-vpcs | grep Name
fi

