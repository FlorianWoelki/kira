runuser -u "$1" -- timeout -s KILL 4 java "$2"