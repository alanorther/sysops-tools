#!/bin/bash

FILE=ip_range_file.txt

for ip in 10.11.{120..127}.{0..255}
do
  IP_WITH_DASHES=$(echo $ip | sed 's/\./-/g')
  echo "=ip-$IP_WITH_DASHES.diprod.us-east-1.aws.vrv:$ip" >> $FILE
  #EXAMPLE DNS# =ip-10-2-85-53.diprod.us-east-1.aws.vrv:10.2.85.53
done
