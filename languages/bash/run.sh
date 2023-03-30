#runuser -l "$1" -c -- "unshare -n -r timeout -s KILL 10 /bin/bash $2"
runuser -l "$1" -c -- "nice timeout -s KILL 3 prlimit --nproc=64 --nofile=2048 --fsize=10000000 blocksyscalls /bin/bash $2"