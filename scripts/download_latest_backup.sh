#!/bin/bash

aws s3 cp "s3://lanbros-mc-backup$(aws s3 ls s3://lanbros-mc-backup/ | awk -F' ' '{print $NF}' | head -n1)" .
