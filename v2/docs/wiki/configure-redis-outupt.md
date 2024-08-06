>Binlogo will send data to redis's list and insert data from the right side of the list, CMD is RPUSH

>Redis's publish and subscribe mode will not be saved when consumers are offline,
> so the list structure is used to store data. 
> Because redis itself is not a message queue, it is not recommended to use redis for more complex services or services requiring higher data consistency

- Configure pipeline to Redis

> Generally, username does not need to be configured.
> Redis 6.0 and above have the concept of username
> 
![avatar](/docs/wiki/assets/pic/config_output_redis.png)

- [Redis Consumer Example Code](https://github.com/jin06/binlogo/tree/master/examples/redis/main.go)









