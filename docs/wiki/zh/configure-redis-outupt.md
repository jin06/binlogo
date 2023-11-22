>  binlogo会把数据发送到Redis的list，从list的右边插入数据，即 RPUSH
>  Redis的发布订阅模式在消费者离线时不会保存，所以采用的是list结构来存储数据。因为Redis本身并不是消息队列，在更加复杂或者对数据一致性等要求更高的业务上不建议使用Redis
> 

- 配置流水线到Redis

>  username一般不需要配置，Redis6.0及以上的版本中有用户名的概念 
> 
![avatar](/docs/assets/pic/config_output_redis.png)

- [Redis 消费者示例代码](https://github.com/jin06/binlogo/tree/master/examples/redis/main.go)










