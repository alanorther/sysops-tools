#!/bin/sh

echo "Syncing Staging..."
ssh puppet.ops.us-west-2.aws.vrv -C "cd /etc/puppet/environments/master/ && git pull origin master && librarian-puppet install"
echo "Syncing Prod West..."
ssh puppet.prod.us-west-1.aws.vrv -C "cd /etc/puppet/environments/master/ && git pull origin master && librarian-puppet install"
echo "Syncing Prod East..."
ssh puppet.prod.us-east-1.aws.vrv -C "cd /etc/puppet/environments/master/ && git pull origin master && librarian-puppet install"
echo "Syncing Prod UK..."
ssh puppet.ops.eu-west-1.aws.vrv -C "cd /etc/puppet/environments/master/ && git pull origin master && librarian-puppet install"
