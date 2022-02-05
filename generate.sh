# 生成普通序列化
protoc -I=/mnt/d/gopath/src/ -I=. --go_out=. ./*.proto
# 生成grpc
protoc -I=/mnt/d/gopath/src/ -I=. --go-grpc_out=. ./tag.proto
# 生成gateway
protoc -I=/mnt/d/gopath/src/ -I=. --grpc-gateway_out=logtostderr=true:. ./tag.proto