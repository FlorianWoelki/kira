#!/bin/bash

# Create runnable users and apply limits.
for i in {0..49}
do
  current_user="user$i"
  # Adds a new user into the runners group for to the directory `/tmp/<current_user>`.
  useradd -m $current_user --base-dir "/tmp/$current_user" -G runners --home "/tmp/$current_user"

  # Set the limits for the user.
  echo "$current_user soft nproc 256" >> /etc/security/limits.conf
  echo "$current_user hard nproc 256" >> /etc/security/limits.conf
  echo "$current_user soft nofile 2048" >> /etc/security/limits.conf
  echo "$current_user hard nofile 2048" >> /etc/security/limits.conf
  echo "$current_user soft fsize 30000000" >> /etc/security/limits.conf
  echo "$current_user hard fsize 30000000" >> /etc/security/limits.conf
  echo "$current_user soft memlock 40960" >> /etc/security/limits.conf
  echo "$current_user hard memlock 40960" >> /etc/security/limits.conf
done
