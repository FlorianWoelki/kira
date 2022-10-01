#!/bin/bash

# Create runnable users and apply limits.
for i in {0..49}
do
  current_user="user$i"
  # Adds a new user into the runners group for to the directory `/tmp/<current_user>`.
  useradd -m $current_user -G runners -d "/tmp/$current_user"

  # Set the limits for the user.
  echo "$current_user soft nproc 256" >> /etc/security/limits.conf
  echo "$current_user hard nproc 256" >> /etc/security/limits.conf
  echo "$current_user soft nofile 4096" >> /etc/security/limits.conf
  echo "$current_user hard nofile 4096" >> /etc/security/limits.conf
  echo "$current_user soft fsize 10120" >> /etc/security/limits.conf
  echo "$current_user hard fsize 10120" >> /etc/security/limits.conf
  echo "$current_user soft memlock 40960" >> /etc/security/limits.conf
  echo "$current_user hard memlock 40960" >> /etc/security/limits.conf
done
