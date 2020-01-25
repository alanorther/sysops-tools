#!/bin/sh

REG=$1
VPCNAME=$2

if [ $# != 2 ]; then
  echo "Either too many or too few arguments, see usage"
  echo "USAGE: ./get-vpc-id.sh REGION VPCNAME"
  exit 1
fi

# Get VPC ID
#VPCID=`aws --output text --region $REG ec2 describe-vpcs --filter Name=tag-value,Values=ads-$STACKNAME-vpc-$SECOCT | grep vpc | awk '{ print $7 }' | head -1`
VPCID=`aws --output text --region $REG ec2 describe-vpcs --filter Name=tag-value,Values=$VPCNAME | grep VPCS | awk '{ print $7 }'`


#echo "VPCID: $VPCID"

if [ -z $VPCID ]; then
  echo "No VPCID found"
  echo "Please make sure your REGION is correct and your VPCNAME fits one of the following:"
  ./get-vpc-names.sh
  exit 1
fi

echo $VPCID

#SGID=`aws --output text --region $REG ec2 create-security-group --group-name vpc-$STACKNAME-allow-all --description "Allow all security group for $STACKNAME" --vpc-id $VPCID`
#echo "Security Group ID: $SGID"

#if [ -z $SGID ]; then
#  echo "Nof SGID Found"
#  exit 1
#fi

#aws --output text --region $REG ec2 authorize-security-group-ingress --group-id $SGID --protocol all --port all --cidr 0.0.0.0/0
#aws --output text --region $REG ec2 create-tags --resources $SGID --tags Key=Name,Value=vpc-$STACKNAME-allow-all


