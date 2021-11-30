> - Binlogo will send the message to the exchange of rabbitmq. The topic mode is used. The routing key is the format of database and data table, such as database.table.
> - Exchange name is created according to the configuration of the user pipeline.
> - RabbitMQ message queues can be created on the rabbitmq console or on the console
> - Reference [RabbitMQ Topic Pattern](https://www.rabbitmq.com/tutorials/tutorial-five-go.html)

- Configure pipeline to RabbitMQ

![avatar](/docs/assets/pic/config_output_rabbit.en.png)

- Create message queues directly in your code

> [RabbitMQ Consumer Example Code](https://github.com/jin06/binlogo/tree/master/examples/rabbitmq/main.go)

![avatar](/docs/assets/pic/config_output_rabbit2.png)

- Or creating message queues and viewing messages in rabbitmq

![avatar](/docs/assets/pic/config_output_rabbit3.png)

![avatar](/docs/assets/pic/config_output_rabbit4.png)









