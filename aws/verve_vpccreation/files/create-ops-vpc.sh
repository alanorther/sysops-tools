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


TEMPDIR="../../verve_aws/manifests/template"
STACKDIR="../../verve_aws/manifests/$STACKNAME"
JSONFILE="../templates/vpc-$STACKNAME-$ENV-template.json"


if [ $ENV == "prod" ]; then
  AZ="us-east-1c"
  REG=`echo $AZ | sed 's/.$//'`
  AMIID="ami-e9357083"
  KEYNAME="vrv-prod-mgmt"
  JSONTEMP="../templates/vpc-virg-template.json"
elif [ $ENV == "stage" ]; then
  AZ="us-west-2c"
  REG=`echo $AZ | sed 's/.$//'`
  AMIID="ami-0f21323f"
  KEYNAME="js-ssh"
  JSONTEMP="../templates/vpc-oreg-template.json"
elif [ $ENV == "uk" ]; then
  AZ="eu-west-1a"
  REG=`echo $AZ | sed 's/.$//'`
  AMIID="ami-89a320fa"
  KEYNAME="uk-ssh"
  JSONTEMP="../templates/vpc-uk-template.json"
else
  echo "Invalid Environment. See usage."
  echo "USAGE: $USAGE"
  exit 1
fi

if [ $ENV == "uk" ]; then
  sed "s/STACKNAME/$STACKNAME/g; s/NETWORK/$NETWORK/g; s/SECOCT/$SECOCT/g; s/AZ/$AZ/g; s/AMIID/$AMIID/g; s/KEYNAME/$KEYNAME/g; s/REG/$REG/g" $JSONTEMP > $JSONFILE
else
  sed "s/STACKNAME/$STACKNAME/g; s/NETWORK/$NETWORK/g; s/SECOCT/$SECOCT/g; s/AZ/$AZ/g; s/NAMESRV/$NAMESRV/g; s/PEERROUTE/$PEERROUTE/g; s/AMIID/$AMIID/g; s/KEYNAME/$KEYNAME/g; s/PEERVPC/$PEERVPC/g; s/ROUTEDEST/$ROUTEDEST/g; s/REG/$REG/g" $JSONTEMP > $JSONFILE
fi

echo "Copying cloud formation config to s3..."
#s3cmd put $JSONFILE s3://verve-cf-templates/$ENV/

  echo "Stack is auto deploying with the following command:"
  echo "aws --region $REG cloudformation create-stack --stack-name $STACKNAME --template-url https://s3.amazonaws.com/verve-cf-templates/$ENV/vpc-$STACKNAME-$ENV-template.json"
  echo "In the future if you wish to redepoly you can use the following URL:"
  echo "https://s3.amazonaws.com/verve-cf-templates/$ENV/vpc-$STACKNAME-$ENV-template.json"

