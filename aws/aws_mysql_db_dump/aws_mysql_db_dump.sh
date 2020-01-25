#!/bin/bash

MYUSERNAME=vlsm_admin
MYPASSWORD=
MYHOST=vpc-prod-pois-db-0.cpqtuzct7fyu.us-east-1.rds.amazonaws.com
MYSQLDUMP=`which mysqldump`

MYSQL=$(mysql -N -u${MYUSERNAME} -p${MYPASSWORD} -h${MYHOST} <<<"SHOW DATABASES" | grep -v mysql | grep -v information_schema | grep -v performance_schema | tr "\n" " ")

${MYSQLDUMP} -v -u${MYUSERNAME} -p${MYPASSWORD} -h${MYHOST} --hex-blob --routines --triggers --events --skip-add-drop-table --databases --skip-lock-tables ${MYSQL} > DB-Dump.sql
