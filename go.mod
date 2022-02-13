module tag-service

go 1.17

require (
	github.com/elazarl/go-bindata-assetfs v1.0.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.7.3
	github.com/opentracing/opentracing-go v1.1.0
	github.com/sirupsen/logrus v1.8.1
	github.com/soheilhy/cmux v0.1.4
	github.com/uber/jaeger-client-go v2.22.1+incompatible
	go.etcd.io/etcd/client/v3 v3.5.2
	golang.org/x/net v0.0.0-20211123203042-d83791d6bcd9
	google.golang.org/genproto v0.0.0-20220118154757-00ab72f36ad5
	google.golang.org/grpc v1.44.0
	google.golang.org/protobuf v1.27.1
)

require (
	github.com/codahale/hdrhistogram v0.0.0-20161010025455-3a0bb77429bd // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/uber/jaeger-lib v2.2.0+incompatible // indirect
	go.etcd.io/etcd/api/v3 v3.5.2 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.2 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.17.0 // indirect
	golang.org/x/sys v0.0.0-20210603081109-ebe580a85c40 // indirect
	golang.org/x/text v0.3.6 // indirect
	google.golang.org/grpc/examples v0.0.0-20220210231334-75fd0240ac41 // indirect
)

replace github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0
