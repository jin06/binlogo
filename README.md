Binlogo
=====================================
[![Go Reference](https://pkg.go.dev/badge/github.com/jin06/binlogo)](https://pkg.go.dev/github.com/jin06/binlogo)
[![Go Report Card](https://goreportcard.com/badge/github.com/jin06/binlogo)](https://goreportcard.com/report/github.com/jin06/binlogo)
[![codecov](https://codecov.io/gh/jin06/binlogo/branch/master/graph/badge.svg)](https://codecov.io/gh/jin06/binlogo)
![GitHub commit activity](https://img.shields.io/github/commit-activity/m/jin06/binlogo)
</br>
[中文](README_zh.md) | [English](README.md)

Binlogo is the distributed, visualized application based on MySQL binlog. In short, binlogo is to process the data
changes of MySQL into easily understand messages and output to different systems according to the user's
configuration. Here are part of advantages:

* Distributed, multi node improves availability of the whole system.
* Visualization, can complete common operations and observe the status of the whole cluster in the control website 
* It can be output to multiple queues or other applications, and new outputs are constantly added

### Get Started

* Install redis. Binlogo V2 relies on redis, so you must install redis first.

* Install binlogo. Binlogo's download address: [Download Address](https://github.com/jin06/binlogo/releases)

* Message Format： [Data format of binlogo output](/docs/1.0.*/message-format.md)

* Start binlogo.

   Edit config. ${binlogo-path}/etc/binlogo.yaml
  
  
      # ./binlogo server --config ./etc/binlogo.yaml


* Open up your browser to http://127.0.0.1:8081/console to view the console website

* Create Pipeline:
[config.go](etc%2Fconfig.go)
> Follow the steps.

![avatar](/docs/wiki/assets/pic/create_pipe_step1.en.png)

![avatar](/docs/wiki/assets/pic/create_pipe_step2.en.png)

* Run pipeline.

> Click button to run the pipeline instance.

![avatar](/docs/wiki/assets/pic/run_pipeline_step1.en.png)

* Operation condition.

> You can see the operation condition of pipeline.


![avatar](/docs/wiki/assets/pic/pipeline_condition_step1.en.png)

![avatar](/docs/wiki/assets/pic/pipeline_condition_step2.en.png)

* See the output

> Insert some into mysql, watch the ouput on stdout.

![avatar](/docs/wiki/assets/pic/output_step1.en.png)

![avatar](/docs/wiki/assets/pic/output_step2.en.png)

* Configuration output to Kafka

> * High performance, possible data loss.
>   * acks=1
>   * enable.idempotence=false
>   * compression.type=snappy
>   * retries=0
> * For reliability performance:
>   * acks=-1
>   * enable.idempotence=true
>   * retries=3 or larger one

![avatar](/docs/wiki/assets/pic/output_kafka_step1.en.png)

![avatar](/docs/wiki/assets/pic/output_kafka_step2.en.png)

### Docker

 [Docker Hub](https://hub.docker.com/r/jin06/binlogo)


      # docker pull jin06/binlogo
      # docker run -e REDIS_ADDR=127.0.0.1 -e REDIS_PORT=6379 -e REDIS_PASSWORD=yourpassword -e REDIS_DB=0  --name BinlogoNode -it -d -p 8081:8081 jin06/binlogo:latest



![avatar](/docs/wiki/assets/pic/docker_step1.en.png)

### Kubernetes

- [doc](/docs/wiki/instanll-kubernetes.md)


### Other outputs

* [HTTP](/docs/wiki/configure-http-output.md)
* [RabbitMQ](/docs/wiki/configure-rabbitmq-outupt.md)
* [Kafka](/docs/wiki/configure-kafka-output.md)
* [Redis](/docs/wiki/configure-redis-outupt.md)
* [AliCloud RocketMQ](/docs/wiki/configure-rocketmq-outupt.md)

### Docs

* [docs link](/docs/wiki)

### Questions

* To Report bug: [GitHub Issue](https://github.com/jin06/binlogo/issues)
* Contact author: jlonmyway@gmail.com
