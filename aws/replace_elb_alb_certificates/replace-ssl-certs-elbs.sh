#!/bin/sh
IFS=$'\n'
#set -x

# Gather base variables
reg=$1
oldcert=$2
newcert=$3
action=$4

# Check basic usage
if [ $# != 4 ]; then
  echo  "Usage: ./script.sh <region> <oldcertarn> <newcertarn> <plan|apply>"
  exit 1
fi

#### ELB Based changes
for elb in `aws --region $reg elb describe-load-balancers --output table | grep LoadBalancerName | awk '{ print $4 }'`; do
  for list in `aws --region $reg elb describe-load-balancers --load-balancer-name $elb --output text | grep LISTENER`; do
    listprot=`echo $list | awk '{ print $5 }'`
    listport=`echo $list | awk '{ print $4 }'`
    certname=`echo $list | awk '{ print $6 }'`
    if [ "$listprot" == "HTTPS" ] && [ "$certname" == "$oldcert" ]; then
      if [ "$action" == "plan" ]; then
          echo "############ ELB Change found ##############"  
          echo "ELB $elb is using cert $certname on listener port $listport"
          cmd="/usr/local/bin/aws --region $reg elb set-load-balancer-listener-ssl-certificate --load-balancer-name $elb --load-balancer-port $listport --ssl-certificate-id $newcert"
          echo "Command that will run: $cmd"
          echo "########################################\n"
      elif [ "$action" == "apply" ]; then
          echo "############ ELB Action taken ##############"
          echo "ELB $elb will be changed to $newcert on listener port $listport"
          cmd="/usr/local/bin/aws --region $reg elb set-load-balancer-listener-ssl-certificate --load-balancer-name $elb --load-balancer-port $listport --ssl-certificate-id $newcert"
          echo "Using command: $cmd"
          aws --region $reg elb set-load-balancer-listener-ssl-certificate --load-balancer-name $elb --load-balancer-port $listport --ssl-certificate-id $newcert
          echo "########################################\n"
      fi
    fi
  done
done

#### ALB Based changes
for alb in `aws --region $reg elbv2 describe-target-groups --output text | grep LOADBALANCERARNS | awk '{ print $2 }'`; do
  for alblist in `aws --region $reg elbv2 describe-listeners --load-balancer-arn $alb --output text | grep LISTENERS`; do
    alblistprot=`echo $alblist | awk '{ print $5 }'`
    alblistport=`echo $alblist | awk '{ print $4 }'`
    alblistarn=`echo $alblist | awk '{ print $2 }'`
    if [ "$alblistprot" == "HTTPS" ]; then
      albcert=`aws --region $reg elbv2 describe-listener-certificates --listener-arn "$alblistarn" --output text | awk '{ print $2 }'`
      if [ "$albcert" == "$oldcert" ]; then
        if [ "$action" == "plan" ]; then
          echo "############ ALB Change found ##############"
          echo "ALB taget group on $alb is using cert $certname on listener port $alblistport"
          albcmd="aws --region $reg elbv2 modify-listener --listener-arn $alblistarn --certificates CertificateArn=$newcert"
          echo "Command that will run: $albcmd"
          echo "########################################\n"
        elif [ "$action" == "apply" ]; then
          echo "############ ALB Action taken ##############"
          echo "ALB target group on $alb will be changed to $newcert on listener port $alblistport"
          albcmd="aws --region $reg elbv2 modify-listener --listener-arn $alblistarn --certificates CertificateArn=$newcert"
          echo "Using command: $albcmd"
          aws --region $reg elbv2 modify-listener --listener-arn $alblistarn --certificates CertificateArn=$newcert
          echo "########################################\n"
        fi
      fi
    fi
  done
done
