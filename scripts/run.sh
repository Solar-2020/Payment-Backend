#!/bin/sh
exec_file=$1
log_filename=$2
touch $log_filename
#$exec_file | xargs -I{} echo -e "[$(date +'%d-%m-%Y %H:%M:%S')] " {} >> $log_filename  # causes container death on load
$exec_file >> $log_filename
