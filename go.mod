module github.com/jin06/binlogo

go 1.16

replace github.com/coreos/bbolt v1.3.6 => go.etcd.io/bbolt v1.3.6

replace google.golang.org/grpc v1.42.0 => google.golang.org/grpc v1.26.0

require (
	github.com/Shopify/sarama v1.30.0
	github.com/coreos/bbolt v1.3.6 // indirect
	github.com/coreos/etcd v3.3.27+incompatible
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/gin-gonic/gin v1.7.4
	github.com/go-mysql-org/go-mysql v1.3.0
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/hashicorp/consul/api v1.11.0
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/prometheus/client_golang v1.11.0 // indirect
	github.com/shirou/gopsutil/v3 v3.21.10
	github.com/siddontang/go-log v0.0.0-20190221022429-1e957dd83bed // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/soheilhy/cmux v0.1.5 // indirect
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.9.0
	github.com/tmc/grpc-websocket-proxy v0.0.0-20201229170055-e5319fda7802 // indirect
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	google.golang.org/grpc v1.42.0 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)
