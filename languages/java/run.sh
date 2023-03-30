dir=$(dirname "$3")
file=$(basename "$3")
# runuser -l "$1" -c -- "unshare -n -r timeout -s KILL 10 java -classpath $dir $file"
runuser -l "$1" -c -- "nice timeout -s KILL 10 prlimit --nproc=64 --nofile=2048 --fsize=10000000 blocksyscalls $dir $file"