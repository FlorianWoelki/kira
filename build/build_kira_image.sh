set GOARCH=amd64
set GOOS=linux
go build -o kira ../main.go

docker build --no-cache -t kira:v0.0.1 .
docker tag kira:v0.0.1 florianwoelki/kira
docker login
docker push florianwoelki/kira
