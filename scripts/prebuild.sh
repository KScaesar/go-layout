#!/bin/bash
set -euo pipefail

USER="caesar"
GROUP="caesar"
SERVICE="myapp"
WORK_DIR="/home/$USER/$SERVICE"

# log
sudo mkdir -p /var/log/$SERVICE
sudo touch /var/log/$SERVICE/stdout.log /var/log/$SERVICE/stderr.log
sudo chown -R $USER:$GROUP /var/log/$SERVICE

# systemd
if [ -f /etc/systemd/system/$SERVICE.service ]; then
  sudo rm /etc/systemd/system/$SERVICE.service
fi
sudo ln -s $WORK_DIR/$SERVICE.service /etc/systemd/system/$SERVICE.service

# logrotate
if [ -f /etc/logrotate.d/$SERVICE.conf ]; then
  sudo rm /etc/logrotate.d/$SERVICE.conf
fi
sudo ln -s $WORK_DIR/$SERVICE.conf /etc/logrotate.d/$SERVICE.conf
sudo chown root:$GROUP $WORK_DIR/$SERVICE.conf

# crontab
# Use sed to find the line with cron.daily and change the time to midnight
if sudo sed -i '/cron\.daily/s/^\([0-9]*\) \([0-9]*\)/0 0/' /etc/crontab; then
    echo "Successfully changed the daily job execution time to 00:00"
    sudo grep 'cron.daily' /etc/crontab
else
    echo "Modification failed. Exiting."
    exit 1
fi

# update service
sudo systemctl daemon-reload
sudo systemctl restart logrotate.service
sudo systemctl restart cron.service
sudo systemctl restart $SERVICE.service
