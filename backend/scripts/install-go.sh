#!/bin/bash

# fail immediately if error
set -e

# lock global
touch /tmp/install.lock

# lock
touch /tmp/install-go.lock

# install golang
apt-get update
apt-get install -y golang

# unlock global
rm /tmp/install.lock

# unlock
rm /tmp/install-go.lock
