#!/bin/bash

# Creates a group of users for managing privileges and permissions.
groupadd runners
# Changes the group for runners to the languages directory.
chgrp -R runners /kira/languages/

# Manages control of users running any sort of script.
# Removes read, write, and execute permission for the runners group.
chmod o-rwx /kira/
# Add execute privilege for the runners group.
chmod g+x /kira/languages/
