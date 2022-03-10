protoc --gofast_out=protos/pb protos/pbdef/*.proto --proto_path=protos/pbdef
go run ../../tools/metactl/main.go --module=github.com/jageros/hawox/example/mygame