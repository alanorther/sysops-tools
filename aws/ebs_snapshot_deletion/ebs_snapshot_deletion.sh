#!/bin/bash

# Check the amount of arguments and give usage.
if [ $# -ne 3 ]; then
    printf "\nWrong Amount of Arguments. Requires 3 and got `echo $#`\n"
    printf "USAGE: ebs_snapshot_deletion.sh <ACCOUNT NUMBER> <SNAPSHOT-START-TIME YYYY-MM-DD Format> <DELETE or DRY-RUN all caps>\n"
    printf "Example: ebs_snapshot_deletion.sh 102126644248 2018-12-30 DELETE\n\n"
    exit 1
fi

# Set to dry-run unless the 3rd argument is DELETE
if [ $3 = "DELETE" ]; then
   RUN=""
else
   RUN="--dry-run"
fi

# Create the list of snapshots to delete and list them.
snapshots_to_delete=$(aws ec2 describe-snapshots --owner-ids $1 --query "Snapshots[?StartTime<='$2'].SnapshotId" --output text)
echo "List of snapshots to delete: $snapshots_to_delete"

# Actual Deletion or Dry-run
for snap in $snapshots_to_delete; do
  printf "Running: aws ec2 delete-snapshot --snapshot-id $snap $RUN\n"
  aws ec2 delete-snapshot --snapshot-id $snap $RUN
done
