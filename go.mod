module github.com/jin06/binlogo

go 1.16

require (
	github.com/Shopify/sarama v1.29.1
	github.com/coreos/etcd v3.3.25+incompatible // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/gin-gonic/gin v1.7.2
	github.com/google/uuid v1.2.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/hashicorp/consul/api v1.10.1
	github.com/heetch/confita v0.10.0 // indirect
	github.com/jmoiron/sqlx v1.2.0 // indirect
	github.com/siddontang/go-log v0.0.0-20190221022429-1e957dd83bed
	github.com/siddontang/go-mysql v0.0.0-20200424072754-803944a6e4ea
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.9.0
	go.etcd.io/etcd v3.3.25+incompatible
	sigs.k8s.io/yaml v1.3.0 // indirect

)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
