> - binlogo会把消息发送到RabbitMQ的exchange中，使用topic模式，routing key 是数据库加数据表的格式，形如database.table。
> - exchange name 根据用户流水线的配置创建。
> - RabbitMQ 的消息队列可以在RabbitMQ的控制台创建，或者在控制台创建
> - 参考 [RabbitMQ Topic Pattern](https://www.rabbitmq.com/tutorials/tutorial-five-go.html)

- 配置流水线到RabbitMQ

![avatar](/docs/assets/pic/config_output_rabbit.en.png)

- 在代码中直接创建消息队列

> [RabbitMQ 消费者示例代码](https://github.com/jin06/binlogo/tree/master/examples/rabbitmq/main.go)

![avatar](/docs/assets/pic/config_output_rabbit2.png)

- 也可以在RabbitMQ创建消息队列和查看消息

![avatar](/docs/assets/pic/config_output_rabbit3.png)

![avatar](/docs/assets/pic/config_output_rabbit4.png)









