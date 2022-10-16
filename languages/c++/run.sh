g++ -Wall -Wextra -Werror -O2 -std=c++17 -pedantic -o "$3" "$2"
runuser -u "$1" -- unshare -n -r timeout -s KILL 10 "$3"
