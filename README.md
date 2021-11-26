Binlogo
=====================================
[![Go Reference](https://pkg.go.dev/badge/github.com/jin06/binlogo)](https://pkg.go.dev/github.com/jin06/binlogo)
[![Go Report Card](https://goreportcard.com/badge/github.com/jin06/binlogo)](https://goreportcard.com/report/github.com/jin06/binlogo)
[![codecov](https://codecov.io/gh/jin06/binlogo/branch/master/graph/badge.svg)](https://codecov.io/gh/jin06/binlogo)
![GitHub commit activity](https://img.shields.io/github/commit-activity/m/jin06/binlogo)
</br>
[中文](README_zh.md) | [English](README.md)

Binlogo is the distributed, visualized application based on MySQL binlog.
In short, binlogo is to process the data changes of MySQL into easy to
understand messages and send them to different places according to
the user's configuration. This is some advantages:

* Distributed, multi node improves availability.
* Visualization, can complete common operations and
  observe the status of the whole cluster in the control background

### Get Started

* Install etcd. Binlogo relies on etcd, so you must install etcd first.

* Install binlogo. Binlogo's download address: [Download Address](https://github.com/jin06/binlogo/releases)

* Message Format： [Data format of binlogo output](/docs/1.0.*/message-format.md)

* Start binlogo.
  * Edit config. ${binlogo-path}/configs/binlogo.yaml
    
    ![avatar](/docs/assets/pic/edit_config_step1.en.png)
    
  * > $ ./binlogo server --config ./configs/binlogo.yaml 

* Open browser: http://127.0.0.1:9999/console

* Create Pipeline: 

> Follow the steps. 

![avatar](/docs/assets/pic/create_pipe_step1.en.png)

![avatar](/docs/assets/pic/create_pipe_step2.en.png)

* Run pipeline.

> Click button to run the pipeline instance. 

![avatar](/docs/assets/pic/run_pipeline_step1.en.png)

* Operation condition.

> You can see the operation condition of pipeline.
 
 
![avatar](/docs/assets/pic/pipeline_condition_step1.en.png)

![avatar](/docs/assets/pic/pipeline_condition_step2.en.png)

* See the output 

> Insert some into mysql, watch the ouput on stdout.
 
![avatar](/docs/assets/pic/output_step1.en.png)

![avatar](/docs/assets/pic/output_step2.en.png)

* Configuration output to Kafka

> * High performance, possible data loss. 
>   *  acks=1 
>   *  enable.idempotence=false
>   *  compression.type=snappy
>   *  retries=0
> * For reliability performance: 
>   * acks=-1
>   * enable.idempotence=true
>   * retries=3 or larger one

![avatar](/docs/assets/pic/output_kafka_step1.en.png)

![avatar](/docs/assets/pic/output_kafka_step2.en.png)

### Docker

> $ docker pull jin06/binlogo
> </br>
> $ docker run -e "ETCD_ENDPOINTS=172.17.0.3:2379" --name BinlogoNode -it -d -p 9999:9999 jin06/binlogo:latest 
> </br>
> Open browser access http://127.0.0.1:9999/console
> </br>
> I started five nodes with docker. The following is a screenshot
> 

![avatar](/docs/assets/pic/docker_step1.en.png)


### Docs

* [文档链接](https://github.com/jin06/binlogo/wiki)

### Questions
* To Report bug: [GitHub Issue](https://github.com/jin06/binlogo/issues)
* Contact author: jinlong4696@163.com
