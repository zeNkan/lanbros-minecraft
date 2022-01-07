#!/bin/bash

backed_up_archive="$(date --iso-8601=hour).tar.gz"

tar --force-local -czvf $backed_up_archive lanbros_survival/lanbros_survival lanbros_survival/level.dat
aws s3 cp $backed_up_archive s3://lanbros-mc-backup/

rm -f $backed_up_archive
