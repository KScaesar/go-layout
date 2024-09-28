#!/bin/bash
set -euo pipefail

CRONTAB_FILE="/etc/crontab"

# Use sed to find the line with cron.daily and change the time to midnight
if sudo sed -i '/cron\.daily/s/^\([0-9]*\) \([0-9]*\)/0 0/' "$CRONTAB_FILE"; then
    echo "Successfully changed the daily job execution time to 00:00"
    sudo grep 'cron.daily' "$CRONTAB_FILE"
else
    echo "Modification failed. Exiting."
    exit 1
fi

sudo systemctl restart cron
