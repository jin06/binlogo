Binlogo
=====================================
[中文](README_zh.md) | [English](README.md)

Binlogo是一个基于MySQL的binlog开发的数据同步中间件。同步binlog的事件，把其中涉及到表的增量数据
处理成易于理解的文本信息，并将这些增量数据发送到消息队列等地方。常见应用如，MySQL到异构数据库的同步等。

* 分布式，多节点提高可用性。集群的数据存储在ETCD中。
* 图形化控制台(支持中文和英文)，可以在控制台完成各项配置，开启或停止同步等。可以在控制台观测数据流水线、集群、节点、全局事件等。

### 快速开始

* 安装ETCD。Binlogo的数据存储和选主等功能都是依赖于ETCD。

* 安装Binlogo. Binlogo的下载地址: [Download Address](https://github.com/jin06/binlogo/releases)

* 数据格式： [Binlogo输出的数据格式](/docs/1.0.*/message-format.md)

* 启动Binlogo
  * 配置文件：解压下载的程序，在程序目录configs中有配置文件： ${binlogo-path}/configs/binlogo.yaml
    
  * 修改etcd地址为你的etcd地址；
    
  * 修改node名称（如果node名称为空，会自动获取本机的hostname）
    
  ![avatar](/docs/assets/pic/edit_config_step1.en.png)

  * > $ ./binlogo server --config ./configs/binlogo.yaml

* 打开浏览器，进入控制台 : http://127.0.0.1:9999/console

* 创建流水线:

> 按照以下步骤就可以创建一个流水线，GTID模式需要MySQL已经开启了GTID。

![avatar](/docs/assets/pic/create_pipe_step1.en.png)

![avatar](/docs/assets/pic/create_pipe_step2.en.png)

* 运行流水线：

> 在流水线列表中点击运行按钮，或在流水线详情中运行。

![avatar](/docs/assets/pic/run_pipeline_step1.en.png)

* 在详情中观测流水线的运行状况

> 以下是一些流水线详情中的信息，流水线当前状态、同步到的binlog位置，一些流水线的事件等。


![avatar](/docs/assets/pic/pipeline_condition_step1.en.png)

![avatar](/docs/assets/pic/pipeline_condition_step2.en.png)

* 观察输出，这里演示的就直接输出到标准输出。

> 在MySQL中进行一些增删改的操作，在程序的标准输出中看看变化。

![avatar](/docs/assets/pic/output_step1.en.png)

![avatar](/docs/assets/pic/output_step2.en.png)

* 配置输出到kafka，下面一些kafka的建议配置

> * 为了更好的性能表现，可以配置如下。这种模式下有数据丢失风险。 
>   *  acks=1
>   *  enable.idempotence=false
>   *  compression.type=snappy
>   *  retries=0
> * 为了更好的可靠性，可以配置如下。其中如果enable.idempotence配置为true，会保证发送到kafka的数据是幂等。
>   * acks=-1
>   * enable.idempotence=true
>   * retries=3 or larger one

![avatar](/docs/assets/pic/output_kafka_step1.en.png)

![avatar](/docs/assets/pic/output_kafka_step2.en.png)

### Docker

- [Docker Hub](https://hub.docker.com/r/jin06/binlogo)

> $ docker pull jin06/binlogo
> </br>
> $ docker run -e "ETCD_ENDPOINTS=172.17.0.3:2379" --name BinlogoNode -it -d -p 9999:9999 jin06/binlogo:latest
> </br>
> Open browser access http://127.0.0.1:9999/console
> </br>
> 以下是我测试时使用docker启动的一个集群
>

![avatar](/docs/assets/pic/docker_step1.en.png)

### Kubernetes

- [文档](/docs/1.0.*/zh/instanll-kubernetes.md)

### 其他输出端

* [HTTP](/docs/1.0.*/zh/configure-http-output.md)
* [RabbitMQ](/docs/1.0.*/zh/configure-rabbitmq-outupt.md)
* [Kafka](/docs/1.0.*/zh/configure-kafka-output.md)
* [Redis](/docs/1.0.*/zh/configure-redis-outupt.md)
* [AliCloud RocketMQ](/docs/1.0.*/zh/configure-rocketmq-outupt.md)

### 文档

* [docs link](https://github.com/jin06/binlogo/wiki)

### 建议&问题反馈
* 欢迎提交bug或需求: [GitHub Issue](https://github.com/jin06/binlogo/issues)
* 如果有疑问可以直接联系我: jinlong4696@163.com
