#!/bin/bash

backed_up_archive="$(date --iso-8601=hour).tar.gz"

tar --force-local -czvf $backed_up_archive server_root/lanbros_survival
/usr/local/aws-cli/v2/current/bin/aws s3 cp $backed_up_archive "s3://lanbros-mc-backup/$(hostname)/"

rm -f $backed_up_archive
