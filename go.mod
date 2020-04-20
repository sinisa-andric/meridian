module github.com/c12s/meridian

replace github.com/coreos/go-systemd/journal => ../../coreos/go-systemd/journal

go 1.13

require (
	github.com/c12s/apollo v0.0.0-20191225203730-f26e394a8564
	github.com/c12s/scheme v0.0.0-20200211232926-6490b386ab3f
	github.com/c12s/stellar-go v0.0.0-20191220161710-a82c2c7bb52e
	github.com/coreos/etcd v3.3.18+incompatible
	github.com/coreos/go-systemd/journal v0.0.0-00010101000000-000000000000 // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/protobuf v1.3.3
	github.com/google/uuid v1.1.1 // indirect
	go.uber.org/zap v1.13.0 // indirect
	golang.org/x/net v0.0.0-20200202094626-16171245cfb2
	google.golang.org/grpc v1.25.1
	gopkg.in/yaml.v2 v2.2.8
)
