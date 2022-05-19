protoc --gofast_out=protos/pb protos/pbdef/*.proto --proto_path=protos/pbdef
go run ../../tools/metactl/main.go --module=git.hawtech.cn/jager/hawox/example/mygame