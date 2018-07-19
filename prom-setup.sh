#! /bin/bash

apt-get install wget tar systemd -y

wget https://github.com/prometheus/prometheus/releases/download/v2.3.2/prometheus-2.3.2.linux-amd64.tar.gz \
    -O /tmp/prom.tar.gz

tar -xvzf /tmp/prom.tar.gz
sudo ln -s /prometheus-2.3.2.linux-amd64/prometheus /usr/bin

tee config.yml <<-EOF
global:
    scrape_interval: 10s
EOF

tee /etc/systemd/system/prometheus.service <<-EOF
[Unit]
Description=prometheus
Documentation=https://prometheus.io/docs/introduction/overview/
After=network-online.target
Wants=network-online.target

[Service]
ExecStart=/bin/sh prometheus --config-file=/config.yml
Restart=always

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable prometheus
sudo systemctl start prometheus
