module github.com/aschmidt75/wgmesh

go 1.15

replace github.com/aschmidt75/wgmesh/cmd => ./cmd

require (
	github.com/armon/go-metrics v0.4.1 // indirect
	github.com/cristalhq/jwt/v3 v3.1.0
	github.com/google/btree v1.1.2 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-msgpack v0.5.5 // indirect
	github.com/hashicorp/go-msgpack/v2 v2.1.2 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-sockaddr v1.0.6 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/hashicorp/memberlist v0.5.1
	github.com/hashicorp/serf v0.10.1
	github.com/mdlayher/socket v0.5.1 // indirect
	github.com/miekg/dns v1.1.59 // indirect
	github.com/sirupsen/logrus v1.9.3
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	golang.org/x/tools v0.20.0 // indirect
	golang.zx2c4.com/wireguard v0.0.0-20231211153847-12269c276173 // indirect
	golang.zx2c4.com/wireguard/wgctrl v0.0.0-20230429144221-925a1e7659e6
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240429193739-8cf5692501f6 // indirect
	google.golang.org/grpc v1.63.2
	google.golang.org/protobuf v1.34.0
	gopkg.in/yaml.v2 v2.4.0
	gortc.io/stun v1.23.0
)
