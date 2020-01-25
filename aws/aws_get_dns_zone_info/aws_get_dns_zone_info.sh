#!/bin/bash
# Run the command with the domain as an argument. Ex: ./aws_get_dns_zone_info.sh website.com

zonename=$1
hostedzoneid=$(aws route53 list-hosted-zones | jq -r ".HostedZones[] | select(.Name == \"$zonename.\") | .Id" | cut -d'/' -f3)
aws route53 list-resource-record-sets --hosted-zone-id $hostedzoneid --output json | jq -jr '.ResourceRecordSets[] | "\(.Name) \t\(.TTL) \t\(.Type) \t\(.ResourceRecords[].Value)\n"'
