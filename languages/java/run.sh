javac "$2"

dir=$(dirname "$3")
file=$(basename "$3")
runuser -l "$1" -c -- "unshare -n -r timeout -s KILL 10 java -classpath $dir $file"