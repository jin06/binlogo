### Configure pipeline to Kafka

![avatar](/docs/assets/pic/config_output_kafka.png)

> Binlogo lists some important configurations. Users can configure high-performance or highly available Kafka according to their business

- The following configuration is a high-performance configuration with the risk of data loss
    - acks=1
      > The meaning of this configuration is that Kafka only needs to return after writing data to the primary partition. The risk of data loss is that if the node where the partition is located crashes and the data just written is not synchronized to the replica partition, the data is lost in Kafka. Of course, data is not 100% lost.
    - enable.idempotence=false
      > This is an idempotent configuration of Kafka. Each time a message is sent to Kafka, there will be a batch number. If the configuration is true, Kafka will check the idempotency of the received data according to the batch number to ensure that Kafka receives data only once.
    - compression.type=snappy
      > For compression algorithms, you can check the Kafka official website, check the compression ratio and performance of various compression algorithms, and select the data compression method suitable for your business
    - retries=0
      > Retries counts.

- The following configuration is a highly available configuration. It is possible to send more data
    - acks=-1
      > If it is configured as - 1, it means that Kafka needs to write the data to all partitions(primary and follows) before returning.
    - enable.idempotence=false
      > Idempotent configuration
    - retries=3 or larger one
      > The greater the number of retries, the more reliable it must be

- The following configuration is based on the high availability configuration to ensure that the data in Kafka is received only once.
    - acks=-1
    - enable.idempotence=true
      > The difference from the highly available configuration is that this configuration enables idempotency verification
    - retries=3 or larger one
