#!/bin/sh

USAGE="./createvpc.sh STACKNAME NETWORK ENVIRONMENT #where STACKNAME is the name of the new vpc, NETWORK is a 4 octet network ie 10.10.0.0, and ENVIRONMENT is either prod or stage"

if [ $# != 3 ]; then
  echo "Not enough variables. Please check usage."
  echo "USAGE: $USAGE"
  exit 1
fi

STACKNAME=$1
NETWORK=$2
ENV=$3

SECOCT=`echo $NETWORK | cut -d. -f2`

JSONFILE="../templates/vpc-$STACKNAME-$ENV-template.json"

if [ $ENV == "prod" ]; then
  AZ="us-east-1c"
  REG=`echo $AZ | sed 's/.$//'`
  NAMESRV="10.1.9.177"
  PEERROUTE="rtb-9ac1d0f8"
  AMIID="ami-e9357083"
  KEYNAME="vrv-prod-mgmt"
  PEERVPC="vpc-ab271cc0"
  ROUTEDEST="10.1.0.0"
  JSONTEMP="../templates/vpc-template.json"
elif [ $ENV == "stage" ]; then
  AZ="us-west-2c"
  REG=`echo $AZ | sed 's/.$//'`
  NAMESRV="10.61.1.103"
  PEERROUTE="rtb-b3cbb4d6"
  AMIID="ami-0f21323f"
  KEYNAME="js-ssh"
  PEERVPC="vpc-6583f500"
  ROUTEDEST="10.61.0.0"
  JSONTEMP="../templates/vpc-template.json"
elif [ $ENV == "uk" ]; then
  AZ="eu-west-1a"
  REG=`echo $AZ | sed 's/.$//'`
  NAMESRV="10.200.1.233"
  PEERROUTE="rtb-165d8a72"
  AMIID="ami-89a320fa"
  KEYNAME="uk-ssh"
  PEERVPC="vpc-835ebae7"
  ROUTEDEST="10.200.0.0"
  JSONTEMP="../templates/vpc-template.json"
else
  echo "Invalid Environment. See usage."
  echo "USAGE: $USAGE"
  exit 1
fi

sed "s/STACKNAME/$STACKNAME/g; s/NETWORK/$NETWORK/g; s/SECOCT/$SECOCT/g; s/AZ/$AZ/g; s/NAMESRV/$NAMESRV/g; s/PEERROUTE/$PEERROUTE/g; s/AMIID/$AMIID/g; s/KEYNAME/$KEYNAME/g; s/PEERVPC/$PEERVPC/g; s/ROUTEDEST/$ROUTEDEST/g; s/REG/$REG/g" $JSONTEMP > $JSONFILE

echo "Copying cloud formation config to s3..."
s3cmd put $JSONFILE s3://verve-cf-templates/$ENV/

  echo "Stack is auto deploying with the following command:"
  echo "aws --region $REG cloudformation create-stack --stack-name $STACKNAME --template-url https://s3.amazonaws.com/verve-cf-templates/$ENV/vpc-$STACKNAME-$ENV-template.json"
  echo "In the future if you wish to redepoly you can use the following URL:"
  echo "https://s3.amazonaws.com/verve-cf-templates/$ENV/vpc-$STACKNAME-$ENV-template.json"

  aws --region $REG cloudformation create-stack --stack-name $STACKNAME --template-url https://s3.amazonaws.com/verve-cf-templates/$ENV/vpc-$STACKNAME-$ENV-template.json

echo "Sleeping 60s for initial vpc creation"
echo "This ensures that the vpc is there when we add the additional security group"
sleep 60

echo "Adding the correct vpc-$STACKNAME-allow-all sd"
# Get VPC ID
VPCID=`aws --output text --region $REG ec2 describe-vpcs --filter Name=tag-value,Values=vpc-$STACKNAME | grep vpc | awk '{ print $7 }' | head -1`
echo "VPCID: $VPCID"

if [ -z $VPCID ]; then
  echo "No VPCID found"
  exit 1
fi

SGID=`aws --output text --region $REG ec2 create-security-group --group-name vpc-$STACKNAME-allow-all --description "Allow all security group for $STACKNAME" --vpc-id $VPCID`
echo "Security Group ID: $SGID"

if [ -z $SGID ]; then
  echo "Nof SGID Found"
  exit 1
fi

aws --output text --region $REG ec2 authorize-security-group-ingress --group-id $SGID --protocol all --port all --cidr 0.0.0.0/0
aws --output text --region $REG ec2 create-tags --resources $SGID --tags Key=Name,Value=vpc-$STACKNAME-allow-all

echo "Now we need to disable the source and destination checking on the gateway node we created"

for line in `aws --output text --region $REG ec2 describe-instances --filter Name=vpc-id,Values=$VPCID | grep NETWORK | awk -F"\t" '{ print $4}'`; do
	aws --output text --region $REG ec2 modify-network-interface-attribute --network-interface-id $line --no-source-dest-check
done
echo "All Done! Enjoy."
