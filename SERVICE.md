# Setup moneybotsapi Go binary as service

## Pre-requisites

### Make folder for logs

```sh
cd ~
mkdir -p logs/moneybotsapi
cd logs/moneybotsapi
pwd
# /home/ec2-user/logs/moneybotsapi
```

## Create a service

```bash
sudo nano /etc/systemd/system/moneybotsapi.service
```

#### Service File

```bash

[Unit]
Description=Moneybots API
Documentation=https://github.com/nsvirk/moneybotsapi
Wants=network.target
After=network.target

[Service]
Type=simple
User=ec2-user
Group=ec2-user

WorkingDirectory=/home/ec2-user/apps/moneybotsapi
ExecStart=/usr/local/bin/moneybotsapi
Restart=on-failure
RestartSec=10
StandardOutput=file:/home/ec2-user/logs/moneybotsapi/log.log
StandardError=file:/home/ec2-user/logs/moneybotsapi/error.log

[Install]
WantedBy=multi-user.target

```

## Manage the service

```bash
# daemon reload
sudo systemctl daemon-reload

# start the service
sudo systemctl start moneybotsapi
sudo systemctl restart moneybotsapi
sudo systemctl stop moneybotsapi
sudo systemctl status moneybotsapi
sudo systemctl enable moneybotsapi
sudo systemctl disable moneybotsapi

# check logs -r for reverse means current logs
sudo journalctl -u moneybotsapi -r

# check logs -f for follow means live logs
sudo journalctl -u moneybotsapi -f

# find running processes
ps aux | grep -i moneybotsapi

# kill running processes
pkill -f moneybotsapi

# check environment variables
sudo nano /etc/environment
```

## Build the app

```bash

# clone the repo for the first time
cd ~/apps
git clone https://github.com/nsvirk/moneybotsapi.git
cd moneybotsapi

# update the repo every next time
cd ~/apps/moneybotsapi && git pull

# build the app
go build -o moneybotsapi .

# move the binary to /usr/local/bin
sudo mv moneybotsapi /usr/local/bin
cd /usr/local/bin

# delete the binary
sudo rm /usr/local/bin/moneybotsapi
cd /usr/local/bin
```

## Use the Makefile

```bash
# use the make all command to pull, build and deploy the app
make all
```
