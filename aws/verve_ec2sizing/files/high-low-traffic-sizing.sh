#!/bin/sh

role=$1
traffic=$2

USAGE="$0 <dsh role group> <traffic pattern (high|low)>"

if [ $# != 2 ]; then
  echo "Number of arguments wrong, see usage:"
  echo $USAGE
  exit 1
fi

# HIGH/LOW config for RTB-BIDDER

if [ $role == "role-rtb-bidder" ]; then
  #nodes="i-7ce4f4d5 i-49e4f4e0 i-72e4f4db i-74e4f4dd i-d5d56406 i-d8d5640b i-d2d56401 i-d6d56405 i-e92a2269 i-962b2316 i-c616eb42 i-f416eb70"
  nodes="i-7ce4f4d5 i-49e4f4e0 i-72e4f4db i-74e4f4dd i-d5d56406 i-d8d5640b i-d2d56401 i-d6d56405"
  if [ $traffic == "low" ]; then
    aws elb deregister-instances-from-load-balancer --load-balancer-name vpc-prod-rtb-0 --instances $nodes
    sleep 300
    aws ec2 stop-instances --instance-ids $nodes
    sleep 120
    puppet agent -t
    su - deploy -c "dsh -w dns1 sudo puppet agent -t"
  elif [ $traffic == "high" ]; then
    aws ec2 start-instances --instance-ids $nodes
    sleep 300
    aws elb register-instances-with-load-balancer --load-balancer-name vpc-prod-rtb-0 --instances $nodes
    sleep 120
    puppet agent -t
    su - deploy -c "dsh -w dns1 sudo puppet agent -t"
  fi
fi

#HIGH/LOW config for RTB-LOG

if [ $role == "role-rtb-log" ]; then
  #nodes="i-ff0a1b7f i-fa0a1b7a i-c80d1c48 i-507801d0 i-294fb4ad i-214fb4a5 i-5e0f55e8 i-0e0c56b8"
  nodes="i-ff0a1b7f i-fa0a1b7a i-c80d1c48 i-214fb4a5 i-5e0f55e8 i-0e0c56b8"
  if [ $traffic == "low" ]; then
    #for line in `aws ec2 describe-instances --output text --instance-ids $nodes | egrep 'TAGS.host'|awk '{ print $3 }'`;do su - deploy -c "ssh $line 'sudo /sbin/service supervisord stop'";done 
    #sleep 600
    aws ec2 stop-instances --instance-ids $nodes
    sleep 120
    puppet agent -t
    su - deploy -c "dsh -w dns1 sudo puppet agent -t"
  elif [ $traffic == "high" ]; then
    aws ec2 start-instances --instance-ids $nodes
    sleep 120
    puppet agent -t
    su - deploy -c "dsh -w dns1 sudo puppet agent -t"
  fi
fi

#HIGH/LOW config for ADCEL

if [ $role == "role-adcel" ]; then
  nodes="i-e6154157 i-9d15412c i-9a15412b i-f1f11678 i-fdf11674 i-fef11677"
  if [ $traffic == "low" ]; then
    aws elb deregister-instances-from-load-balancer --load-balancer-name adcel --instances $nodes
    sleep 600
    aws ec2 stop-instances --instance-ids $nodes
  elif [ $traffic == "high" ]; then
    aws ec2 start-instances --instance-ids $nodes
    sleep 600
    aws elb register-instances-with-load-balancer --load-balancer-name adcel --instances $nodes
  fi
fi

#HIGH/LOW config for NSQHUB

if [ $role == "role-nsqhub" ]; then
  #nodes="i-258b8df1 i-248b8df0 i-4decc4f2 i-4cecc4f3"
  nodes="i-4cecc4f3 i-5f669cdb i-bb669c3f i-150b1b95"
  if [ $traffic == "low" ]; then
    aws ec2 stop-instances --instance-ids $nodes
    sleep 120
    puppet agent -t
    su - deploy -c "dsh -w dns1 sudo puppet agent -t"
  elif [ $traffic == "high" ]; then
    aws ec2 start-instances --instance-ids $nodes
    sleep 120
    puppet agent -t
    su - deploy -c "dsh -w dns1 sudo puppet agent -t"
  fi
fi

#HIGH/LOW config for FCAP

if [ $role == "role-rtb-fcap" ]; then
  nodes="i-512f50ee i-0442acb5 i-a643ad17 i-675834d6 i-2d59359c"
  if [ $traffic == "low" ]; then
    aws ec2 stop-instances --instance-ids $nodes
  elif [ $traffic == "high" ]; then
    aws ec2 start-instances --instance-ids $nodes
  fi
fi
