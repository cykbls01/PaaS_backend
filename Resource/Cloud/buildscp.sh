echo "building"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
echo "build finished"
scp -r -P 2226 $GOPATH/src/Cloud root@al.sijia.ooo:~/Cloud1