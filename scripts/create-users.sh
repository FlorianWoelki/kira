#!/bin/bash

# Create runnable users and apply limits.
for i in {1..3}
do
  current_user="user$i"
  # Adds a new user into the runners group for to the directory `/tmp/<current_user>`.
  useradd -m $current_user -G runners -d "/tmp/$current_user"

  # Set the limits for the user.
  echo "$current_user soft nproc 64" >> /etc/security/limits.conf
  echo "$current_user hard nproc 64" >> /etc/security/limits.conf
  echo "$current_user soft nofile 2048" >> /etc/security/limits.conf
  echo "$current_user hard nofile 2048" >> /etc/security/limits.conf
  echo "$current_user soft fsize 5120" >> /etc/security/limits.conf
  echo "$current_user hard fsize 5120" >> /etc/security/limits.conf
  echo "$current_user soft memlock 40960" >> /etc/security/limits.conf
  echo "$current_user hard memlock 40960" >> /etc/security/limits.conf
done
