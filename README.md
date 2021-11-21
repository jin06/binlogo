Binlogo
=====================================
[中文](README_zh.md) | [English](README.md)

Binlogo is the distributed, visualized application based on MySQL binlog.
In short, binlogo is to process the data changes of MySQL into easy to
understand messages and send them to different places according to
the user's configuration. This is some advantages:

* Distributed, multi node improves availability.
* Visualization, can complete common operations and
  observe the status of the whole cluster in the control background

### get-started

* Install etcd. Binlogo relies on etcd, so you must install etcd first.

* Install binlogo. Binlogo's download address: [Download Address]()

* Start binlogo.
  > $ ./binlogo server --config ./configs/binlogo.yaml 

* Open browser: http://127.0.0.1:9999/console

* Create Pipeline: 

> Follow the stemps. 

![avatar](/docs/assets/pic/create_pipe_step1.en.png)

![avatar](/docs/assets/pic/create_pipe_step2.en.png)

1. Run pipeline.

> Click button to run the pipeline instance. 

![avatar](/docs/assets/pic/create_pipe_step2.en.png)

2. Operation condition.

> You can see the operation condition of pipeline.
 
 
![avatar](/docs/assets/pic/pipeline_condition_step1.en.png)

![avatar](/docs/assets/pic/pipeline_condition_step2.en.png)

3. See the output 

> Insert some into mysql, watch the ouput on stdout.
 
![avatar](/docs/assets/pic/output_step1.en.png)

![avatar](/docs/assets/pic/output_step2.en.png)

4. Configuration output to Kafka

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

### docs

* [docs link](https://github.com/jin06/binlogo/wiki)

### questions
* To Report bug: https://github.com/jin06/binlogo/issues[GitHub Issue]
* Contact author: jinlong4696@163.com
