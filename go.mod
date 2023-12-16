module github.com/sinisa-andric/meridian

go 1.21.5

require (
	github.com/c12s/meridian v0.0.0-20200522172917-981b013b70ff
	github.com/c12s/scheme v0.0.0-20200211232926-6490b386ab3f
	github.com/c12s/stellar-go v0.0.0-20191220161710-a82c2c7bb52e
	go.etcd.io/etcd v3.3.27+incompatible
	golang.org/x/net v0.17.0
	google.golang.org/grpc v1.59.0
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/coreos/bbolt v1.3.8 // indirect
	github.com/coreos/etcd v3.3.27+incompatible // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/go-redis/redis v6.15.7+incompatible // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/btree v1.1.2 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/jonboulle/clockwork v0.4.0 // indirect
	github.com/nats-io/gnatsd v1.4.1 // indirect
	github.com/nats-io/nats-server v1.4.1 // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.30.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_golang v1.11.1 // indirect
	github.com/soheilhy/cmux v0.1.5 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
	github.com/tmc/grpc-websocket-proxy v0.0.0-20220101234140-673ab2c3ae75 // indirect
	github.com/xiang90/probing v0.0.0-20221125231312-a49e3df8f510 // indirect
	golang.org/x/time v0.5.0 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)

replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.5

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/MilosSimic/pipes v0.0.0-20191204181555-9ca78dc556f9 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/coreos/go-semver v0.3.1 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/uuid v1.3.1 // indirect
	github.com/nats-io/go-nats v1.7.3-0.20190608183121-73ffc26dfe70 // indirect
	github.com/nats-io/nkeys v0.1.0 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/redis/go-redis/v9 v9.3.0
	github.com/rs/xid v1.2.1 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.18.1 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/genproto v0.0.0-20231002182017-d307bd883b97 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20231002182017-d307bd883b97 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231002182017-d307bd883b97 // indirect
	google.golang.org/protobuf v1.31.0
)
