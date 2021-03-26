#!/bin/bash

service nscd start

exec /root/paste -f /root/config.yaml
