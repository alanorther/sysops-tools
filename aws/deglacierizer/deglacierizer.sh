#!/bin/bash

# Confirm that the S3 path and restoration days arguments are set.
[[ $# != 2 ]] && echo "Usage: deglacierize <s3 path> <number of days to restore>" && exit 1
S3PATH=$1
DAYS_TO_RESTORE=$2
[[ "${S3PATH: -1}" == "/" ]] || S3PATH="${S3PATH}/"

# Confirm that the S3 path is correct and set bucket variable.
if [[ "$S3PATH" =~ ^s3://([^/]+)/(.*)?/$ ]]; then
    BUCKET=${BASH_REMATCH[1]}
else
    echo "S3 path must be in the format 's3://bucket/path'"
fi

# Set temp file names and locations.
TMPFILE1=/tmp/deglacierize.$$.1
TMPFILE2=/tmp/deglacierize.$$.2
TMPFILE3=/tmp/deglacierize.$$.3
TMPFILE4=/tmp/deglacierize.$$.4

quit() {
    rm -f $TMPFILE1 $TMPFILE2 $TMPFILE3 $TMPFILE4
    exit $1
}

printf "Getting recursive list of files in ${S3PATH}.......\n\n"
aws s3 ls --recursive "${S3PATH}"  > $TMPFILE1
[[ $? != 0 ]] && echo "Error when querying S3; quitting" && quit 1

printf "Pulling S3 keys from the list of S3 objects.......\n\n"
sleep 1
awk '{if ($4) print $4}' $TMPFILE1 > $TMPFILE2

printf "Finding all files that are in Glacier.....\n\n"
while read KEY; do
    if aws s3api list-objects --bucket "$BUCKET" --prefix "$KEY" --query 'Contents[].{StorageClass: StorageClass}' | grep GLACIER > /dev/null
        then echo "$KEY" >> "$TMPFILE3"
    [[ $? != 0 ]] && echo "Error!" && quit 1
    fi;
done < $TMPFILE2

if [ ! -f "$TMPFILE3" ]; then
    printf "There are no files in Glacier within this S3 location.\n\n"
    quit 0
fi;

printf "Filtering files that are already restored or being restored."
while read KEY; do
    if aws s3api head-object --bucket "$BUCKET" --key "$KEY" | grep Restore
      then :
    else
      echo "$KEY" >> "$TMPFILE4"
    fi;
done < $TMPFILE3

NUM_FILES_TEMP2=$(wc -l $TMPFILE2 | cut -d" " -f1)
NUM_FILES_TEMP3=$(wc -l $TMPFILE3 | cut -d" " -f1)
NUM_FILES_TEMP4=$(wc -l $TMPFILE4 | cut -d" " -f1)

### Used for testing - Uncomment these 2 lines and comment out everything below except for the quit line.
#printf "Restoring ${NUM_FILES_TEMP3} files from Glacier. Originally ${NUM_FILES_TEMP2} files in S3 list. Restoring for ${DAYS_TO_RESTORE} days.\n\n"

read -p "Restoring ${NUM_FILES_TEMP4} files from original ${NUM_FILES_TEMP2} in list for ${DAYS_TO_RESTORE} days; proceed? [y|N]: " go
[[ $go != "y" ]] && echo "Aborting!" && quit 1

while read KEY; do
    echo -n "Restoring s3://$BUCKET/$KEY... "
    aws s3api restore-object \
        --bucket "$BUCKET" \
        --key "$KEY" \
    --restore-request "Days=$DAYS_TO_RESTORE"
    [[ $? != 0 ]] && echo "Error!" && quit 1
    echo "ok!"
done < $TMPFILE4

quit 0
