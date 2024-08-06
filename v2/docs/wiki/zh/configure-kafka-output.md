### 配置流水线到Kafka

![avatar](/docs/wiki/assets/pic/config_output_kafka.png)

> binlogo列出了一些比较重要的配置，用户可以根据自己的业务配置高性能或者高可用的kafka

- 以下配置是一个高性能的配置，有丢失数据的风险
    - acks=1
      > 该配置的含义是，kafka只需要在主分片写入数据后即可返回。丢失数据的风险在于，如果此时该分片所在的节点崩溃，刚才写入的数据没有同步到副本分片中，那么这条数据在kafka就是丢失的。当然，数据不是百分之百丢失。
    - enable.idempotence=false
      > 这个是kafka的一个幂等性配置。每次给kafka发送消息时会有一个批次号，如果配置为true，kafka会根据该批次号对接收的数据进行幂等性校验，以此保证kafka只接收一次数据。
    - compression.type=snappy
      > 压缩算法，可以查看kafka官网，查看各种压缩算法的压缩比和性能等，选择适合自己业务的数据压缩方式
    - retries=0
      > 数据重试次数.

- 以下配置是一个高可用的配置，有可能多发送数据
    - acks=-1
      > 配置成-1，意味着需要kafka把该条数据写入到所有的副本分片后才返回。
    - enable.idempotence=false
      > 幂等性配置
    - retries=3 or larger one
      > 重试次数，越大肯定越可靠

- 以下配置在高可用配置基础上，保证数据只接收一次。
    - acks=-1
    - enable.idempotence=true
      > 和高可用配置的区别就是这个配置，开启幂等性的校验
    - retries=3 or larger one
