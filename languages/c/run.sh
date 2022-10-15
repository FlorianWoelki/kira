# TODO: Refactor to more flexible way
runuser -u "$1" -- mkdir /tmp/"$1"/executable
gcc -Wall -Wextra -Werror -O2 -std=c99 -pedantic -o /tmp/"$1"/executable/code "$2"
runuser -u "$1" -- unshare -n -r timeout -s KILL 10 /tmp/"$1"/executable/code