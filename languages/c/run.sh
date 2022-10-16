gcc -Wall -Wextra -Werror -O2 -std=c99 -pedantic -o "$3" "$2"
runuser -l "$1" -c -- "unshare -n -r timeout -s KILL 10 $3"