# MQ GO HTTP SDK  
Alyun MQ Documents: http://www.aliyun.com/product/ons

Aliyun MQ Console: https://ons.console.aliyun.com


## Use

1. import `github.com/aliyunmq/mq-http-go-sdk`
2. setup GOPATH:
    - MAC/LINUX: export GOPATH={dir}, export GOBIN=$GOPATH/bin
    - WINDOWS: set GOPATH={dir}
2. go get -x -v

## Note
1. Http consumer only support timer msg (less than 3 days), no matter the msg is produced from http or tcp protocol.
2. Order is only supported at special server cluster.

## Sample (github)

[Publish Message](https://github.com/aliyunmq/mq-http-samples/blob/master/go/producer.go)

[Consume Message](https://github.com/aliyunmq/mq-http-samples/blob/master/go/consumer.go)

[Transaction Message](https://github.com/aliyunmq/mq-http-samples/blob/master/go/trans_producer.go)

[Publish Order Message](https://github.com/aliyunmq/mq-http-samples/blob/master/go/order_producer.go)

[Consume Order Message](https://github.com/aliyunmq/mq-http-samples/blob/master/go/order_consumer.go)


## Sample (code.aliyun.com)

[Publish Message](https://code.aliyun.com/aliware_rocketmq/mq-http-samples/blob/master/go/producer.go)

[Consume Message](https://code.aliyun.com/aliware_rocketmq/mq-http-samples/blob/master/go/consumer.go)

[Transaction Message](https://code.aliyun.com/aliware_rocketmq/mq-http-samples/blob/master/go/trans_producer.go)

[Publish Order Message](https://code.aliyun.com/aliware_rocketmq/mq-http-samples/blob/master/go/order_producer.go)

[Consume Order Message](https://code.aliyun.com/aliware_rocketmq/mq-http-samples/blob/master/go/order_consumer.go)